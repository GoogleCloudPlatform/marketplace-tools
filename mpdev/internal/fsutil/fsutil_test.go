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

package fsutil

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

func TestCanonicalDirAllowsParentSymlink(t *testing.T) {
	base := t.TempDir()
	realSub := filepath.Join(base, "real_parent", "pkg")
	if err := os.MkdirAll(realSub, 0755); err != nil {
		t.Fatalf("os.MkdirAll(%q) failed: %v", realSub, err)
	}
	linkParent := filepath.Join(base, "symlink_parent")
	if err := os.Symlink(filepath.Join(base, "real_parent"), linkParent); err != nil {
		t.Fatalf("os.Symlink failed: %v", err)
	}

	viaLink := filepath.Join(linkParent, "pkg")
	got, err := CanonicalDir(viaLink)
	if err != nil {
		t.Fatalf("CanonicalDir(%q) unexpected error: %v", viaLink, err)
	}
	want, err := filepath.EvalSymlinks(realSub)
	if err != nil {
		t.Fatalf("filepath.EvalSymlinks(%q) failed: %v", realSub, err)
	}
	if got != want {
		t.Errorf("CanonicalDir(%q) = %q, want %q", viaLink, got, want)
	}
}

func TestRejectSymlinksInTree(t *testing.T) {
	t.Run("clean_tree_and_dot_git_skipped_while_gitignore_kept", func(t *testing.T) {
		root, err := CanonicalDir(t.TempDir())
		if err != nil {
			t.Fatalf("CanonicalDir failed: %v", err)
		}
		if err := os.MkdirAll(filepath.Join(root, ".git", "hooks"), 0755); err != nil {
			t.Fatalf("MkdirAll .git failed: %v", err)
		}
		if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte("*.tmp\n"), 0644); err != nil {
			t.Fatalf("WriteFile .gitignore failed: %v", err)
		}
		if err := os.WriteFile(filepath.Join(root, "main.tf"), []byte("variable \"x\" {}\n"), 0644); err != nil {
			t.Fatalf("WriteFile main.tf failed: %v", err)
		}
		if err := RejectSymlinksInTree(root, true); err != nil {
			t.Errorf("RejectSymlinksInTree(%q, true) = %v, want nil", root, err)
		}
	})

	t.Run("rejects_file_symlink_and_dangling_symlink", func(t *testing.T) {
		root, err := CanonicalDir(t.TempDir())
		if err != nil {
			t.Fatalf("CanonicalDir failed: %v", err)
		}
		if err := os.Symlink("/nonexistent/target", filepath.Join(root, "dangling.yaml")); err != nil {
			t.Fatalf("os.Symlink failed: %v", err)
		}
		if err := RejectSymlinksInTree(root, true); !errors.Is(err, ErrSymlink) {
			t.Errorf("RejectSymlinksInTree(%q, true) = %v, want ErrSymlink", root, err)
		}
	})

	t.Run("rejects_gitignore_symlink_even_when_skip_dot_git_true", func(t *testing.T) {
		root, err := CanonicalDir(t.TempDir())
		if err != nil {
			t.Fatalf("CanonicalDir failed: %v", err)
		}
		outside := filepath.Join(t.TempDir(), "secret.txt")
		if err := os.WriteFile(outside, []byte("secret"), 0600); err != nil {
			t.Fatalf("WriteFile failed: %v", err)
		}
		if err := os.Symlink(outside, filepath.Join(root, ".gitignore")); err != nil {
			t.Fatalf("os.Symlink failed: %v", err)
		}
		if err := RejectSymlinksInTree(root, true); !errors.Is(err, ErrSymlink) {
			t.Errorf("RejectSymlinksInTree(%q, true) = %v, want ErrSymlink", root, err)
		}
	})

	t.Run("rejects_dot_git_symlink", func(t *testing.T) {
		root, err := CanonicalDir(t.TempDir())
		if err != nil {
			t.Fatalf("CanonicalDir failed: %v", err)
		}
		if err := os.Symlink(t.TempDir(), filepath.Join(root, ".git")); err != nil {
			t.Fatalf("os.Symlink failed: %v", err)
		}
		if err := RejectSymlinksInTree(root, true); !errors.Is(err, ErrSymlink) {
			t.Errorf("RejectSymlinksInTree(%q, true) = %v, want ErrSymlink", root, err)
		}
	})

	t.Run("rejects_named_pipe_special_file", func(t *testing.T) {
		root, err := CanonicalDir(t.TempDir())
		if err != nil {
			t.Fatalf("CanonicalDir failed: %v", err)
		}
		fifo := filepath.Join(root, "pipe.fifo")
		if err := syscall.Mkfifo(fifo, 0600); err != nil {
			t.Fatalf("syscall.Mkfifo(%q) failed: %v", fifo, err)
		}
		if err := RejectSymlinksInTree(root, true); !errors.Is(err, ErrSpecialFile) {
			t.Errorf("RejectSymlinksInTree(%q, true) = %v, want ErrSpecialFile", root, err)
		}
	})
}

func TestValidateContainedRelPath(t *testing.T) {
	root, err := CanonicalDir(t.TempDir())
	if err != nil {
		t.Fatalf("CanonicalDir failed: %v", err)
	}
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(root, "escape_dir")); err != nil {
		t.Fatalf("os.Symlink failed: %v", err)
	}

	tests := []struct {
		name    string
		relPath string
		wantErr error
	}{
		{name: "absolute_path", relPath: "/etc/passwd", wantErr: ErrPathEscapesRoot},
		{name: "dot_dot_root", relPath: "../outside", wantErr: ErrPathEscapesRoot},
		{name: "dot_dot_after_symlink", relPath: "escape_dir/../normal.txt", wantErr: ErrPathEscapesRoot},
		{name: "backslash_path", relPath: `sub\file.txt`, wantErr: ErrPathEscapesRoot},
		{name: "nul_byte", relPath: "file\x00.txt", wantErr: ErrPathEscapesRoot},
		{name: "ancestor_symlink", relPath: "escape_dir/sub/file.txt", wantErr: ErrSymlink},
		{name: "leaf_symlink", relPath: "escape_dir", wantErr: ErrSymlink},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ValidateContainedRelPath(root, tc.relPath)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("ValidateContainedRelPath(%q, %q) = %v, want %v", root, tc.relPath, err, tc.wantErr)
			}
		})
	}
}

func TestReadAndOverwriteRegularFileNoFollow(t *testing.T) {
	root, err := CanonicalDir(t.TempDir())
	if err != nil {
		t.Fatalf("CanonicalDir failed: %v", err)
	}
	canaryDir := t.TempDir()
	canaryFile := filepath.Join(canaryDir, "canary.txt")
	const canaryContent = "ORIGINAL_CANARY_DO_NOT_TOUCH"
	if err := os.WriteFile(canaryFile, []byte(canaryContent), 0600); err != nil {
		t.Fatalf("WriteFile canary failed: %v", err)
	}

	if err := os.Symlink(canaryFile, filepath.Join(root, "link.txt")); err != nil {
		t.Fatalf("os.Symlink failed: %v", err)
	}
	if _, _, err := ReadRegularFileNoFollow(root, "link.txt"); !errors.Is(err, ErrSymlink) {
		t.Errorf("ReadRegularFileNoFollow(link.txt) = %v, want ErrSymlink", err)
	}
	if err := OverwriteRegularFileNoFollow(root, "link.txt", []byte("OVERWRITTEN"), 0644); !errors.Is(err, ErrSymlink) {
		t.Errorf("OverwriteRegularFileNoFollow(link.txt) = %v, want ErrSymlink", err)
	}
	gotCanary, err := os.ReadFile(canaryFile)
	if err != nil {
		t.Fatalf("ReadFile canary failed: %v", err)
	}
	if string(gotCanary) != canaryContent {
		t.Errorf("canary modified: got %q, want %q", string(gotCanary), canaryContent)
	}

	regularPath := filepath.Join(root, "regular.txt")
	if err := os.WriteFile(regularPath, []byte("before"), 0640); err != nil {
		t.Fatalf("WriteFile regular failed: %v", err)
	}
	if err := OverwriteRegularFileNoFollow(root, "regular.txt", []byte("after"), 0640); err != nil {
		t.Fatalf("OverwriteRegularFileNoFollow(regular.txt) failed: %v", err)
	}
	gotData, _, err := ReadRegularFileNoFollow(root, "regular.txt")
	if err != nil {
		t.Fatalf("ReadRegularFileNoFollow(regular.txt) failed: %v", err)
	}
	if string(gotData) != "after" {
		t.Errorf("ReadRegularFileNoFollow(regular.txt) = %q, want %q", string(gotData), "after")
	}
}

func TestCopyDirNoSymlinks(t *testing.T) {
	t.Run("copies_clean_subdir_and_preserves_gitignore_excluding_dot_git", func(t *testing.T) {
		srcRoot := t.TempDir()
		if err := os.MkdirAll(filepath.Join(srcRoot, ".git"), 0755); err != nil {
			t.Fatalf("MkdirAll .git failed: %v", err)
		}
		if err := os.WriteFile(filepath.Join(srcRoot, ".git", "config"), []byte("gitcfg"), 0600); err != nil {
			t.Fatalf("WriteFile .git/config failed: %v", err)
		}
		pkgDir := filepath.Join(srcRoot, "examples", "pkg")
		if err := os.MkdirAll(filepath.Join(pkgDir, "nested"), 0755); err != nil {
			t.Fatalf("MkdirAll nested failed: %v", err)
		}
		if err := os.WriteFile(filepath.Join(pkgDir, ".gitignore"), []byte("*.bak\n"), 0644); err != nil {
			t.Fatalf("WriteFile .gitignore failed: %v", err)
		}
		if err := os.WriteFile(filepath.Join(pkgDir, "nested", "app.yaml"), []byte("kind: Config\n"), 0644); err != nil {
			t.Fatalf("WriteFile app.yaml failed: %v", err)
		}

		dstDir := filepath.Join(t.TempDir(), "out")
		if err := CopyDirNoSymlinks(srcRoot, "examples/pkg", dstDir, false); err != nil {
			t.Fatalf("CopyDirNoSymlinks failed: %v", err)
		}
		if _, err := os.Lstat(filepath.Join(dstDir, ".gitignore")); err != nil {
			t.Errorf("expected .gitignore in dstDir, got err: %v", err)
		}
		if _, err := os.Lstat(filepath.Join(dstDir, "nested", "app.yaml")); err != nil {
			t.Errorf("expected nested/app.yaml in dstDir, got err: %v", err)
		}
	})

	t.Run("rejects_ancestor_directory_symlink_without_creating_dest", func(t *testing.T) {
		secretRoot := t.TempDir()
		secretSub := filepath.Join(secretRoot, "gcloud")
		if err := os.MkdirAll(secretSub, 0700); err != nil {
			t.Fatalf("MkdirAll secretSub failed: %v", err)
		}
		if err := os.WriteFile(filepath.Join(secretSub, "adc.json"), []byte("SECRET_ADC"), 0600); err != nil {
			t.Fatalf("WriteFile adc.json failed: %v", err)
		}

		srcRoot := t.TempDir()
		if err := os.Symlink(secretRoot, filepath.Join(srcRoot, "escape")); err != nil {
			t.Fatalf("os.Symlink escape failed: %v", err)
		}

		dstDir := filepath.Join(t.TempDir(), "out")
		err := CopyDirNoSymlinks(srcRoot, "escape/gcloud", dstDir, false)
		if !errors.Is(err, ErrSymlink) {
			t.Errorf("CopyDirNoSymlinks(escape/gcloud) = %v, want ErrSymlink", err)
		}
		if _, statErr := os.Lstat(dstDir); !errors.Is(statErr, os.ErrNotExist) {
			t.Errorf("dstDir %q exists after failed copy (err=%v), want not exist", dstDir, statErr)
		}
	})

	t.Run("rejects_clean_dest_containing_symlink_without_deleting_target", func(t *testing.T) {
		canary := filepath.Join(t.TempDir(), "canary.txt")
		if err := os.WriteFile(canary, []byte("KEEP"), 0600); err != nil {
			t.Fatalf("WriteFile canary failed: %v", err)
		}
		srcRoot := t.TempDir()
		if err := os.WriteFile(filepath.Join(srcRoot, "a.yaml"), []byte("a: 1\n"), 0644); err != nil {
			t.Fatalf("WriteFile a.yaml failed: %v", err)
		}
		dstDir := filepath.Join(t.TempDir(), "existing")
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			t.Fatalf("MkdirAll dstDir failed: %v", err)
		}
		if err := os.Symlink(canary, filepath.Join(dstDir, "inner_link")); err != nil {
			t.Fatalf("os.Symlink inner_link failed: %v", err)
		}

		err := CopyDirNoSymlinks(srcRoot, ".", dstDir, true)
		if !errors.Is(err, ErrSymlink) {
			t.Errorf("CopyDirNoSymlinks with symlink inside cleanDst = %v, want ErrSymlink", err)
		}
		got, err := os.ReadFile(canary)
		if err != nil || string(got) != "KEEP" {
			t.Errorf("canary = %q (err=%v), want %q", string(got), err, "KEEP")
		}
	})

	t.Run("rolls_back_staging_when_post_stage_errors", func(t *testing.T) {
		srcRoot := t.TempDir()
		if err := os.WriteFile(filepath.Join(srcRoot, "a.yaml"), []byte("a: 1\n"), 0644); err != nil {
			t.Fatalf("WriteFile a.yaml failed: %v", err)
		}
		parent := t.TempDir()
		dstDir := filepath.Join(parent, "out")
		boom := errors.New("post-stage failure")
		err := StageAndCommitCopyDir(srcRoot, ".", dstDir, false, func(stageDir string) error {
			return boom
		})
		if !errors.Is(err, boom) {
			t.Errorf("StageAndCommitCopyDir = %v, want %v", err, boom)
		}
		entries, readErr := os.ReadDir(parent)
		if readErr != nil {
			t.Fatalf("ReadDir parent failed: %v", readErr)
		}
		if len(entries) != 0 {
			t.Errorf("parent dir has leftover entries %v, want empty", entries)
		}
	})
}
