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

// This package receives an autogen.yaml file from stdin and converts
// to the schema used by the deploymentmanager-autogen
// See: https://github.com/GoogleCloudPlatform/deploymentmanager-autogen/blob/master/example-config/solution.yaml
package main

import (
	"io/ioutil"
	"os"

	autogen "github.com/GoogleCloudPlatform/marketplace-tools/mpdev/autogen/generated"
	yaml2json "github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	var c command
	cmd := &cobra.Command{
		Use:  "autogen",
		RunE: c.autogenConvert,
	}

	cmd.Flags().StringVar(&c.PartnerID, "partnerId", "", "partner id for autogen template")
	cmd.MarkFlagRequired("partnerId")
	cmd.Flags().StringVar(&c.SolutionID, "solutionId", "", "solution id for autogen template")
	cmd.MarkFlagRequired("solutionId")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

type command struct {
	PartnerID  string
	SolutionID string
}

func (c *command) autogenConvert(_ *cobra.Command, _ [] string) error {
	buff, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return errors.Wrap(err, "failed to parse yaml from stdin")
	}

	// Convert yaml to json because the protobuf generated code
	// only has annotations for json, not yaml
	jsonBytes, err := yaml2json.YAMLToJSON(buff)
	if err != nil {
		return err
	}

	var autogenSpec autogen.DeploymentPackageAutogenSpec
	err = protojson.Unmarshal(jsonBytes, &autogenSpec)
	if err != nil {
		return errors.Wrap(err, "failed to convert yaml/json to DeploymentPackageAutogenSpec")
	}

	deploymentPackageInput := autogen.DeploymentPackageInput{
		Spec:       &autogenSpec,
		PartnerId:  c.PartnerID,
		SolutionId: c.SolutionID,
	}

	out, err := protojson.Marshal(&deploymentPackageInput)
	if err != nil {
		return err
	}

	_, err = os.Stdout.Write(out)
	if err != nil {
		return err
	}

	return nil
}
