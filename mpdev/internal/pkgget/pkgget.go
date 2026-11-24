// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package pkgget implements a hardened replacement for kpt pkg get that blocks
// Git option/transport injection and refuses to dereference or copy symlinks.
package pkgget

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/GoogleCloudPlatform/marketplace-tools/mpdev/internal/fsutil"
	"sigs.k8s.io/kustomize/kyaml/fieldmeta"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/kio/filters"
	"sigs.k8s.io/kustomize/kyaml/kio/kioutil"
	"sigs.k8s.io/kustomize/kyaml/setters"
	"sigs.k8s.io/kustomize/kyaml/setters2/settersutil"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

const (
	kptFileName    = "Kptfile"
	kptAPIVersion  = "kpt.dev/v1alpha1"
	shortHandRef   = "$kpt-set"
	legacyPathAnno = "config.k8s.io/path"
)

// Options configures a hardened package fetch invocation.
type Options struct {
	// Pattern is the filename pattern used when reading resources from stdin ("-").
	Pattern string

	// AutoSet enables automatic KPT_SET_* and gcloud.* setter substitution.
	AutoSet bool

	// Clean removes an existing destination directory after validating that it
	// contains no symlinks.
	Clean bool

	// AllowLocalRepo permits cloning from local filesystem paths and file:// URLs.
	// It is always false in the CLI and used only by hermetic unit tests.
	AllowLocalRepo bool

	// Stdin is the reader used when args[0] is "-". Defaults to os.Stdin.
	Stdin io.Reader

	// Stdout receives status messages. Defaults to io.Discard when nil.
	Stdout io.Writer
}

// GitSpec records the upstream Git origin written to Kptfile.
type GitSpec struct {
	Commit    string `yaml:"commit,omitempty"`
	Repo      string `yaml:"repo,omitempty"`
	Directory string `yaml:"directory,omitempty"`
	Ref       string `yaml:"ref,omitempty"`
}

// GitTarget holds a validated Git source specification and local destination.
type GitTarget struct {
	GitSpec
	Destination string
}

type kptFileDoc struct {
	yaml.ResourceMeta `yaml:",inline"`
	Upstream          kptUpstream   `yaml:"upstream,omitempty"`
	PackageMeta       interface{}   `yaml:"packageMetadata,omitempty"`
	Dependencies      []interface{} `yaml:"dependencies,omitempty"`
	OpenAPI           interface{}   `yaml:"openAPI,omitempty"`
	Functions         interface{}   `yaml:"functions,omitempty"`
}

type kptUpstream struct {
	Type  string       `yaml:"type,omitempty"`
	Git   GitSpec      `yaml:"git,omitempty"`
	Stdin kptStdinSpec `yaml:"stdin,omitempty"`
}

type kptStdinSpec struct {
	FilenamePattern string `yaml:"filenamePattern,omitempty"`
	Original        string `yaml:"original,omitempty"`
}

// Run executes a hardened package fetch for positional arguments
// [REPO_URI[.git]/PKG_PATH[@VERSION] | "-", LOCAL_DEST_DIRECTORY].
func Run(args []string, opts Options) error {
	if len(args) != 2 {
		return fmt.Errorf("expected exactly 2 arguments, got %d", len(args))
	}
	if err := validateSafeChars("package source", args[0]); err != nil {
		return err
	}
	if err := validateSafeChars("destination directory", args[1]); err != nil {
		return err
	}
	if err := rejectLeadingHyphen("destination directory", args[1]); err != nil {
		return err
	}

	if args[0] == "-" {
		return runFromStdin(args[1], opts)
	}

	if err := rejectLeadingHyphen("package source", args[0]); err != nil {
		return err
	}

	target, err := ParseGitSpec(args[0], args[1], opts.AllowLocalRepo, opts.Clean)
	if err != nil {
		return err
	}

	stdout := opts.Stdout
	if stdout == nil {
		stdout = io.Discard
	}
	fmt.Fprintf(stdout, "fetching package %s from %s to %s\n",
		target.Directory, target.Repo, target.Destination)

	tmpClone, commitSHA, err := cloneAndVerifyRepo(target, opts.AllowLocalRepo)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpClone)

	pkgName := filepath.Base(target.Destination)
	postStage := func(stageDir string) error {
		if err := upsertGitKptfile(stageDir, pkgName, target.GitSpec, commitSHA); err != nil {
			return err
		}
		if opts.AutoSet {
			if err := PerformAutoSetters(stageDir); err != nil {
				return err
			}
		}
		return nil
	}

	return fsutil.StageAndCommitCopyDir(tmpClone, target.Directory, target.Destination, opts.Clean, postStage)
}

// ParseGitSpec parses and validates REPO_URI[.git]/PKG_PATH[@VERSION] and
// resolves LOCAL_DEST_DIRECTORY without dereferencing leaf symlinks.
func ParseGitSpec(rawURI, rawDest string, allowLocal, cleanDst bool) (GitTarget, error) {
	if err := validateSafeChars("package URI", rawURI); err != nil {
		return GitTarget{}, err
	}
	if err := rejectLeadingHyphen("package URI", rawURI); err != nil {
		return GitTarget{}, err
	}
	if err := validateSafeChars("destination directory", rawDest); err != nil {
		return GitTarget{}, err
	}
	if err := rejectLeadingHyphen("destination directory", rawDest); err != nil {
		return GitTarget{}, err
	}

	repo, rawDir, ref, err := splitRepoDirRef(rawURI)
	if err != nil {
		return GitTarget{}, err
	}

	if err := ValidateRepoURI(repo, allowLocal); err != nil {
		return GitTarget{}, err
	}
	if err := validateGitRef(ref); err != nil {
		return GitTarget{}, err
	}
	cleanDir, err := validateAndCleanSubDir(rawDir)
	if err != nil {
		return GitTarget{}, err
	}
	dest, err := resolveDestination(rawDest, repo, cleanDir, cleanDst)
	if err != nil {
		return GitTarget{}, err
	}

	return GitTarget{
		GitSpec: GitSpec{
			Repo:      repo,
			Directory: cleanDir,
			Ref:       ref,
		},
		Destination: dest,
	}, nil
}

// ValidateRepoURI enforces the remote transport allowlist (https, http, ssh,
// and SCP-style [user@]host:path) and rejects leading hyphens, control
// characters, ext::/fd:: remote helpers, and (unless allowLocal is true) local
// paths and file:// URLs.
func ValidateRepoURI(repo string, allowLocal bool) error {
	if err := validateSafeChars("repository URI", repo); err != nil {
		return err
	}
	if err := rejectLeadingHyphen("repository URI", repo); err != nil {
		return err
	}
	if strings.Contains(repo, "::") {
		return fmt.Errorf("git remote helper transport is not permitted in %q", repo)
	}

	if strings.Contains(repo, "://") {
		u, err := url.Parse(repo)
		if err != nil {
			return fmt.Errorf("invalid repository URL %q: %w", repo, err)
		}
		if u.Opaque != "" || u.RawQuery != "" || u.Fragment != "" {
			return fmt.Errorf("unsupported query, fragment, or opaque form in repository URL %q", repo)
		}
		scheme := strings.ToLower(u.Scheme)
		switch scheme {
		case "https", "http", "ssh":
			if u.Host == "" || u.Hostname() == "" {
				return fmt.Errorf("repository URL %q is missing a host", repo)
			}
			if err := validateUnescapedHostOrUser("repository host", u.Host); err != nil {
				return err
			}
			if err := validateUnescapedHostOrUser("repository hostname", u.Hostname()); err != nil {
				return err
			}
			if u.User != nil {
				if err := validateUnescapedHostOrUser("repository user", u.User.String()); err != nil {
					return err
				}
			}
			if u.Path == "" || u.Path == "/" {
				return fmt.Errorf("repository URL %q is missing a repository path", repo)
			}
			unescapedPath, unescapeErr := url.PathUnescape(u.Path)
			if unescapeErr != nil {
				return fmt.Errorf("invalid percent-encoding in repository URL %q: %w", repo, unescapeErr)
			}
			if err := validateSafeChars("repository URL path", unescapedPath); err != nil {
				return err
			}
			if err := rejectTraversalSegments(u.Path); err != nil {
				return err
			}
			if err := rejectTraversalSegments(unescapedPath); err != nil {
				return err
			}
			return nil

		case "file":
			if !allowLocal {
				return fmt.Errorf("file:// repository URL %q is not permitted", repo)
			}
			if u.Host != "" && u.Host != "localhost" {
				return fmt.Errorf("unsupported host %q in file:// URL", u.Host)
			}
			if u.Path == "" {
				return fmt.Errorf("empty path in file:// URL %q", repo)
			}
			return rejectTraversalSegments(u.Path)

		default:
			return fmt.Errorf("unsupported repository URL scheme %q in %q", u.Scheme, repo)
		}
	}

	// Check SCP-style SSH syntax: [user@]host:path (no '/' or '\' before ':').
	if colonIdx := strings.IndexByte(repo, ':'); colonIdx > 0 && !strings.ContainsAny(repo[:colonIdx], "/\\") {
		hostSpec := repo[:colonIdx]
		repoPath := repo[colonIdx+1:]
		if strings.Count(hostSpec, "@") > 1 {
			return fmt.Errorf("invalid user@host specification in %q", repo)
		}
		host := hostSpec
		if atIdx := strings.IndexByte(hostSpec, '@'); atIdx >= 0 {
			user := hostSpec[:atIdx]
			host = hostSpec[atIdx+1:]
			if user == "" {
				return fmt.Errorf("empty SSH username in %q", repo)
			}
			if err := rejectLeadingHyphen("SSH username", user); err != nil {
				return err
			}
		}
		if len(host) <= 1 || strings.ContainsAny(host, " \t") {
			return fmt.Errorf("invalid SSH host %q in %q", host, repo)
		}
		if err := rejectLeadingHyphen("SSH host", host); err != nil {
			return err
		}
		if repoPath == "" {
			return fmt.Errorf("empty repository path in %q", repo)
		}
		if err := rejectLeadingHyphen("SSH repository path", repoPath); err != nil {
			return err
		}
		if err := rejectTraversalSegments(repoPath); err != nil {
			return err
		}
		return nil
	}

	if !allowLocal {
		return fmt.Errorf("local filesystem repository path %q is not permitted; specify an https://, http://, or ssh:// URL", repo)
	}
	if strings.ContainsRune(repo, '\\') {
		return fmt.Errorf("backslash not permitted in repository path %q", repo)
	}
	return nil
}

func splitRepoDirRef(rawURI string) (repo, dir, ref string, err error) {
	if gitRepo, rest, count := splitAtGitDelimiter(rawURI); count > 0 {
		if count > 1 {
			return "", "", "", fmt.Errorf("ambiguous package URI %q containing multiple '.git' repository delimiters", rawURI)
		}
		repo = strings.TrimSuffix(gitRepo, "/")
		if strings.Count(rest, "@") > 1 {
			return "", "", "", fmt.Errorf("ambiguous '@version' in package URI %q", rawURI)
		}
		if strings.Contains(rest, "@") {
			subParts := strings.SplitN(rest, "@", 2)
			dir = subParts[0]
			ref = strings.TrimSuffix(subParts[1], "/")
			if ref == "" {
				return "", "", "", fmt.Errorf("empty '@version' in package URI %q", rawURI)
			}
		} else {
			dir = rest
			ref = "master"
		}
		if dir == "" {
			dir = "/"
		}
		return repo, dir, ref, nil
	}

	schemeIdx := strings.Index(rawURI, "://")
	if schemeIdx <= 0 || strings.Count(rawURI, "://") > 1 {
		return "", "", "", fmt.Errorf("ambiguous repo/dir@version in %q: specify '.git' in argument", rawURI)
	}
	restAfterScheme := rawURI[schemeIdx+3:]
	slashIdx := strings.IndexByte(restAfterScheme, '/')
	if slashIdx <= 0 {
		return "", "", "", fmt.Errorf("ambiguous repo/dir@version in %q: specify '.git' in argument", rawURI)
	}
	authorityPart := rawURI[:schemeIdx+3+slashIdx]
	pathAndRef := rawURI[schemeIdx+3+slashIdx:]
	if strings.Count(pathAndRef, "@") > 1 {
		return "", "", "", fmt.Errorf("ambiguous '@version' in package URI %q", rawURI)
	}

	rawPath := pathAndRef
	ref = "master"
	if strings.Contains(pathAndRef, "@") {
		subParts := strings.SplitN(pathAndRef, "@", 2)
		rawPath = subParts[0]
		ref = strings.TrimSuffix(subParts[1], "/")
		if ref == "" {
			return "", "", "", fmt.Errorf("empty '@version' in package URI %q", rawURI)
		}
	}

	u, parseErr := url.Parse(authorityPart + rawPath)
	if parseErr != nil || !strings.EqualFold(u.Hostname(), "github.com") {
		return "", "", "", fmt.Errorf("ambiguous repo/dir@version in %q: specify '.git' in argument", rawURI)
	}

	trimmedPath := strings.Trim(u.Path, "/")
	if trimmedPath == "" {
		return "", "", "", fmt.Errorf("ambiguous repo/dir@version in %q: specify '.git' in argument", rawURI)
	}
	segs := strings.Split(trimmedPath, "/")
	if len(segs) < 2 || segs[0] == "" || segs[1] == "" {
		return "", "", "", fmt.Errorf("ambiguous repo/dir@version in %q: specify '.git' in argument", rawURI)
	}

	repoURL := &url.URL{
		Scheme: u.Scheme,
		User:   u.User,
		Host:   u.Host,
		Path:   "/" + segs[0] + "/" + segs[1],
	}
	repo = repoURL.String()
	if len(segs) > 2 {
		dir = path.Join(segs[2:]...)
	} else {
		dir = "/"
	}
	return repo, dir, ref, nil
}

// splitAtGitDelimiter locates ".git" only when it ends the repository name
// (followed by '/', '@', or end-of-string), avoiding false splits on ".github",
// ".gitignore", or ".gitmodules" path components.
func splitAtGitDelimiter(rawURI string) (repo, rest string, count int) {
	searchFrom := 0
	firstIdx := -1
	for searchFrom < len(rawURI) {
		relIdx := strings.Index(rawURI[searchFrom:], ".git")
		if relIdx < 0 {
			break
		}
		idx := searchFrom + relIdx
		after := idx + len(".git")
		if after == len(rawURI) || rawURI[after] == '/' || rawURI[after] == '@' {
			count++
			if firstIdx < 0 {
				firstIdx = idx
			}
		}
		searchFrom = after
	}
	if firstIdx < 0 {
		return "", "", 0
	}
	return rawURI[:firstIdx], rawURI[firstIdx+len(".git"):], count
}

func validateGitRef(ref string) error {
	if err := validateSafeChars("git ref", ref); err != nil {
		return err
	}
	if err := rejectLeadingHyphen("git ref", ref); err != nil {
		return err
	}
	if strings.Contains(ref, "..") || strings.Contains(ref, "@{") || strings.ContainsAny(ref, " \t~^:?*[\\") {
		return fmt.Errorf("invalid characters or sequence in git ref %q", ref)
	}
	if strings.HasSuffix(ref, ".lock") || strings.HasSuffix(ref, ".") || strings.HasSuffix(ref, "/") {
		return fmt.Errorf("invalid suffix in git ref %q", ref)
	}
	return nil
}

func validateAndCleanSubDir(rawDir string) (string, error) {
	if rawDir == "" {
		return "/", nil
	}
	if err := validateSafeChars("package subdirectory", rawDir); err != nil {
		return "", err
	}
	if err := rejectLeadingHyphen("package subdirectory", rawDir); err != nil {
		return "", err
	}
	trimmed := strings.TrimLeft(rawDir, "/")
	if trimmed != "" {
		if err := rejectLeadingHyphen("package subdirectory", trimmed); err != nil {
			return "", err
		}
	}
	if err := rejectTraversalSegments(rawDir); err != nil {
		return "", err
	}
	for _, seg := range strings.Split(trimmed, "/") {
		if strings.HasPrefix(seg, "-") {
			return "", fmt.Errorf("package subdirectory segment %q must not start with '-'", seg)
		}
	}
	if strings.HasPrefix(rawDir, "/") {
		return path.Clean(rawDir), nil
	}
	cleaned := path.Clean(rawDir)
	if cleaned == "." {
		return "/", nil
	}
	return cleaned, nil
}

func resolveDestination(rawDest, repo, subDir string, cleanDst bool) (string, error) {
	if err := fsutil.RejectLeafSymlink(rawDest); err != nil {
		return "", err
	}
	cleaned := filepath.Clean(rawDest)

	info, err := os.Lstat(cleaned)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}
		parent := filepath.Dir(cleaned)
		if err := fsutil.RejectLeafSymlink(parent); err != nil {
			return "", err
		}
		parentInfo, statErr := os.Stat(parent)
		if statErr != nil {
			if errors.Is(statErr, os.ErrNotExist) {
				return "", fmt.Errorf("parent directory %s does not exist", parent)
			}
			return "", statErr
		}
		if !parentInfo.IsDir() {
			return "", fmt.Errorf("parent %s is not a directory", parent)
		}
		return cleaned, nil
	}

	if info.Mode()&os.ModeSymlink != 0 {
		return "", fmt.Errorf("destination %q: %w", rawDest, fsutil.ErrSymlink)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("LOCAL_DEST_DIRECTORY %q must be a directory", rawDest)
	}

	// If cleanDst is requested on an explicit destination directory (not "."),
	// StageAndCommitCopyDir will validate and replace cleaned directly.
	if cleanDst && cleaned != "." && cleaned != string(filepath.Separator) {
		return cleaned, nil
	}

	// Match kpt pkg get behavior when pointing at "." or an existing parent dir:
	// derive a child directory from the package subdirectory or repository basename.
	base := derivePkgBaseName(repo, subDir)
	if base == "" || base == "." || base == ".." || base == "/" {
		return "", fmt.Errorf("unable to derive destination directory name from %q and %q", repo, subDir)
	}
	derived := filepath.Join(cleaned, base)
	if err := fsutil.RejectLeafSymlink(derived); err != nil {
		return "", err
	}
	if _, err := os.Lstat(derived); err == nil && !cleanDst {
		return "", fmt.Errorf("destination directory %s already exists", derived)
	}
	return derived, nil
}

func derivePkgBaseName(repo, subDir string) string {
	trimmedSub := strings.Trim(subDir, "/")
	if trimmedSub != "" && trimmedSub != "." {
		return path.Base(trimmedSub)
	}
	trimmedRepo := strings.TrimSuffix(strings.TrimSuffix(repo, "/"), ".git")
	if idx := strings.LastIndexAny(trimmedRepo, "/:"); idx >= 0 {
		return trimmedRepo[idx+1:]
	}
	return path.Base(trimmedRepo)
}

func cloneAndVerifyRepo(target GitTarget, allowLocal bool) (string, string, error) {
	gitProgram, err := exec.LookPath("git")
	if err != nil {
		return "", "", fmt.Errorf("locating git binary: %w", err)
	}

	var candidates []string
	trimmedSub := strings.Trim(target.Directory, "/")
	if trimmedSub != "" && trimmedSub != "." && !strings.Contains(target.Ref, "refs") {
		prefixed := path.Join(trimmedSub, target.Ref)
		if validateGitRef(prefixed) == nil && prefixed != target.Ref {
			candidates = append(candidates, prefixed)
		}
	}
	candidates = append(candidates, target.Ref)

	var lastErr error
	for _, candidateRef := range candidates {
		tmpClone, commitSHA, cloneErr := cloneSingleRef(gitProgram, target.Repo, candidateRef, target.Directory, allowLocal)
		if cloneErr == nil {
			return tmpClone, commitSHA, nil
		}
		lastErr = cloneErr
		if errors.Is(cloneErr, fsutil.ErrSymlink) || errors.Is(cloneErr, fsutil.ErrSpecialFile) || errors.Is(cloneErr, fsutil.ErrPathEscapesRoot) {
			return "", "", cloneErr
		}
	}
	return "", "", fmt.Errorf("cloning git repository %s at %s: %w", target.Repo, target.Ref, lastErr)
}

func cloneSingleRef(gitProgram, repo, ref, subDir string, allowLocal bool) (string, string, error) {
	tmpClone, err := os.MkdirTemp("", "mpdev-pkg-get-*")
	if err != nil {
		return "", "", fmt.Errorf("creating temporary clone directory: %w", err)
	}
	succeeded := false
	defer func() {
		if !succeeded {
			_ = os.RemoveAll(tmpClone)
		}
	}()

	env := sanitizedGitEnv(allowLocal)
	cfgArgs := gitConfigArgs(allowLocal)

	runGit := func(args ...string) ([]byte, error) {
		fullArgs := make([]string, 0, len(cfgArgs)+len(args)+2)
		fullArgs = append(fullArgs, "-C", tmpClone)
		fullArgs = append(fullArgs, cfgArgs...)
		fullArgs = append(fullArgs, args...)
		cmd := exec.Command(gitProgram, fullArgs...)
		cmd.Env = env
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if runErr := cmd.Run(); runErr != nil {
			return nil, fmt.Errorf("git %s failed: %w (%s)", args[0], runErr, strings.TrimSpace(out.String()))
		}
		return out.Bytes(), nil
	}

	initArgs := make([]string, 0, len(cfgArgs)+3)
	initArgs = append(initArgs, cfgArgs...)
	initArgs = append(initArgs, "init", "--", tmpClone)
	initCmd := exec.Command(gitProgram, initArgs...)
	initCmd.Env = env
	if initOut, initErr := initCmd.CombinedOutput(); initErr != nil {
		return "", "", fmt.Errorf("git init failed: %w (%s)", initErr, strings.TrimSpace(string(initOut)))
	}

	if _, err := runGit("config", "core.symlinks", "false"); err != nil {
		return "", "", err
	}
	if _, err := runGit("remote", "add", "--", "origin", repo); err != nil {
		return "", "", err
	}

	shallowErr := func() error {
		if _, fetchErr := runGit("fetch", "--depth=1", "--", "origin", ref); fetchErr != nil {
			return fetchErr
		}
		_, resetErr := runGit("reset", "--hard", "FETCH_HEAD", "--")
		return resetErr
	}()
	if shallowErr != nil {
		if _, fetchErr := runGit("fetch", "--", "origin"); fetchErr != nil {
			return "", "", fetchErr
		}
		if _, resetErr := runGit("reset", "--hard", ref, "--"); resetErr != nil {
			if strings.HasPrefix(ref, "origin/") {
				return "", "", resetErr
			}
			if _, branchErr := runGit("reset", "--hard", "origin/"+ref, "--"); branchErr != nil {
				return "", "", resetErr
			}
		}
	}

	// Inspect parent Git index for mode-120000 symlinks before touching submodules.
	if err := checkGitIndexNoSymlinks(runGit, subDir, false); err != nil {
		return "", "", err
	}

	if _, statErr := os.Lstat(filepath.Join(tmpClone, ".gitmodules")); statErr == nil {
		if err := validateGitmodules(tmpClone, runGit, allowLocal); err != nil {
			return "", "", err
		}
		if _, subErr := runGit("submodule", "update", "--init", "--recursive", "--"); subErr != nil {
			return "", "", subErr
		}
		if err := checkGitIndexNoSymlinks(runGit, subDir, true); err != nil {
			return "", "", err
		}
	}

	revOut, err := runGit("rev-parse", "--verify", "HEAD")
	if err != nil {
		return "", "", err
	}
	commitSHA := strings.TrimSpace(string(revOut))
	if commitSHA == "" {
		return "", "", errors.New("git rev-parse returned empty commit SHA")
	}

	succeeded = true
	return tmpClone, commitSHA, nil
}

func checkGitIndexNoSymlinks(runGit func(args ...string) ([]byte, error), subDir string, recurseSubmodules bool) error {
	lsArgs := []string{"ls-files", "-s", "-z"}
	if recurseSubmodules {
		lsArgs = append(lsArgs, "--recurse-submodules")
	}
	out, err := runGit(lsArgs...)
	if err != nil {
		return err
	}

	targetSub := strings.Trim(filepath.ToSlash(subDir), "/")
	records := bytes.Split(out, []byte{0})
	for _, rawRec := range records {
		if len(rawRec) == 0 {
			continue
		}
		rec := string(rawRec)
		tabIdx := strings.IndexByte(rec, '\t')
		if tabIdx <= 0 {
			return fmt.Errorf("malformed git ls-files record %q", rec)
		}
		fields := strings.Fields(rec[:tabIdx])
		if len(fields) == 0 {
			return fmt.Errorf("missing mode in git ls-files record %q", rec)
		}
		mode := fields[0]
		entryPath := rec[tabIdx+1:]
		if !isIndexEntryRelevant(entryPath, targetSub) {
			continue
		}
		switch mode {
		case "120000":
			return fmt.Errorf("repository entry %q has git symlink mode 120000: %w", entryPath, fsutil.ErrSymlink)
		case "100644", "100755", "160000":
			continue
		default:
			return fmt.Errorf("repository entry %q has unsupported git mode %s: %w", entryPath, mode, fsutil.ErrSpecialFile)
		}
	}
	return nil
}

func isIndexEntryRelevant(entryPath, targetSub string) bool {
	if entryPath == ".gitmodules" {
		return true
	}
	if targetSub == "" || targetSub == "." {
		return true
	}
	if entryPath == targetSub {
		return true
	}
	if strings.HasPrefix(entryPath, targetSub+"/") {
		return true
	}
	// Reject ancestor directory symlinks (e.g. entryPath="escape" when targetSub="escape/secretdir").
	if strings.HasPrefix(targetSub, entryPath+"/") {
		return true
	}
	return false
}

func gitConfigArgs(allowLocal bool) []string {
	args := []string{"-c", "core.symlinks=false"}
	if allowLocal {
		args = append(args, "-c", "protocol.file.allow=always")
	} else {
		args = append(args, "-c", "protocol.file.allow=never")
	}
	return args
}

func sanitizedGitEnv(allowLocal bool) []string {
	rawEnv := os.Environ()
	filtered := make([]string, 0, len(rawEnv)+2)
	for _, kv := range rawEnv {
		key, _, _ := strings.Cut(kv, "=")
		switch {
		case key == "GIT_DIR",
			key == "GIT_WORK_TREE",
			key == "GIT_INDEX_FILE",
			key == "GIT_OBJECT_DIRECTORY",
			key == "GIT_ALTERNATE_OBJECT_DIRECTORIES",
			key == "GIT_COMMON_DIR",
			key == "GIT_ALLOW_PROTOCOL",
			key == "GIT_TERMINAL_PROMPT",
			strings.HasPrefix(key, "GIT_CONFIG"):
			continue
		default:
			filtered = append(filtered, kv)
		}
	}
	filtered = append(filtered, "GIT_TERMINAL_PROMPT=0")
	if !allowLocal {
		filtered = append(filtered, "GIT_ALLOW_PROTOCOL=https:http:ssh")
	}
	return filtered
}

func upsertGitKptfile(stageDir, pkgName string, gitSpec GitSpec, commitSHA string) error {
	doc := kptFileDoc{
		ResourceMeta: yaml.ResourceMeta{
			APIVersion: kptAPIVersion,
			Kind:       kptFileName,
			ObjectMeta: yaml.ObjectMeta{Name: pkgName},
		},
	}

	existing, _, err := fsutil.ReadRegularFileNoFollow(stageDir, kptFileName)
	if err == nil {
		dec := yaml.NewDecoder(bytes.NewReader(existing))
		if decodeErr := dec.Decode(&doc); decodeErr != nil {
			return fmt.Errorf("parsing existing %s: %w", kptFileName, decodeErr)
		}
		if doc.APIVersion == "" {
			doc.APIVersion = kptAPIVersion
		}
		if doc.Kind == "" {
			doc.Kind = kptFileName
		}
		if doc.ObjectMeta.Name == "" {
			doc.ObjectMeta.Name = pkgName
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	gitSpec.Commit = commitSHA
	doc.Upstream = kptUpstream{
		Type: "git",
		Git:  gitSpec,
	}

	return writeFormattedKptfile(stageDir, doc, 0644)
}

func writeFormattedKptfile(stageDir string, doc kptFileDoc, perm os.FileMode) error {
	rawYAML, err := yaml.Marshal(doc)
	if err != nil {
		return fmt.Errorf("marshaling %s: %w", kptFileName, err)
	}
	rNode, err := yaml.Parse(string(rawYAML))
	if err != nil {
		return fmt.Errorf("formatting %s: %w", kptFileName, err)
	}
	formatted, err := rNode.String()
	if err != nil {
		return fmt.Errorf("serializing %s: %w", kptFileName, err)
	}
	return fsutil.WriteRegularFileNoFollow(stageDir, kptFileName, []byte(formatted), perm)
}

// PerformAutoSetters applies KPT_SET_* environment variables and non-empty
// gcloud configuration values to YAML resources inside pkgDir without running
// Starlark or container functions.
func PerformAutoSetters(pkgDir string) error {
	fieldmeta.SetShortHandRef(shortHandRef)

	applySetter := func(name, value string) (int, error) {
		rw := &kio.LocalPackageReadWriter{
			PackagePath:           pkgDir,
			KeepReaderAnnotations: false,
			IncludeSubpackages:    true,
		}
		s1 := &setters.PerformSetters{Name: name, Value: value, SetBy: "kpt"}
		if err := (kio.Pipeline{
			Inputs:  []kio.Reader{rw},
			Filters: []kio.Filter{s1},
			Outputs: []kio.Writer{rw},
		}).Execute(); err != nil {
			return 0, err
		}

		s2 := &settersutil.FieldSetter{
			Name:          name,
			Value:         value,
			SetBy:         "kpt",
			ResourcesPath: pkgDir,
			OpenAPIPath:   filepath.Join(pkgDir, kptFileName),
		}
		if err := (kio.Pipeline{
			Inputs:  []kio.Reader{rw},
			Filters: []kio.Filter{s2},
			Outputs: []kio.Writer{rw},
		}).Execute(); err != nil {
			return 0, err
		}
		return s1.Count + s2.Count, nil
	}

	for _, envKV := range os.Environ() {
		if !strings.HasPrefix(envKV, "KPT_SET_") {
			continue
		}
		k, v, ok := strings.Cut(envKV, "=")
		if !ok {
			continue
		}
		setterName := strings.ToLower(strings.TrimPrefix(k, "KPT_SET_"))
		if setterName == "" {
			continue
		}
		if _, err := applySetter(setterName, v); err != nil {
			return err
		}
	}

	gcloudProps := []string{"compute.region", "compute.zone", "core.project"}
	projectID := ""
	for _, prop := range gcloudProps {
		out, err := exec.Command("gcloud", "config", "list", "--format", fmt.Sprintf("value(%s)", prop)).Output()
		if err != nil {
			continue
		}
		val := strings.TrimSpace(string(out))
		if val == "" {
			continue
		}
		count, err := applySetter("gcloud."+prop, val)
		if err != nil {
			return err
		}
		if prop == "core.project" && count > 0 {
			projectID = val
		}
	}

	if projectID != "" {
		out, err := exec.Command("gcloud", "projects", "describe", projectID, "--format", "value(projectNumber)").Output()
		if err == nil {
			if projNum := strings.TrimSpace(string(out)); projNum != "" {
				if _, setErr := applySetter("gcloud.project.projectNumber", projNum); setErr != nil {
					return setErr
				}
			}
		}
	}

	return nil
}

func runFromStdin(rawDest string, opts Options) error {
	pattern := opts.Pattern
	if pattern == "" {
		pattern = filters.DefaultFilenamePattern
	}
	if err := validateSafeChars("filename pattern", pattern); err != nil {
		return err
	}
	if filepath.IsAbs(pattern) || strings.HasPrefix(pattern, "/") || strings.ContainsRune(pattern, '\\') {
		return fmt.Errorf("filename pattern %q must be a relative path: %w", pattern, fsutil.ErrPathEscapesRoot)
	}
	if err := rejectTraversalSegments(pattern); err != nil {
		return err
	}

	if err := fsutil.RejectLeafSymlink(rawDest); err != nil {
		return err
	}
	dstAbs, err := filepath.Abs(rawDest)
	if err != nil {
		return err
	}
	dstBase := filepath.Base(dstAbs)
	if dstBase == "." || dstBase == ".." || dstBase == string(filepath.Separator) {
		return fmt.Errorf("invalid destination directory %q: %w", rawDest, fsutil.ErrPathEscapesRoot)
	}
	parentReal, err := fsutil.CanonicalDir(filepath.Dir(dstAbs))
	if err != nil {
		return fmt.Errorf("destination parent directory: %w", err)
	}
	dstTarget := filepath.Join(parentReal, dstBase)

	if info, lstatErr := os.Lstat(dstTarget); lstatErr == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("destination %q: %w", rawDest, fsutil.ErrSymlink)
		}
		if !opts.Clean {
			return fmt.Errorf("destination directory %s already exists", dstTarget)
		}
		if !info.IsDir() {
			return fmt.Errorf("destination %q exists and is not a directory", dstTarget)
		}
		if err := fsutil.RejectSymlinksInTree(dstTarget, false); err != nil {
			return err
		}
	} else if !errors.Is(lstatErr, os.ErrNotExist) {
		return lstatErr
	}

	input := opts.Stdin
	if input == nil {
		input = os.Stdin
	}

	stageDir, err := os.MkdirTemp(parentReal, ".mpdev-pkg-stdin-*")
	if err != nil {
		return err
	}
	committed := false
	defer func() {
		if !committed {
			_ = os.RemoveAll(stageDir)
		}
	}()
	if err := os.Chmod(stageDir, 0700); err != nil {
		return err
	}

	var buf bytes.Buffer
	var guardErr error
	fileSetter := &filters.FileSetter{FilenamePattern: pattern, Mode: fmt.Sprintf("%d", 0600)}
	pathGuard := kio.FilterFunc(func(nodes []*yaml.RNode) ([]*yaml.RNode, error) {
		for _, n := range nodes {
			meta, metaErr := n.GetMeta()
			if metaErr != nil {
				return nil, metaErr
			}
			for _, key := range []string{kioutil.PathAnnotation, legacyPathAnno} {
				if p, ok := meta.Annotations[key]; ok && p != "" {
					if _, pathErr := fsutil.ValidateContainedRelPath(stageDir, p); pathErr != nil {
						guardErr = fmt.Errorf("unsafe resource path annotation %q: %w", p, pathErr)
						return nil, guardErr
					}
				}
			}
		}
		return nodes, nil
	})

	err = (kio.Pipeline{
		Inputs:  []kio.Reader{&kio.ByteReader{Reader: input}},
		Filters: []kio.Filter{fileSetter, pathGuard, filters.FormatFilter{}},
		Outputs: []kio.Writer{
			kio.ByteWriter{Writer: &buf, KeepReaderAnnotations: true},
			kio.LocalPackageWriter{PackagePath: stageDir},
		},
	}).Execute()
	if guardErr != nil {
		return guardErr
	}
	if err != nil {
		return err
	}

	doc := kptFileDoc{
		ResourceMeta: yaml.ResourceMeta{
			APIVersion: kptAPIVersion,
			Kind:       kptFileName,
			ObjectMeta: yaml.ObjectMeta{Name: dstBase},
		},
		Upstream: kptUpstream{
			Type: "stdin",
			Stdin: kptStdinSpec{
				FilenamePattern: fileSetter.FilenamePattern,
				Original:        buf.String(),
			},
		},
	}
	if err := writeFormattedKptfile(stageDir, doc, 0600); err != nil {
		return err
	}
	if err := fsutil.RejectSymlinksInTree(stageDir, false); err != nil {
		return err
	}

	if opts.Clean {
		if _, lstatErr := os.Lstat(dstTarget); lstatErr == nil {
			if err := os.RemoveAll(dstTarget); err != nil {
				return err
			}
		}
	}
	if err := os.Rename(stageDir, dstTarget); err != nil {
		return err
	}
	committed = true
	return nil
}

func validateSafeChars(label, s string) error {
	if s == "" {
		return fmt.Errorf("%s must not be empty", label)
	}
	for i := 0; i < len(s); i++ {
		b := s[i]
		if b < 0x20 || b == 0x7f {
			return fmt.Errorf("%s contains control character 0x%02x", label, b)
		}
	}
	return nil
}

func rejectLeadingHyphen(label, s string) error {
	if strings.HasPrefix(s, "-") || strings.HasPrefix(strings.TrimSpace(s), "-") {
		return fmt.Errorf("%s %q must not start with '-'", label, s)
	}
	return nil
}

func validateUnescapedHostOrUser(label, raw string) error {
	if err := rejectLeadingHyphen(label, raw); err != nil {
		return err
	}
	unescaped, err := url.PathUnescape(raw)
	if err != nil {
		return fmt.Errorf("invalid percent-encoding in %s %q: %w", label, raw, err)
	}
	if err := validateSafeChars(label, unescaped); err != nil {
		return err
	}
	return rejectLeadingHyphen(label, unescaped)
}

func validateGitmodules(tmpClone string, runGit func(args ...string) ([]byte, error), allowLocal bool) error {
	realClone, err := fsutil.CanonicalDir(tmpClone)
	if err != nil {
		return err
	}
	if _, _, err := fsutil.ReadRegularFileNoFollow(realClone, ".gitmodules"); err != nil {
		return fmt.Errorf("validating .gitmodules: %w", err)
	}
	gmPath := filepath.Join(realClone, ".gitmodules")
	out, err := runGit("config", "--file", gmPath, "--null", "--list")
	if err != nil {
		return fmt.Errorf("inspecting .gitmodules: %w", err)
	}
	for _, entry := range bytes.Split(out, []byte{0}) {
		if len(entry) == 0 {
			continue
		}
		kv := string(entry)
		key, val, ok := strings.Cut(kv, "\n")
		if !ok {
			return fmt.Errorf("malformed .gitmodules config entry %q", kv)
		}
		keyLower := strings.ToLower(key)
		switch {
		case strings.HasSuffix(keyLower, ".url"):
			if err := ValidateRepoURI(val, allowLocal); err != nil {
				return fmt.Errorf("unsafe submodule URL in .gitmodules (%s): %w", key, err)
			}
		case strings.HasSuffix(keyLower, ".path"):
			if _, err := validateAndCleanSubDir(val); err != nil {
				return fmt.Errorf("unsafe submodule path in .gitmodules (%s): %w", key, err)
			}
		case strings.HasSuffix(keyLower, ".update"):
			if strings.HasPrefix(strings.TrimSpace(val), "!") {
				return fmt.Errorf("unsafe submodule update command in .gitmodules (%s)", key)
			}
		}
	}
	return nil
}

func rejectTraversalSegments(p string) error {
	if strings.ContainsRune(p, '\\') {
		return fmt.Errorf("backslash not permitted in path %q: %w", p, fsutil.ErrPathEscapesRoot)
	}
	for _, seg := range strings.Split(p, "/") {
		if seg == ".." {
			return fmt.Errorf("parent traversal segment '..' not permitted in %q: %w", p, fsutil.ErrPathEscapesRoot)
		}
	}
	return nil
}
