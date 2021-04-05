// Copyright 2020 Google LLC
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

package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"k8s.io/utils/exec"
)

var (
	tmpDirOverride = os.Getenv("MPDEV_TMP_DIR")
)

// ZipDirectory zips the given directory to the given zipFile and returns
// the path of the zipFile which is directory/zipFile
func ZipDirectory(executor exec.Interface, zipFile string, directory string) error {
	if directory == "" || zipFile == "" {
		return fmt.Errorf("directory: %s or zipFile: %s cannot be empty string", directory, zipFile)
	}

	cmd := executor.Command("zip", "-r", zipFile, ".")
	cmd.SetDir(directory)
	cmd.SetStdout(os.Stdout)
	cmd.SetStderr(os.Stderr)

	err := cmd.Run()
	return err
}

// OsTempDir gets os.TempDir() (usually provided by $TMPDIR) but expands any symlinks found within it.
// This wrapper function can prevent problems with docker-for-mac trying to use /var/..., which is not typically
// shared/mounted. It will be expanded via the /var symlink to /private/var/...
func OsTempDir() (string, error) {
	if tmpDirOverride != "" {
		return tmpDirOverride, nil
	}

	dirName := os.TempDir()
	tmpDir, err := filepath.EvalSymlinks(dirName)
	if err != nil {
		return "", err
	}
	tmpDir = filepath.Clean(tmpDir)
	return tmpDir, nil
}

// CreateTmpDir creates a temporary directory and returns its path as a string.
func CreateTmpDir(prefix string) (string, error) {
	systemTempDir, err := OsTempDir()
	if err != nil {
		return "", err
	}
	fullPath, err := ioutil.TempDir(systemTempDir, prefix)
	if err != nil {
		return "", err
	}
	return fullPath, nil
}
