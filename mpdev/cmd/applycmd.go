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
	"io"
	"os"
	"path/filepath"

	"github.com/GoogleCloudPlatform/marketplace-tools/mpdev/internal/apply"
	"github.com/GoogleCloudPlatform/marketplace-tools/mpdev/internal/docs"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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
		RunE:    c.RunE,
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

// RunE Executes the `apply` command
func (c *command) RunE(_ *cobra.Command, _ []string) (err error) {
	registry := apply.NewRegistry(exec.New())
	for _, file := range c.Filenames {
		objs, err := decodeFile(file)
		if err != nil {
			return err
		}

		dir := filepath.Dir(file)

		for _, obj := range objs {
			resource, err := apply.UnstructuredToResource(obj)
			if err != nil {
				return err
			}
			registry.RegisterResource(resource, dir)
		}
	}

	err = registry.Apply(c.DryRun)

	return err
}

func decodeFile(file string) ([]apply.Unstructured, error) {
	var objs []apply.Unstructured

	var f *os.File
	var err error
	if file == "-" {
		f = os.Stdin
	} else {
		f, err = os.Open(file)
		if err != nil {
			return objs, err
		}
		defer f.Close()
	}

	dec := yaml.NewDecoder(f)
	for err == nil {
		var m apply.Unstructured
		err = dec.Decode(&m)
		if err == nil {
			objs = append(objs, m)
		}
	}

	if err != io.EOF {
		return objs, errors.Wrap(err, "failed to parse yaml")
	}

	return objs, nil
}
