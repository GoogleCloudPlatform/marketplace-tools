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
	"path/filepath"
	"strings"

	"github.com/GoogleCloudPlatform/marketplace-tools/mpdev/internal/util"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// DeploymentManagerAutogenTemplate generates a deployment manager template
// given an autogen.yaml file.
type DeploymentManagerAutogenTemplate struct {
	BaseResource
	AutogenSpec interface{}
	PartnerID   string `json:"partnerId"`
	SolutionID  string `json:"solutionId"`

	outDir string
}

// Apply generates a deployment manager template from an autogen file.
func (dm *DeploymentManagerAutogenTemplate) Apply(registry Registry) error {
	dir, err := ioutil.TempDir("", "autogen")
	if err != nil {
		return err
	}
	dm.outDir = dir

	inputDir, err := ioutil.TempDir("", "autogenInput")
	if err != nil {
		return err
	}
	defer os.RemoveAll(inputDir)

	inputFile, err := os.Create(filepath.Join(inputDir, "autogen.yaml"))
	if err != nil {
		return err
	}

	enc := yaml.NewEncoder(inputFile)
	err = enc.Encode(dm.AutogenSpec)
	if err != nil {
		return errors.Wrap(err, "failed to write autogen spec to temp file")
	}

	err = dm.runAutogen(registry, inputDir)
	return err
}

func (dm *DeploymentManagerAutogenTemplate) runAutogen(registry Registry, inputDir string) error {
	autogenImg := "gcr.io/cloud-marketplace-tools/dm/autogen"
	args := []string{"--input_type", "YAML", "--single_input", "/autogen/autogen.yaml",
		"--output_type", "PACKAGE", "--output", "/tmp/out"}

	cp := newContainerProcess(
		registry.GetExecutor(),
		autogenImg,
		args,
		[]mount{
			&bindMount{src: dm.outDir, dst: "/tmp/out"},
			&bindMount{src: inputDir, dst: "/autogen"},
		},
	)
	cmd := cp.getCommand()
	cmd.SetStderr(os.Stderr)
	cmd.SetStdout(os.Stdout)

	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to execute autogen")
	}

	fmt.Printf("Wrote autogen output to directory: %s\n", dm.outDir)

	return nil
}

// DeploymentManagerTemplate saves a referenced Deployment Manager
// template to GCS or the local filesystem
type DeploymentManagerTemplate struct {
	BaseResource
	DeploymentManagerRef Reference
	// Uploads to gcs if file path prefixed with "gs://". Otherwise will
	// zip to given local file path.
	ZipFilePath string
}

// GetDependencies returns dependencies for DeploymentManagerTemplate
func (dm *DeploymentManagerTemplate) GetDependencies() (r []Reference) {
	r = append(r, dm.DeploymentManagerRef)
	return r
}

// Apply uploads a Deployment Manager template to GCS.
func (dm *DeploymentManagerTemplate) Apply(registry Registry) error {
	dmRef := registry.GetResource(dm.DeploymentManagerRef)
	if dmRef == nil {
		return fmt.Errorf("autogen template not found %+v", dm.DeploymentManagerRef)
	}

	dmTemplate, ok := dmRef.(*DeploymentManagerAutogenTemplate)
	if !ok {
		return fmt.Errorf("referenced autogen template is not correct type %+v", dm.DeploymentManagerRef)
	}

	var localZipPath string
	isGCSUpload := strings.HasPrefix(dm.ZipFilePath, "gs://")
	if isGCSUpload {
		localZipPath = filepath.Join(dmTemplate.outDir, "dm_template.zip")
	} else {
		var err error
		localZipPath, err = registry.ResolveFilePath(dm, dm.ZipFilePath)
		if err != nil {
			return errors.Wrapf(err, "failed to resolve path to zipFile: %s", dm.ZipFilePath)
		}
	}

	executor := registry.GetExecutor()
	err := util.ZipDirectory(executor, localZipPath, dmTemplate.outDir)
	if err != nil {
		return errors.Wrapf(err, "failed to zip DM template to %s", localZipPath)
	}
	fmt.Printf("DM template zipped to %s\n", localZipPath)

	if isGCSUpload {
		cmd := executor.Command("gsutil", "cp", localZipPath, dm.ZipFilePath)
		cmd.SetStdout(os.Stdout)
		cmd.SetStderr(os.Stderr)

		fmt.Printf("Uploading DM template to GCS from:%s to:%s\n", localZipPath, dm.ZipFilePath)

		err = cmd.Run()
		if err != nil {
			return errors.Wrap(err, "failed to copy DM template to GCS")
		}

		fmt.Printf("Uploaded DM template to GCS path: %s\n", dm.ZipFilePath)
	}

	return nil
}
