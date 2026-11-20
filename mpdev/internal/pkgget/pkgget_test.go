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

package pkgget

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/marketplace-tools/mpdev/internal/fsutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateAndParseRejectsMaliciousInputs(t *testing.T) {
	destDir := filepath.Join(t.TempDir(), "out")

	rejectedRemoteOnly := []struct {
		name string
		uri  string
		dest string
	}{
		{
			name: "leading hyphen source option injection",
			uri:  "--upload-pack=touch /tmp/pwn",
			dest: destDir,
		},
		{
			name: "leading hyphen destination",
			uri:  "https://github.com/GoogleCloudPlatform/marketplace-tools.git/examples/tf",
			dest: "-rf",
		},
		{
			name: "ref option injection after at-sign",
			uri:  "https://github.com/GoogleCloudPlatform/marketplace-tools.git/examples/tf@--upload-pack=id",
			dest: destDir,
		},
		{
			name: "leading hyphen subdirectory",
			uri:  "https://github.com/GoogleCloudPlatform/marketplace-tools.git/-dash",
			dest: destDir,
		},
		{
			name: "subdirectory parent traversal",
			uri:  "https://github.com/GoogleCloudPlatform/marketplace-tools.git/examples/../../etc",
			dest: destDir,
		},
		{
			name: "file scheme disallowed in CLI mode",
			uri:  "file:///tmp/localrepo.git/pkg",
			dest: destDir,
		},
		{
			name: "bare local path disallowed in CLI mode",
			uri:  "/tmp/localrepo.git/pkg",
			dest: destDir,
		},
		{
			name: "ext remote helper transport",
			uri:  "ext::sh -c touch% /tmp/pwn.git/pkg",
			dest: destDir,
		},
		{
			name: "git daemon cleartext scheme",
			uri:  "git://example.com/org/repo.git/pkg",
			dest: destDir,
		},
		{
			name: "ssh host option injection",
			uri:  "ssh://-oProxyCommand=id/org/repo.git/pkg",
			dest: destDir,
		},
		{
			name: "ssh user option injection",
			uri:  "ssh://-oProxyCommand=id@example.com/org/repo.git/pkg",
			dest: destDir,
		},
		{
			name: "percent-encoded ssh host option injection (CVE-2018-17456 style)",
			uri:  "ssh://%2doProxyCommand=id/org/repo.git/pkg",
			dest: destDir,
		},
		{
			name: "percent-encoded ssh user option injection",
			uri:  "ssh://%2du@example.com/org/repo.git/pkg",
			dest: destDir,
		},
		{
			name: "scp style host option injection",
			uri:  "-oProxyCommand=id:org/repo.git/pkg",
			dest: destDir,
		},
		{
			name: "scp style user option injection",
			uri:  "-u@github.com:org/repo.git/pkg",
			dest: destDir,
		},
		{
			name: "control character in URI",
			uri:  "https://github.com/org/repo.git/pkg\n--upload-pack=id",
			dest: destDir,
		},
		{
			name: "NUL byte in destination",
			uri:  "https://github.com/org/repo.git/pkg",
			dest: destDir + "\x00evil",
		},
	}

	for _, tc := range rejectedRemoteOnly {
		t.Run(tc.name, func(t *testing.T) {
			err := Run([]string{tc.uri, tc.dest}, Options{AllowLocalRepo: false})
			require.Error(t, err)
			if !strings.ContainsRune(tc.dest, 0) {
				_, statErr := os.Lstat(tc.dest)
				assert.ErrorIs(t, statErr, os.ErrNotExist)
			} else {
				_, statErr := os.Lstat(destDir)
				assert.ErrorIs(t, statErr, os.ErrNotExist)
			}
		})
	}

	// Even when AllowLocalRepo is true (unit tests), option/transport/traversal payloads fail.
	for _, uri := range []string{
		"ext::sh -c id.git/pkg",
		"ssh://-oProxyCommand=id/org/repo.git/pkg",
		"/tmp/repo.git/pkg@--upload-pack=id",
		"/tmp/repo.git/-dash",
		"/tmp/repo.git/a/../../b",
	} {
		_, err := ParseGitSpec(uri, destDir, true, false)
		require.Error(t, err, "expected error for %q", uri)
	}
}

func TestParseGitSpecValidForms(t *testing.T) {
	tmp := t.TempDir()

	t1, err := ParseGitSpec(
		"https://github.com/GoogleCloudPlatform/marketplace-tools.git/examples/deployment-manager/autogen/singlevm/",
		filepath.Join(tmp, "wordpress"),
		false,
		false,
	)
	require.NoError(t, err)
	assert.Equal(t, "https://github.com/GoogleCloudPlatform/marketplace-tools", t1.Repo)
	assert.Equal(t, "/examples/deployment-manager/autogen/singlevm", t1.Directory)
	assert.Equal(t, "master", t1.Ref)
	assert.Equal(t, filepath.Join(tmp, "wordpress"), t1.Destination)

	t2, err := ParseGitSpec(
		"https://github.com/GoogleCloudPlatform/marketplace-tools/examples/terraform/minimal-example@v1.2.3",
		filepath.Join(tmp, "tfpkg"),
		false,
		false,
	)
	require.NoError(t, err)
	assert.Equal(t, "https://github.com/GoogleCloudPlatform/marketplace-tools", t2.Repo)
	assert.Equal(t, "examples/terraform/minimal-example", t2.Directory)
	assert.Equal(t, "v1.2.3", t2.Ref)

	t3, err := ParseGitSpec(
		"git@github.com:GoogleCloudPlatform/marketplace-tools.git/examples/tf@main",
		filepath.Join(tmp, "scppkg"),
		false,
		false,
	)
	require.NoError(t, err)
	assert.Equal(t, "git@github.com:GoogleCloudPlatform/marketplace-tools", t3.Repo)
	assert.Equal(t, "/examples/tf", t3.Directory)
	assert.Equal(t, "main", t3.Ref)

	t4, err := ParseGitSpec(
		"https://github.com/GoogleCloudPlatform/marketplace-tools.git/.github/workflows@main",
		filepath.Join(tmp, "ghworkflows"),
		false,
		false,
	)
	require.NoError(t, err)
	assert.Equal(t, "https://github.com/GoogleCloudPlatform/marketplace-tools", t4.Repo)
	assert.Equal(t, "/.github/workflows", t4.Directory)
	assert.Equal(t, "main", t4.Ref)
}

func TestRunRejectsRepoEmbeddedSymlinks(t *testing.T) {
	canaryDir := t.TempDir()
	canarySecret := "TOP_SECRET_HOST_CREDENTIAL_DO_NOT_READ"
	canaryFile := filepath.Join(canaryDir, "secret.txt")
	require.NoError(t, os.WriteFile(canaryFile, []byte(canarySecret), 0600))
	nestedCanaryDir := filepath.Join(canaryDir, "secretdir")
	require.NoError(t, os.MkdirAll(nestedCanaryDir, 0700))
	require.NoError(t, os.WriteFile(filepath.Join(nestedCanaryDir, "token.yaml"), []byte(canarySecret), 0600))

	t.Run("direct absolute file symlink", func(t *testing.T) {
		repoDir := initTestGitRepo(t, func(work string) {
			require.NoError(t, os.MkdirAll(filepath.Join(work, "pkg"), 0755))
			require.NoError(t, os.WriteFile(filepath.Join(work, "pkg", "ok.yaml"), []byte("kind: ConfigMap\n"), 0644))
			require.NoError(t, os.Symlink(canaryFile, filepath.Join(work, "pkg", "leak.yaml")))
		})
		dst := filepath.Join(t.TempDir(), "fetched")
		err := Run([]string{repoDir + ".git/pkg@master", dst}, Options{AllowLocalRepo: true})
		require.ErrorIs(t, err, fsutil.ErrSymlink)
		_, statErr := os.Lstat(dst)
		assert.ErrorIs(t, statErr, os.ErrNotExist)
	})

	t.Run("relative out-of-tree file symlink", func(t *testing.T) {
		repoDir := initTestGitRepo(t, func(work string) {
			require.NoError(t, os.MkdirAll(filepath.Join(work, "pkg"), 0755))
			require.NoError(t, os.Symlink("../../../secret.txt", filepath.Join(work, "pkg", "rel.yaml")))
		})
		dst := filepath.Join(t.TempDir(), "fetched")
		err := Run([]string{repoDir + ".git/pkg@master", dst}, Options{AllowLocalRepo: true})
		require.ErrorIs(t, err, fsutil.ErrSymlink)
		_, statErr := os.Lstat(dst)
		assert.ErrorIs(t, statErr, os.ErrNotExist)
	})

	t.Run("cross-subdirectory in-repo file symlink", func(t *testing.T) {
		repoDir := initTestGitRepo(t, func(work string) {
			require.NoError(t, os.MkdirAll(filepath.Join(work, "private"), 0755))
			require.NoError(t, os.WriteFile(filepath.Join(work, "private", "creds.yaml"), []byte("secret: internal\n"), 0644))
			require.NoError(t, os.MkdirAll(filepath.Join(work, "pkg"), 0755))
			require.NoError(t, os.Symlink("../private/creds.yaml", filepath.Join(work, "pkg", "creds.yaml")))
		})
		dst := filepath.Join(t.TempDir(), "fetched")
		err := Run([]string{repoDir + ".git/pkg@master", dst}, Options{AllowLocalRepo: true})
		require.ErrorIs(t, err, fsutil.ErrSymlink)
		_, statErr := os.Lstat(dst)
		assert.ErrorIs(t, statErr, os.ErrNotExist)
	})

	t.Run("multi-hop chained symlink", func(t *testing.T) {
		repoDir := initTestGitRepo(t, func(work string) {
			require.NoError(t, os.MkdirAll(filepath.Join(work, "pkg"), 0755))
			require.NoError(t, os.Symlink(canaryFile, filepath.Join(work, "pkg", "hop2.yaml")))
			require.NoError(t, os.Symlink("hop2.yaml", filepath.Join(work, "pkg", "hop1.yaml")))
		})
		dst := filepath.Join(t.TempDir(), "fetched")
		err := Run([]string{repoDir + ".git/pkg@master", dst}, Options{AllowLocalRepo: true})
		require.ErrorIs(t, err, fsutil.ErrSymlink)
		_, statErr := os.Lstat(dst)
		assert.ErrorIs(t, statErr, os.ErrNotExist)
	})

	t.Run("leaf directory symlink", func(t *testing.T) {
		repoDir := initTestGitRepo(t, func(work string) {
			require.NoError(t, os.MkdirAll(filepath.Join(work, "pkg"), 0755))
			require.NoError(t, os.Symlink(nestedCanaryDir, filepath.Join(work, "pkg", "sub")))
		})
		dst := filepath.Join(t.TempDir(), "fetched")
		err := Run([]string{repoDir + ".git/pkg@master", dst}, Options{AllowLocalRepo: true})
		require.ErrorIs(t, err, fsutil.ErrSymlink)
		_, statErr := os.Lstat(dst)
		assert.ErrorIs(t, statErr, os.ErrNotExist)
	})

	t.Run("ancestor directory symlink in requested package path (Case B)", func(t *testing.T) {
		repoDir := initTestGitRepo(t, func(work string) {
			require.NoError(t, os.Symlink(canaryDir, filepath.Join(work, "escape")))
		})
		dst := filepath.Join(t.TempDir(), "fetched")
		err := Run([]string{repoDir + ".git/escape/secretdir@master", dst}, Options{AllowLocalRepo: true})
		require.ErrorIs(t, err, fsutil.ErrSymlink)
		_, statErr := os.Lstat(dst)
		assert.ErrorIs(t, statErr, os.ErrNotExist)
	})

	t.Run("symlink with spaces and quotes in filename", func(t *testing.T) {
		repoDir := initTestGitRepo(t, func(work string) {
			require.NoError(t, os.MkdirAll(filepath.Join(work, "pkg"), 0755))
			require.NoError(t, os.Symlink(canaryFile, filepath.Join(work, "pkg", `spaced "quoted" link.yaml`)))
		})
		dst := filepath.Join(t.TempDir(), "fetched")
		err := Run([]string{repoDir + ".git/pkg@master", dst}, Options{AllowLocalRepo: true})
		require.ErrorIs(t, err, fsutil.ErrSymlink)
		_, statErr := os.Lstat(dst)
		assert.ErrorIs(t, statErr, os.ErrNotExist)
	})

	t.Run("git submodule containing mode-120000 symlink", func(t *testing.T) {
		subRepo := initTestGitRepo(t, func(work string) {
			require.NoError(t, os.Symlink(canaryFile, filepath.Join(work, "sub_leak.yaml")))
		})
		parentRepo := initTestGitRepoWithSubmodule(t, subRepo, "pkg/submod")
		dst := filepath.Join(t.TempDir(), "fetched")
		err := Run([]string{parentRepo + ".git/pkg@master", dst}, Options{AllowLocalRepo: true})
		require.ErrorIs(t, err, fsutil.ErrSymlink)
		_, statErr := os.Lstat(dst)
		assert.ErrorIs(t, statErr, os.ErrNotExist)
	})

	t.Run("malicious .gitmodules URL rejected before submodule update", func(t *testing.T) {
		repoDir := initTestGitRepo(t, func(work string) {
			require.NoError(t, os.MkdirAll(filepath.Join(work, "pkg"), 0755))
			require.NoError(t, os.WriteFile(filepath.Join(work, "pkg", "cm.yaml"), []byte("kind: ConfigMap\n"), 0644))
			evilGitmodules := "[submodule \"evil\"]\n\tpath = pkg/evil\n\turl = ssh://-oProxyCommand=id/org/repo\n"
			require.NoError(t, os.WriteFile(filepath.Join(work, ".gitmodules"), []byte(evilGitmodules), 0644))
		})
		dst := filepath.Join(t.TempDir(), "fetched")
		err := Run([]string{repoDir + ".git/pkg@master", dst}, Options{AllowLocalRepo: true})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unsafe submodule URL in .gitmodules")
		_, statErr := os.Lstat(dst)
		assert.ErrorIs(t, statErr, os.ErrNotExist)
	})
}

func TestRunRejectsDestinationSymlinks(t *testing.T) {
	repoDir := initTestGitRepo(t, func(work string) {
		require.NoError(t, os.MkdirAll(filepath.Join(work, "pkg"), 0755))
		require.NoError(t, os.WriteFile(filepath.Join(work, "pkg", "cm.yaml"), []byte("apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: demo\n"), 0644))
	})

	t.Run("destination is a symlink to existing directory", func(t *testing.T) {
		victimDir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(victimDir, "keep.txt"), []byte("intact"), 0600))
		linkDst := filepath.Join(t.TempDir(), "dst-link")
		require.NoError(t, os.Symlink(victimDir, linkDst))

		for _, clean := range []bool{false, true} {
			err := Run([]string{repoDir + ".git/pkg@master", linkDst}, Options{AllowLocalRepo: true, Clean: clean})
			require.ErrorIs(t, err, fsutil.ErrSymlink)
			got, readErr := os.ReadFile(filepath.Join(victimDir, "keep.txt"))
			require.NoError(t, readErr)
			assert.Equal(t, "intact", string(got))
		}
	})

	t.Run("existing destination contains child symlink with Clean true", func(t *testing.T) {
		victimFile := filepath.Join(t.TempDir(), "victim.txt")
		require.NoError(t, os.WriteFile(victimFile, []byte("intact"), 0600))

		dst := filepath.Join(t.TempDir(), "existing-dst")
		require.NoError(t, os.MkdirAll(dst, 0755))
		require.NoError(t, os.Symlink(victimFile, filepath.Join(dst, "child-link")))

		err := Run([]string{repoDir + ".git/pkg@master", dst}, Options{AllowLocalRepo: true, Clean: true})
		require.ErrorIs(t, err, fsutil.ErrSymlink)

		got, readErr := os.ReadFile(victimFile)
		require.NoError(t, readErr)
		assert.Equal(t, "intact", string(got))
	})

	t.Run("destination parent is a symlink", func(t *testing.T) {
		victimDir := t.TempDir()
		parentLink := filepath.Join(t.TempDir(), "parent-link")
		require.NoError(t, os.Symlink(victimDir, parentLink))

		err := Run([]string{repoDir + ".git/pkg@master", filepath.Join(parentLink, "newpkg")}, Options{AllowLocalRepo: true})
		require.ErrorIs(t, err, fsutil.ErrSymlink)
	})
}

func TestRunBenignGitPackageWithKptfileAndAutoSet(t *testing.T) {
	kptfileContent := `apiVersion: kpt.dev/v1alpha1
kind: Kptfile
metadata:
  name: singlevm
openAPI:
  definitions:
    io.k8s.cli.setters.image:
      x-k8s-cli:
        setter:
          name: image
          value: wordpress-default
`
	configContent := `apiVersion: dev.marketplace.cloud.google.com/v1alpha1
kind: DeploymentManagerTemplate
metadata:
  name: wordpress
spec:
  image: wordpress-default # {"$kpt-set":"image"}
`
	var headSHA string
	repoDir := initTestGitRepo(t, func(work string) {
		pkgDir := filepath.Join(work, "examples", "singlevm")
		require.NoError(t, os.MkdirAll(filepath.Join(pkgDir, "nested"), 0755))
		require.NoError(t, os.WriteFile(filepath.Join(pkgDir, "Kptfile"), []byte(kptfileContent), 0644))
		require.NoError(t, os.WriteFile(filepath.Join(pkgDir, "configurations.yaml"), []byte(configContent), 0644))
		require.NoError(t, os.WriteFile(filepath.Join(pkgDir, ".gitignore"), []byte("template.zip\n"), 0644))
		require.NoError(t, os.WriteFile(filepath.Join(pkgDir, "nested", "README.md"), []byte("# Single VM\n"), 0644))
	})
	headBytes, err := exec.Command("git", "-C", repoDir, "rev-parse", "--verify", "HEAD").Output()
	require.NoError(t, err)
	headSHA = strings.TrimSpace(string(headBytes))

	// Tag with subdirectory prefix to verify kpt tag resolution compatibility.
	runGitCmd(t, repoDir, "tag", "examples/singlevm/v1.0.0")

	t.Setenv("KPT_SET_IMAGE", "wordpress-hardened-v2")

	parentOut := t.TempDir()
	var stdout bytes.Buffer
	err = Run(
		[]string{repoDir + ".git/examples/singlevm@v1.0.0", parentOut},
		Options{
			AllowLocalRepo: true,
			AutoSet:        true,
			Stdout:         &stdout,
		},
	)
	require.NoError(t, err)

	fetchedDir := filepath.Join(parentOut, "singlevm")
	assert.Contains(t, stdout.String(), "fetching package /examples/singlevm")

	gotCfg, err := os.ReadFile(filepath.Join(fetchedDir, "configurations.yaml"))
	require.NoError(t, err)
	assert.Contains(t, string(gotCfg), "image: wordpress-hardened-v2")

	gotGitignore, err := os.ReadFile(filepath.Join(fetchedDir, ".gitignore"))
	require.NoError(t, err)
	assert.Equal(t, "template.zip\n", string(gotGitignore))

	gotReadme, err := os.ReadFile(filepath.Join(fetchedDir, "nested", "README.md"))
	require.NoError(t, err)
	assert.Equal(t, "# Single VM\n", string(gotReadme))

	gotKptfile, err := os.ReadFile(filepath.Join(fetchedDir, "Kptfile"))
	require.NoError(t, err)
	assert.Contains(t, string(gotKptfile), "commit: "+headSHA)
	assert.Contains(t, string(gotKptfile), "directory: /examples/singlevm")
	assert.Contains(t, string(gotKptfile), "ref: v1.0.0")
	assert.Contains(t, string(gotKptfile), "io.k8s.cli.setters.image")
}

func TestRunStdinMode(t *testing.T) {
	validStream := `apiVersion: v1
kind: ConfigMap
metadata:
  name: demo
data:
  key: value
`

	t.Run("valid stdin stream writes resources and Kptfile", func(t *testing.T) {
		dst := filepath.Join(t.TempDir(), "stdin-pkg")
		err := Run([]string{"-", dst}, Options{
			Stdin: strings.NewReader(validStream),
		})
		require.NoError(t, err)

		cmBytes, err := os.ReadFile(filepath.Join(dst, "demo_configmap.yaml"))
		require.NoError(t, err)
		assert.Contains(t, string(cmBytes), "name: demo")

		kfBytes, err := os.ReadFile(filepath.Join(dst, "Kptfile"))
		require.NoError(t, err)
		assert.Contains(t, string(kfBytes), "type: stdin")
	})

	t.Run("rejects parent traversal in filename pattern", func(t *testing.T) {
		dst := filepath.Join(t.TempDir(), "stdin-pkg")
		err := Run([]string{"-", dst}, Options{
			Pattern: "../escape_%n.yaml",
			Stdin:   strings.NewReader(validStream),
		})
		require.ErrorIs(t, err, fsutil.ErrPathEscapesRoot)
		_, statErr := os.Lstat(dst)
		assert.ErrorIs(t, statErr, os.ErrNotExist)
	})

	t.Run("rejects traversal in resource path annotation", func(t *testing.T) {
		poisoned := `apiVersion: v1
kind: ConfigMap
metadata:
  name: demo
  annotations:
    config.kubernetes.io/path: "../../outside.yaml"
`
		dst := filepath.Join(t.TempDir(), "stdin-pkg")
		err := Run([]string{"-", dst}, Options{
			Stdin: strings.NewReader(poisoned),
		})
		require.ErrorIs(t, err, fsutil.ErrPathEscapesRoot)
		_, statErr := os.Lstat(dst)
		assert.ErrorIs(t, statErr, os.ErrNotExist)
	})

	t.Run("rejects traversal in metadata.name expanded by pattern", func(t *testing.T) {
		poisonedName := `apiVersion: v1
kind: ConfigMap
metadata:
  name: "../../outside"
`
		dst := filepath.Join(t.TempDir(), "stdin-pkg")
		err := Run([]string{"-", dst}, Options{
			Stdin: strings.NewReader(poisonedName),
		})
		require.ErrorIs(t, err, fsutil.ErrPathEscapesRoot)
		_, statErr := os.Lstat(dst)
		assert.ErrorIs(t, statErr, os.ErrNotExist)
	})

	t.Run("rejects symlink destination", func(t *testing.T) {
		victim := t.TempDir()
		linkDst := filepath.Join(t.TempDir(), "link-dst")
		require.NoError(t, os.Symlink(victim, linkDst))

		err := Run([]string{"-", linkDst}, Options{
			Stdin: strings.NewReader(validStream),
		})
		require.ErrorIs(t, err, fsutil.ErrSymlink)
	})
}

func initTestGitRepo(t *testing.T, populate func(workDir string)) string {
	t.Helper()
	repoDir := filepath.Join(t.TempDir(), "repo")
	require.NoError(t, os.MkdirAll(repoDir, 0755))

	runGitCmd(t, repoDir, "init", "-b", "master")
	runGitCmd(t, repoDir, "config", "user.email", "test@example.com")
	runGitCmd(t, repoDir, "config", "user.name", "Test User")
	populate(repoDir)
	runGitCmd(t, repoDir, "add", "-A")
	runGitCmd(t, repoDir, "commit", "-m", "initial commit")
	return repoDir
}

func initTestGitRepoWithSubmodule(t *testing.T, subRepo, subPath string) string {
	t.Helper()
	repoDir := filepath.Join(t.TempDir(), "parent-repo")
	require.NoError(t, os.MkdirAll(repoDir, 0755))

	runGitCmd(t, repoDir, "init", "-b", "master")
	runGitCmd(t, repoDir, "config", "user.email", "test@example.com")
	runGitCmd(t, repoDir, "config", "user.name", "Test User")
	require.NoError(t, os.MkdirAll(filepath.Dir(filepath.Join(repoDir, subPath)), 0755))
	runGitCmd(t, repoDir, "-c", "protocol.file.allow=always", "submodule", "add", subRepo, subPath)
	runGitCmd(t, repoDir, "commit", "-m", "add submodule")
	return repoDir
}

func runGitCmd(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil && len(args) >= 2 && args[0] == "init" && args[1] == "-b" {
		// Fallback for older Git versions without `git init -b`.
		cmd2 := exec.Command("git", "init")
		cmd2.Dir = dir
		out2, err2 := cmd2.CombinedOutput()
		require.NoError(t, err2, "git init failed: %s", string(out2))
		cmd3 := exec.Command("git", "symbolic-ref", "HEAD", "refs/heads/master")
		cmd3.Dir = dir
		out3, err3 := cmd3.CombinedOutput()
		require.NoError(t, err3, "git symbolic-ref failed: %s", string(out3))
		return
	}
	require.NoError(t, err, "git %v failed: %s", args, string(out))
}
