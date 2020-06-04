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

	autogen "github.com/GoogleCloudPlatform/marketplace-tools/autogen/generated"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	yaml2json "github.com/ghodss/yaml"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	cmd := &cobra.Command{
		Use:  "autogen",
		RunE: autogenConvert,
	}

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func autogenConvert(_ *cobra.Command, _ [] string) error {
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

	partnerId := "placeholder"
	solutionId := "placeholder"
	deploymentPackageInput := autogen.DeploymentPackageInput{
		Spec:       &autogenSpec,
		PartnerId:  partnerId,
		SolutionId: solutionId,
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
