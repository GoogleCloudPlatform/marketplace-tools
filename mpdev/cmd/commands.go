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

package cmd

import (
	"github.com/GoogleContainerTools/kpt/commands"
	"github.com/spf13/cobra"
)

// GetMpdevCommands returns all command.
func GetMpdevCommands(name string) []*cobra.Command {
	var c []*cobra.Command
	pkgCmd := commands.GetPkgCommand(name)
	cfgCmd := commands.GetConfigCommand(name)
	applyCmd := GetApplyCommand()

	c = append(c, pkgCmd, cfgCmd, applyCmd)

	// apply cross-cutting issues to command
	commands.NormalizeCommand(c...)
	return c
}
