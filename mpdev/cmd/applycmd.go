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

	"github.com/GoogleCloudPlatform/marketplace-tools/mpdev/internal/apply"
	"github.com/GoogleCloudPlatform/marketplace-tools/mpdev/internal/docs"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// GetApplyCommand returns `apply` command used to create mpdev resources.
func GetApplyCommand() *cobra.Command {
	var c command
	cmd := &cobra.Command{
		Use:     "apply",
		Short:   docs.ApplyShort,
		Long:    docs.ApplyLong,
		Example: docs.ApplyExamples,
		RunE:    c.RunE,
	}

	cmd.Flags().StringSliceVarP(&c.Filenames, "filename", "f", c.Filenames, "that contains the configuration to apply")

	return cmd
}

type command struct {
	Filenames []string
}

// RunE Executes the `apply` command
func (c *command) RunE(_ *cobra.Command, _ []string) error {
	var allObjs []apply.Unstructured
	for _, file := range c.Filenames {
		objs, err := decodeFile(file)
		if err != nil {
			return err
		}
		allObjs = append(allObjs, objs...)
	}

	var resources []apply.Resource
	var refMap = map[apply.Reference]apply.Resource{}
	for _, obj := range allObjs {
		resource, err := apply.UnstructuredToResource(obj)
		if err != nil {
			return err
		}
		refMap[resource.GetReference()] = resource
		resource.SetReferenceMap(refMap)

		resources = append(resources, resource)
	}

	// TODO: Setup execution order for apply command
	// Some resources cannot be applied prior to others
	for _, resource := range resources {
		err := resource.Apply()
		if err != nil {
			return err
		}
	}

	return nil
}

func decodeFile(file string) ([]apply.Unstructured, error) {
	var objs []apply.Unstructured
	f, err := os.Open(file)
	if err != nil {
		return objs, err
	}
	defer f.Close()

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
