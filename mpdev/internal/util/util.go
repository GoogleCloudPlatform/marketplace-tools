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
	"k8s.io/utils/exec"
	"os"
)

// ZipDirectory zips the given directory to the given zipFile and returns
// the path of the zipFile which is directory/zipFile
func ZipDirectory(executor exec.Interface, zipFile string, directory string) error {
	if directory == "" || zipFile == "" {
		return fmt.Errorf("directory: %s or zipFile: %s cannot be empty string", directory, zipFile)
	}

	cmd := executor.Command("/bin/sh", "-c", fmt.Sprintf("cd %s && zip -r %s .", directory, zipFile))
	cmd.SetStdout(os.Stdout)
	cmd.SetStderr(os.Stderr)

	err := cmd.Run()
	return err
}
