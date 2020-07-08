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

package apply

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/GoogleCloudPlatform/marketplace-tools/mpdev/internal/util"
	"github.com/pkg/errors"
)

// DeploymentManagerAutogenTemplate generates a deployment manager template
// given an autogen.yaml file.
type DeploymentManagerAutogenTemplate struct {
	BaseResource
	// TODO(#10) resolve file paths such that file paths are relative to
	// where mpdev.yaml file was read.
	AutogenFile string
	PartnerID   string `json:"partnerId"`
	SolutionID  string `json:"solutionId"`

	outDir string
}

// Apply generates a deployment manager template from an autogen file.
func (dm *DeploymentManagerAutogenTemplate) Apply(registry *registry) error {
	dir, err := ioutil.TempDir("", "autogen")
	if err != nil {
		return err
	}
	dm.outDir = dir

	fmt.Printf("Generating deployment manager template from autogen...\n")
	err = dm.convertAutogenSchema()
	if err != nil {
		return err
	}

	autogenImg := "gcr.io/cloud-marketplace-tools/dm/autogen"
	args := []string{"--input_type", "YAML", "--single_input", "/autogen/autogen.yaml",
		"--output_type", "PACKAGE", "--output", "/autogen"}

	cp := &containerProcess{
		containerImage: autogenImg,
		processArgs:    args,
		mounts:         []mount{&bindMount{src: dm.outDir, dst: "/autogen"}},
	}
	cmd := cp.getCommand()
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to execute autogen")
	}

	fmt.Printf("Wrote autogen output to directory: %s\n", dm.outDir)

	return nil
}

func (dm *DeploymentManagerAutogenTemplate) convertAutogenSchema() error {
	inputFile, err := os.Open(dm.AutogenFile)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outFile, err := os.Create(filepath.Join(dm.outDir, "autogen.yaml"))
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Image name from running `docker build mpdev/autogen -f mpdev/autogen/Dockerfile -t autogen_converter
	autogenConverterImg := "autogen_converter"
	cp := &containerProcess{
		containerImage: autogenConverterImg,
		processArgs:    []string{"--partnerId", dm.PartnerID, "--solutionId", dm.SolutionID},
		mounts:         nil,
	}
	cmd := cp.getCommand()

	cmd.Stderr = os.Stderr
	cmd.Stdin = inputFile
	cmd.Stdout = outFile

	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to execute autogen converter")
	}

	return nil
}

// DeploymentManagerTemplateOnGCS uploads a referenced Deployment Manager
// template to GCS.
type DeploymentManagerTemplateOnGCS struct {
	BaseResource
	DeploymentManagerRef Reference
	GCS                  struct {
		Bucket string
		Object string
	}
}

// GetDependencies returns dependencies for DeploymentManagerTemplateOnGCS
func (dm *DeploymentManagerTemplateOnGCS) GetDependencies() (r []Reference) {
	r = append(r, dm.DeploymentManagerRef)
	return r
}

// Apply uploads a Deployment Manager template to GCS.
func (dm *DeploymentManagerTemplateOnGCS) Apply(registry *registry) error {
	dmRef := registry.getReference(dm.DeploymentManagerRef)
	if dmRef == nil {
		return fmt.Errorf("autogen template not found %+v", dm.DeploymentManagerRef)
	}

	dmTemplate, ok := dmRef.(*DeploymentManagerAutogenTemplate)
	if !ok {
		return fmt.Errorf("referenced autogen template is not correct type %+v", dm.DeploymentManagerRef)
	}

	zipFileName := "dm_template.zip"
	zipFilePath, err := util.ZipDirectory(zipFileName, dmTemplate.outDir)
	if err != nil {
		return errors.Wrap(err, "failed to zip autogen template before uploading to GCS")
	}

	gcsPath := fmt.Sprintf("gs://%s/%s", dm.GCS.Bucket, dm.GCS.Object)

	cmd := exec.Command("gsutil", "cp", zipFilePath, gcsPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Uploading DM template to GCS. Running command: %v\n", cmd)

	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to copy autogen template to GCS")
	}

	fmt.Printf("Uploaded DM template to GCS path: %s\n", gcsPath)

	return nil
}
