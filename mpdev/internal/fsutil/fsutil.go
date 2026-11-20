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

// Package fsutil provides symlink-rejecting, root-contained filesystem
// operations compatible with Go 1.23.
package fsutil

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

var (
	// ErrSymlink indicates that a symbolic link was encountered where symlinks
	// are forbidden.
	ErrSymlink = errors.New("symlink not permitted")

	// ErrSpecialFile indicates that a non-regular, non-directory file (such as
	// a FIFO, socket, or device node) was encountered.
	ErrSpecialFile = errors.New("non-regular file not permitted")

	// ErrPathEscapesRoot indicates that a path is absolute, malformed, or
	// escapes its root directory.
	ErrPathEscapesRoot = errors.New("path escapes root directory")
)

// CanonicalDir resolves any OS-level parent symlinks on dir (such as /var ->
// /private/var on macOS temp directories) and verifies that the target exists
// and is a directory. If dir itself is a leaf symlink, CanonicalDir still
// resolves it; callers that must reject a leaf symlink on a user-supplied
// destination path should call RejectLeafSymlink first.
func CanonicalDir(dir string) (string, error) {
	if dir == "" {
		return "", fmt.Errorf("empty directory path: %w", ErrPathEscapesRoot)
	}
	if strings.IndexByte(dir, 0) >= 0 {
		return "", fmt.Errorf("NUL byte in directory path: %w", ErrPathEscapesRoot)
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("resolving absolute path for %q: %w", dir, err)
	}
	realDir, err := filepath.EvalSymlinks(absDir)
	if err != nil {
		return "", fmt.Errorf("evaluating symlinks for %q: %w", dir, err)
	}
	realDir = filepath.Clean(realDir)
	info, err := os.Lstat(realDir)
	if err != nil {
		return "", fmt.Errorf("stat canonical directory %q: %w", realDir, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return "", fmt.Errorf("canonical root %q: %w", realDir, ErrSymlink)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("path %q is not a directory", dir)
	}
	return realDir, nil
}

// RejectLeafSymlink checks path via os.Lstat (without resolving parent
// symlinks) and returns ErrSymlink if the final path element is a symlink. It
// returns nil if path does not exist.
func RejectLeafSymlink(path string) error {
	if path == "" || strings.IndexByte(path, 0) >= 0 {
		return fmt.Errorf("invalid path %q: %w", path, ErrPathEscapesRoot)
	}
	// Trim trailing slashes so os.Lstat inspects the link itself rather than
	// dereferencing a directory symlink.
	trimmed := strings.TrimRight(path, "/\\")
	if trimmed == "" {
		trimmed = path
	}
	info, err := os.Lstat(trimmed)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("path %q: %w", path, ErrSymlink)
	}
	return nil
}

// RejectSymlinksInTree walks realRoot using filepath.WalkDir and fails closed
// if any entry is a symlink or a non-regular, non-directory file. When
// skipDotGit is true, a real .git directory or submodule .git regular file is
// skipped (while .gitignore, .github, etc. are still inspected), whereas a
// .git symlink is still rejected.
func RejectSymlinksInTree(realRoot string, skipDotGit bool) error {
	if err := verifyCanonicalRoot(realRoot); err != nil {
		return err
	}
	return filepath.WalkDir(realRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(realRoot, path)
		if err != nil {
			return err
		}
		if d.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("entry %q: %w", rel, ErrSymlink)
		}
		relSlash := filepath.ToSlash(rel)
		if skipDotGit && (d.Name() == ".git" || relSlash == ".git" || strings.HasPrefix(relSlash, ".git/")) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			if d.Type().IsRegular() {
				return nil
			}
			return fmt.Errorf("entry %q (mode %s): %w", rel, d.Type(), ErrSpecialFile)
		}
		if !d.IsDir() && !d.Type().IsRegular() {
			return fmt.Errorf("entry %q (mode %s): %w", rel, d.Type(), ErrSpecialFile)
		}
		return nil
	})
}

// ValidateContainedRelPath verifies that relPath is a clean relative path
// strictly inside realRoot (no absolute path, NUL bytes, backslashes, or ".."
// segments) and inspects every existing path segment from realRoot down to the
// target with os.Lstat, rejecting symlinks and special files.
func ValidateContainedRelPath(realRoot, relPath string) (string, error) {
	if err := verifyCanonicalRoot(realRoot); err != nil {
		return "", err
	}
	if relPath == "" {
		return "", fmt.Errorf("empty relative path: %w", ErrPathEscapesRoot)
	}
	if strings.IndexByte(relPath, 0) >= 0 {
		return "", fmt.Errorf("NUL byte in path: %w", ErrPathEscapesRoot)
	}
	if strings.ContainsRune(relPath, '\\') {
		return "", fmt.Errorf("backslash not permitted in relative path %q: %w", relPath, ErrPathEscapesRoot)
	}
	if filepath.IsAbs(relPath) || strings.HasPrefix(relPath, "/") {
		return "", fmt.Errorf("absolute path %q: %w", relPath, ErrPathEscapesRoot)
	}

	rawSlash := filepath.ToSlash(relPath)
	for _, rawSeg := range strings.Split(rawSlash, "/") {
		if rawSeg == ".." {
			return "", fmt.Errorf("parent traversal segment in %q: %w", relPath, ErrPathEscapesRoot)
		}
	}

	cleanRel := filepath.Clean(relPath)
	if cleanRel == "." {
		return realRoot, nil
	}
	if cleanRel == ".." || strings.HasPrefix(cleanRel, ".."+string(filepath.Separator)) || filepath.IsAbs(cleanRel) {
		return "", fmt.Errorf("path %q: %w", relPath, ErrPathEscapesRoot)
	}

	segments := strings.Split(cleanRel, string(filepath.Separator))
	cur := realRoot
	missingAncestor := false
	for i, seg := range segments {
		if seg == "" || seg == "." || seg == ".." {
			return "", fmt.Errorf("invalid segment %q in %q: %w", seg, relPath, ErrPathEscapesRoot)
		}
		cur = filepath.Join(cur, seg)
		if missingAncestor {
			continue
		}
		info, err := os.Lstat(cur)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				missingAncestor = true
				continue
			}
			return "", err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("path %q (at %q): %w", relPath, seg, ErrSymlink)
		}
		isLast := i == len(segments)-1
		if !isLast {
			if !info.IsDir() {
				return "", fmt.Errorf("ancestor %q in %q is not a directory", seg, relPath)
			}
			continue
		}
		if !info.IsDir() && !info.Mode().IsRegular() {
			return "", fmt.Errorf("path %q (mode %s): %w", relPath, info.Mode(), ErrSpecialFile)
		}
	}
	return cur, nil
}

// OpenRegularFileNoFollow validates relPath within realRoot and opens the
// target with O_NOFOLLOW, verifying on the open file descriptor that it is a
// regular file contained within realRoot.
func OpenRegularFileNoFollow(realRoot, relPath string, flags int, perm os.FileMode) (*os.File, os.FileInfo, error) {
	target, err := ValidateContainedRelPath(realRoot, relPath)
	if err != nil {
		return nil, nil, err
	}
	if target == realRoot {
		return nil, nil, fmt.Errorf("path %q refers to root directory, want regular file", relPath)
	}

	f, err := os.OpenFile(target, flags|syscall.O_NOFOLLOW, perm)
	if err != nil {
		if isNoFollowSymlinkErr(err) {
			return nil, nil, fmt.Errorf("opening %q: %w", relPath, ErrSymlink)
		}
		return nil, nil, err
	}

	info, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return nil, nil, err
	}
	if !info.Mode().IsRegular() {
		_ = f.Close()
		return nil, nil, fmt.Errorf("file %q (mode %s): %w", relPath, info.Mode(), ErrSpecialFile)
	}
	if err := verifyOpenFDContained(f, realRoot); err != nil {
		_ = f.Close()
		return nil, nil, err
	}
	return f, info, nil
}

// ReadRegularFileNoFollow reads a regular file inside realRoot without
// following symlinks on any path component or leaf.
func ReadRegularFileNoFollow(realRoot, relPath string) ([]byte, os.FileMode, error) {
	f, info, err := OpenRegularFileNoFollow(realRoot, relPath, os.O_RDONLY, 0)
	if err != nil {
		return nil, 0, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, 0, err
	}
	return data, info.Mode(), nil
}

// OverwriteRegularFileNoFollow overwrites an existing regular file inside
// realRoot without following symlinks. It intentionally omits O_CREAT and
// O_TRUNC at open time so that regular-file and FD-containment checks run
// before truncating a single byte.
func OverwriteRegularFileNoFollow(realRoot, relPath string, data []byte, defaultPerm os.FileMode) error {
	f, _, err := OpenRegularFileNoFollow(realRoot, relPath, os.O_WRONLY, defaultPerm)
	if err != nil {
		return err
	}
	closed := false
	defer func() {
		if !closed {
			_ = f.Close()
		}
	}()

	if err := f.Truncate(0); err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		return err
	}
	closed = true
	return f.Close()
}

// WriteRegularFileNoFollow creates or overwrites a regular file inside
// realRoot without following symlinks.
func WriteRegularFileNoFollow(realRoot, relPath string, data []byte, perm os.FileMode) error {
	f, _, err := OpenRegularFileNoFollow(realRoot, relPath, os.O_WRONLY|os.O_CREATE, perm)
	if err != nil {
		return err
	}
	closed := false
	defer func() {
		if !closed {
			_ = f.Close()
		}
	}()

	if err := f.Truncate(0); err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		return err
	}
	closed = true
	return f.Close()
}

// CopyDirNoSymlinks copies subDir within srcCloneRoot to dstDir while
// rejecting symlinks and non-regular files across ancestor components, the
// source subtree, and any pre-existing destination directory.
func CopyDirNoSymlinks(srcCloneRoot, subDir, dstDir string, cleanDst bool) error {
	return StageAndCommitCopyDir(srcCloneRoot, subDir, dstDir, cleanDst, nil)
}

// StageAndCommitCopyDir copies subDir within srcCloneRoot into a sibling
// staging directory next to dstDir, runs an optional postStage hook on the
// staged tree, re-verifies no symlinks were introduced, and renames the staged
// directory to dstDir. On any error, the staging directory is removed and
// dstDir is left untouched.
func StageAndCommitCopyDir(srcCloneRoot, subDir, dstDir string, cleanDst bool, postStage func(stageDir string) error) error {
	if dstDir == "" || strings.IndexByte(dstDir, 0) >= 0 {
		return fmt.Errorf("invalid destination %q: %w", dstDir, ErrPathEscapesRoot)
	}
	if err := RejectLeafSymlink(dstDir); err != nil {
		return err
	}

	dstAbs, err := filepath.Abs(dstDir)
	if err != nil {
		return fmt.Errorf("resolving destination %q: %w", dstDir, err)
	}
	dstBase := filepath.Base(dstAbs)
	if dstBase == "." || dstBase == ".." || dstBase == string(filepath.Separator) {
		return fmt.Errorf("invalid destination directory %q: %w", dstDir, ErrPathEscapesRoot)
	}

	parentReal, err := CanonicalDir(filepath.Dir(dstAbs))
	if err != nil {
		return fmt.Errorf("destination parent directory: %w", err)
	}
	dstTarget := filepath.Join(parentReal, dstBase)

	if err := checkExistingDestination(dstTarget, cleanDst); err != nil {
		return err
	}

	realCloneRoot, err := CanonicalDir(srcCloneRoot)
	if err != nil {
		return fmt.Errorf("source clone root: %w", err)
	}

	relSub, err := normalizeSubDir(subDir)
	if err != nil {
		return err
	}
	srcSubPath, err := ValidateContainedRelPath(realCloneRoot, relSub)
	if err != nil {
		return fmt.Errorf("validating package subdirectory %q: %w", subDir, err)
	}
	srcSubInfo, err := os.Lstat(srcSubPath)
	if err != nil {
		return err
	}
	if srcSubInfo.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("package subdirectory %q: %w", subDir, ErrSymlink)
	}
	if !srcSubInfo.IsDir() {
		return fmt.Errorf("package subdirectory %q is not a directory", subDir)
	}

	if err := RejectSymlinksInTree(srcSubPath, true); err != nil {
		return err
	}

	stageDir, err := os.MkdirTemp(parentReal, ".mpdev-pkg-stage-*")
	if err != nil {
		return fmt.Errorf("creating staging directory: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			_ = os.RemoveAll(stageDir)
		}
	}()
	if err := os.Chmod(stageDir, 0755); err != nil {
		return err
	}

	err = filepath.WalkDir(srcSubPath, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(srcSubPath, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("entry %q: %w", rel, ErrSymlink)
		}
		relSlash := filepath.ToSlash(rel)
		if d.Name() == ".git" || relSlash == ".git" || strings.HasPrefix(relSlash, ".git/") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			if d.Type().IsRegular() {
				return nil
			}
			return fmt.Errorf("entry %q: %w", rel, ErrSpecialFile)
		}
		if d.IsDir() {
			dirInfo, err := d.Info()
			if err != nil {
				return err
			}
			if dirInfo.Mode()&os.ModeSymlink != 0 {
				return fmt.Errorf("directory %q: %w", rel, ErrSymlink)
			}
			// Ensure owner rwx while staging nested files.
			return os.MkdirAll(filepath.Join(stageDir, rel), dirInfo.Mode().Perm()|0700)
		}
		if !d.Type().IsRegular() {
			return fmt.Errorf("entry %q (mode %s): %w", rel, d.Type(), ErrSpecialFile)
		}

		data, mode, err := ReadRegularFileNoFollow(srcSubPath, rel)
		if err != nil {
			return err
		}
		dstFile, err := ValidateContainedRelPath(stageDir, rel)
		if err != nil {
			return err
		}
		out, err := os.OpenFile(dstFile, os.O_WRONLY|os.O_CREATE|os.O_EXCL|syscall.O_NOFOLLOW, mode.Perm())
		if err != nil {
			if isNoFollowSymlinkErr(err) {
				return fmt.Errorf("creating staged file %q: %w", rel, ErrSymlink)
			}
			return err
		}
		_, writeErr := out.Write(data)
		closeErr := out.Close()
		if writeErr != nil {
			return writeErr
		}
		return closeErr
	})
	if err != nil {
		return err
	}

	if postStage != nil {
		if err := postStage(stageDir); err != nil {
			return err
		}
	}
	if err := RejectSymlinksInTree(stageDir, false); err != nil {
		return err
	}

	if err := checkExistingDestination(dstTarget, cleanDst); err != nil {
		return err
	}
	if cleanDst {
		if _, err := os.Lstat(dstTarget); err == nil { // if NO error
			if err := os.RemoveAll(dstTarget); err != nil {
				return fmt.Errorf("cleaning existing destination %q: %w", dstTarget, err)
			}
		}
	}

	if err := os.Rename(stageDir, dstTarget); err != nil {
		return fmt.Errorf("committing staged directory to %q: %w", dstTarget, err)
	}
	committed = true
	return nil
}

func checkExistingDestination(dstTarget string, cleanDst bool) error {
	info, err := os.Lstat(dstTarget)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("destination %q: %w", dstTarget, ErrSymlink)
	}
	if !cleanDst {
		return fmt.Errorf("destination directory %s already exists", dstTarget)
	}
	if !info.IsDir() {
		return fmt.Errorf("destination %q exists and is not a directory", dstTarget)
	}
	return RejectSymlinksInTree(dstTarget, false)
}

func normalizeSubDir(subDir string) (string, error) {
	if strings.IndexByte(subDir, 0) >= 0 {
		return "", fmt.Errorf("NUL byte in subdirectory: %w", ErrPathEscapesRoot)
	}
	if strings.ContainsRune(subDir, '\\') {
		return "", fmt.Errorf("backslash not permitted in subdirectory %q: %w", subDir, ErrPathEscapesRoot)
	}
	trimmed := strings.TrimPrefix(filepath.ToSlash(subDir), "/")
	if trimmed == "" || trimmed == "." {
		return ".", nil
	}
	for _, seg := range strings.Split(trimmed, "/") {
		if seg == ".." {
			return "", fmt.Errorf("parent traversal segment in subdirectory %q: %w", subDir, ErrPathEscapesRoot)
		}
	}
	return filepath.FromSlash(trimmed), nil
}

func verifyCanonicalRoot(realRoot string) error {
	if realRoot == "" || !filepath.IsAbs(realRoot) || strings.IndexByte(realRoot, 0) >= 0 {
		return fmt.Errorf("root %q must be a non-empty canonical absolute path: %w", realRoot, ErrPathEscapesRoot)
	}
	info, err := os.Lstat(realRoot)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("root %q: %w", realRoot, ErrSymlink)
	}
	if !info.IsDir() {
		return fmt.Errorf("root %q is not a directory", realRoot)
	}
	return nil
}

func isNoFollowSymlinkErr(err error) bool {
	return errors.Is(err, syscall.ELOOP) || errors.Is(err, syscall.EMLINK)
}

func verifyOpenFDContained(f *os.File, realRoot string) error {
	if runtime.GOOS != "linux" {
		return nil
	}
	fdPath, err := os.Readlink(fmt.Sprintf("/proc/self/fd/%d", f.Fd()))
	if err != nil {
		return nil
	}
	prefix := realRoot
	if !strings.HasSuffix(prefix, string(filepath.Separator)) {
		prefix += string(filepath.Separator)
	}
	if !strings.HasPrefix(fdPath, prefix) {
		return fmt.Errorf("open file descriptor resolved outside root (%q): %w", fdPath, ErrPathEscapesRoot)
	}
	return nil
}
