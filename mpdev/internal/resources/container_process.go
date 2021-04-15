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

package resources

import (
	"fmt"
	"k8s.io/utils/exec"
)

type containerProcess struct {
	executor       exec.Interface
	containerImage string
	processArgs    []string
	mounts         []mount
}

// newContainerProcess constructs a command to execute the container process
func newContainerProcess(executor exec.Interface, containerImage string, processArgs []string, mounts []mount) *containerProcess {
	return &containerProcess{
		executor:       executor,
		containerImage: containerImage,
		processArgs:    processArgs,
		mounts:         mounts,
	}
}

type mount interface {
	getMount() string
}

type bindMount struct {
	src string
	dst string
}

func (bm *bindMount) getMount() string {
	return fmt.Sprintf("type=bind,src=%s,dst=%s", bm.src, bm.dst)
}

func (cp *containerProcess) getCommand() exec.Cmd {
	args := []string{"docker", "run", "--rm", "-i"}
	for _, mount := range cp.mounts {
		args = append(args, "--mount", mount.getMount())
	}
	args = append(args, cp.containerImage)
	args = append(args, cp.processArgs...)
	return cp.executor.Command(args[0], args[1:]...)
}
