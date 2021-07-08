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
	"fmt"
	"regexp"

	"github.com/GoogleContainerTools/kpt/commands"
	"github.com/spf13/cobra"
)

// GetMpdevCommands returns all command.
func GetMpdevCommands(name string) (c []*cobra.Command) {
	pkgCmd := commands.GetPkgCommand(name)
	cfgCmd := commands.GetConfigCommand(name)
	fixDocs(regexp.MustCompile(`\bkpt\b`), name, pkgCmd, cfgCmd)
	applyCmd := GetApplyCommand()
	testCmd := GetTestCommand()

	c = append(c, pkgCmd, cfgCmd, applyCmd, testCmd, versionCmd)

	// apply cross-cutting issues to command
	commands.NormalizeCommand(c...)
	return c
}

// Replace occurrences of "kpt" with "mpdev" in docs.
func fixDocs(old *regexp.Regexp, new string, cmd ...*cobra.Command) {
	for _, c := range cmd {
		c.Use = old.ReplaceAllString(c.Use, new)
		c.Short = old.ReplaceAllString(c.Short, new)
		c.Long = old.ReplaceAllString(c.Long, new)
		c.Example = old.ReplaceAllString(c.Example, new)
		fixDocs(old, new, c.Commands()...)
	}
}

// Variables populated by linking. See: https://github.com/bazelbuild/rules_go/blob/master/go/core.rst#defines-and-stamping
// and x_deps in https://github.com/GoogleCloudPlatform/marketplace-tools/blob/master/mpdev/cmd/BUILD.bazel
var (
	version   string
	gitCommit string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of mpdev",
	Run: func(cmd *cobra.Command, args []string) {
		if version == "" {
			version = "unknown version (built from source)"
		}
		if gitCommit == "" {
			gitCommit = "unknown commit (dirty git repo)"
		}
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("GitCommit: %s\n", gitCommit)
	},
}
