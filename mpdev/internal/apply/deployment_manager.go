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
	Spec AutogenSpec

	outDir string
}

// AutogenSpec is defines the spec used for auto-generating deployment packages.
type AutogenSpec struct {
	// Deployment Spec is documented in https://github.com/GoogleCloudPlatform/marketplace-tools/docs/autogen-reference.md
	DeploymentSpec map[string]interface{} `yaml:"deploymentSpec"`
	PackageInfo    PackageInfo            `yaml:"packageInfo"`
}

// PackageInfo describes the software packaged in a deployable solution. PackageInfo
// is metadata displayed on the VM solution details page in the GCP marketplace
// console.
type PackageInfo struct {
	// Version of combined software components
	Version string
	// Name and version of OS
	OsInfo component `yaml:"osInfo"`
	// Names and versions of software components
	Components []component
}

type component struct {
	Name    string
	Version string
}

type convertedSpec struct {
	PartnerID    string                 `yaml:"partnerId"`
	SolutionID   string                 `yaml:"solutionId"`
	Spec         map[string]interface{} `yaml:"spec"`
	PartnerInfo  map[string]interface{} `yaml:"partnerInfo"`
	SolutionInfo map[string]interface{} `yaml:"solutionInfo"`
}

// Apply generates a deployment manager template from an autogen file.
func (dm *DeploymentManagerAutogenTemplate) Apply(registry Registry, dryRun bool) error {
	err := dm.validateSpec()
	if err != nil {
		return err
	}

	if dryRun {
		return nil
	}

	convertedSpec := dm.convertToAutogen()

	dir, err := util.CreateTmpDir("autogen")
	if err != nil {
		return err
	}
	dm.outDir = dir

	inputDir, err := util.CreateTmpDir("autogenInput")
	if err != nil {
		return err
	}
	defer os.RemoveAll(inputDir)

	inputFile, err := os.Create(filepath.Join(inputDir, "autogen.yaml"))
	if err != nil {
		return err
	}

	enc := yaml.NewEncoder(inputFile)
	err = enc.Encode(convertedSpec)
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

	fmt.Printf("Executing autogen container: %s\n", autogenImg)
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to execute autogen container with docker")
	}

	fmt.Printf("Wrote autogen output to directory: %s\n", dm.outDir)

	return nil
}

func (dm *DeploymentManagerAutogenTemplate) convertToAutogen() *convertedSpec {
	// Placeholder fields are either unused by autogen or overidden in the partner
	// portal UI when configuring solution details/metadata
	autogenSpec := &convertedSpec{
		PartnerID:    "placeholder",
		SolutionID:   "solution",
		Spec:         dm.Spec.DeploymentSpec,
		PartnerInfo:  map[string]interface{}{"name": "placeholder"},
		SolutionInfo: map[string]interface{}{"name": "placeholder", "version": dm.Spec.PackageInfo.Version},
	}

	pkgGroups := []struct {
		Type       string `yaml:",omitempty"`
		Components []component
	}{{
		"SOFTWARE_GROUP_OS",
		[]component{dm.Spec.PackageInfo.OsInfo},
	}, {
		Components: dm.Spec.PackageInfo.Components,
	}}

	autogenSpec.SolutionInfo["packagedSoftwareGroups"] = pkgGroups

	return autogenSpec
}

func (dm *DeploymentManagerAutogenTemplate) validateSpec() error {
	packageInfo := dm.Spec.PackageInfo
	osInfo := packageInfo.OsInfo
	if osInfo.Version == "" || osInfo.Name == "" {
		return fmt.Errorf("osInfo version or name not specified. Ensure spec.packageInfo.osInfo in config file is set")
	}
	if len(packageInfo.Components) == 0 {
		return fmt.Errorf("no packageInfo Components. Ensure spec.packageInfo.Components in config file is set")
	}

	// Further deploymentSpec schema checks are done when executing autogen container.
	if len(dm.Spec.DeploymentSpec) == 0 {
		return fmt.Errorf("no deploymentSpec contents. Ensure spec.deploymentSpec in config file is set")
	}
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
func (dm *DeploymentManagerTemplate) Apply(registry Registry, dryRun bool) error {
	dmRef := registry.GetResource(dm.DeploymentManagerRef)
	if dmRef == nil {
		return fmt.Errorf("autogen template not found %+v", dm.DeploymentManagerRef)
	}

	dmTemplate, ok := dmRef.(*DeploymentManagerAutogenTemplate)
	if !ok {
		return fmt.Errorf("referenced autogen template is not correct type %+v", dm.DeploymentManagerRef)
	}

	if dm.ZipFilePath == "" {
		return errors.New("ZipFilePath cannot be empty for DM template")
	}

	if dryRun {
		return nil
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
