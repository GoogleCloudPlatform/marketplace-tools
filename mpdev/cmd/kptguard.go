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

package cmd

import (
	"errors"
	"fmt"

	"github.com/GoogleCloudPlatform/marketplace-tools/mpdev/internal/pkgget"
	"github.com/spf13/cobra"
)

// hardenPkgCommand prunes unmaintained kpt pkg subcommands (leaving only get)
// and replaces the get subcommand execution path with pkgget.Run so that Git
// option/transport injection and repo-embedded symlinks fail closed.
func hardenPkgCommand(pkgCmd *cobra.Command) {
	if pkgCmd == nil {
		return
	}
	pkgCmd.Run = nil
	pkgCmd.Args = func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("unknown command %q for %q", args[0], cmd.CommandPath())
		}
		return nil
	}
	pkgCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	}
	subs := append([]*cobra.Command(nil), pkgCmd.Commands()...)
	for _, sub := range subs {
		if sub.Name() != "get" {
			pkgCmd.RemoveCommand(sub)
			continue
		}
		sub.PreRunE = nil
		sub.PreRun = nil
		sub.Args = cobra.ExactArgs(2)
		if sub.Flags().Lookup("clean") == nil {
			sub.Flags().Bool("clean", false, "Remove existing destination directory before writing package")
		}
		sub.RunE = func(c *cobra.Command, args []string) error {
			pattern, err := c.Flags().GetString("pattern")
			if err != nil {
				return err
			}
			autoSet, err := c.Flags().GetBool("auto-set")
			if err != nil {
				return err
			}
			clean, err := c.Flags().GetBool("clean")
			if err != nil {
				return err
			}
			return pkgget.Run(args, pkgget.Options{
				Pattern:        pattern,
				AutoSet:        autoSet,
				Clean:          clean,
				AllowLocalRepo: false,
				Stdin:          c.InOrStdin(),
				Stdout:         c.OutOrStdout(),
			})
		}
	}
}

// guardCfgSet disables automatic Starlark/function execution on kpt's cfg set
// subcommand, failing closed if --auto-run is missing or explicitly enabled.
func guardCfgSet(cfgCmd *cobra.Command) {
	if cfgCmd == nil {
		return
	}
	for _, sub := range cfgCmd.Commands() {
		if sub.Name() != "set" {
			continue
		}
		if f := sub.Flags().Lookup("auto-run"); f != nil {
			f.DefValue = "false"
			_ = f.Value.Set("false")
			_ = sub.Flags().MarkHidden("auto-run")
		}
		orig := sub.RunE
		sub.RunE = func(c *cobra.Command, args []string) error {
			if orig == nil {
				return errors.New("cfg set handler is unexpectedly nil")
			}
			f := c.Flags().Lookup("auto-run")
			if f == nil {
				return errors.New("cfg set missing expected --auto-run flag; refusing to execute")
			}
			if f.Changed && f.Value.String() != "false" {
				return errors.New("--auto-run is disabled for security; automatic Starlark/function execution is not permitted")
			}
			if err := f.Value.Set("false"); err != nil {
				return fmt.Errorf("disabling --auto-run: %w", err)
			}
			if len(args) == 0 {
				return errors.New("cfg set requires a package directory argument")
			}
			return orig(c, args)
		}
	}
}
