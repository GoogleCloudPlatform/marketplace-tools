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
	"github.com/GoogleCloudPlatform/marketplace-tools/mpdev/internal/docs"
	"github.com/GoogleCloudPlatform/marketplace-tools/mpdev/internal/resources"
	"github.com/spf13/cobra"
	"k8s.io/utils/exec"
)

// GetApplyCommand returns `apply` command used to create mpdev resources.
func GetApplyCommand() *cobra.Command {
	var c command
	cmd := &cobra.Command{
		Use:     "apply -f FILENAME [--dryrun]",
		Short:   docs.ApplyShort,
		Long:    docs.ApplyLong,
		Example: docs.ApplyExamples,
		RunE:    c.applyRunE,
	}

	cmd.Flags().BoolVar(&c.DryRun, "dryrun", c.DryRun, "if set, validates configuration files without creating resource")
	cmd.Flags().StringSliceVarP(&c.Filenames, "filename", "f", c.Filenames, "that contains the configuration to apply")
	_ = cobra.MarkFlagRequired(cmd.Flags(), "filename")

	return cmd
}

type command struct {
	Filenames []string
	DryRun    bool
}

func (c *command) applyRunE(_ *cobra.Command, _ []string) (err error) {
	registry := resources.NewRegistry(exec.New())
	err = resources.PopulateRegistryFromFiles(registry, c.Filenames)
	if err != nil {
		return err
	}

	err = registry.Apply(c.DryRun)
	return err
}
