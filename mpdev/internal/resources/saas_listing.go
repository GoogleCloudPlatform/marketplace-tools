// Copyright 2021 Google LLC
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

package resources

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/GoogleCloudPlatform/marketplace-tools/mpdev/internal/util"
	"github.com/pkg/errors"
	"k8s.io/utils/exec"
)

// SaasListingTestTemplate has the parameters for a SaasListing integration test.
type SaasListingTestTemplate struct {
	BaseResource
	IntegrationTestConfig PartnerIntegrationTestConfig `yaml:"integrationTestConfig"`
	CredFilePath          string                       `yaml:"credFilePath"`
}

// Apply is a no-op.
func (template *SaasListingTestTemplate) Apply(registry Registry, dryRun bool) error {
	return nil
}

type BillingMeteringDriver struct {
	DriverCommand             string `json:"driverCommand"`
	PlanId                    string `json:"planId"`
	BillingMeteringConfigPath string `json:"billingMeteringConfigPath"`
}

type BillingMeteringTestConfig struct {
	Driver BillingMeteringDriver `json:"driver"`
}

// PartnerIntegrationTestConfig has the parameters that need to go into the test container.
type PartnerIntegrationTestConfig struct {
	Provider                         string                      `json:"provider"`
	ProductExternalName              string                      `json:"productExternalName"`
	BillingAccount                   string                      `json:"billingAccount"`
	Plans                            []string                    `json:"plans"`
	ApproveEntitlementTimeoutSeconds int32                       `json:"approveEntitlementTimeoutSeconds"`
	ApprovePlanChangeTimeoutSeconds  int32                       `json:"approvePlanChangeTimeoutSeconds"`
	BillingMeteringTestConfig        []BillingMeteringTestConfig `json:"billingMeteringTestConfig"`
}

// Test takes a SaasListingTestTemplate and starts the partner integration test container.
func (template *SaasListingTestTemplate) Test(registry Registry, dryRun bool) error {
	if dryRun {
		return nil
	}

	dir, err := util.CreateTmpDir("mp-saas-integ-test")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	err = copyFile(registry.GetExecutor(), template.CredFilePath, filepath.Join(dir, "cred.json"))
	if err != nil {
		return err
	}

	// copy all billing metering expectation files to `dir` so that they could be accessible from within the Docker container
	for index, billingMeteringDriverConfig := range template.IntegrationTestConfig.BillingMeteringTestConfig {
		srcFilePath := billingMeteringDriverConfig.Driver.BillingMeteringConfigPath

		// create a new config filename using the index and the original filename's extension.
		// This is required to prevent file name collision after we copy the files to `dir`.
		// Example: index=3 and SrcFilePath=/path/to/file/config.json, then dstFileName=3.json
		dstFileName := strings.Join([]string{strconv.Itoa(index), filepath.Ext(srcFilePath)}, "")
		dstFilePath := filepath.Join(dir, dstFileName)

		err = copyFile(registry.GetExecutor(), srcFilePath, dstFilePath)
		if err != nil {
			return err
		}
		
		template.IntegrationTestConfig.BillingMeteringTestConfig[index].Driver.BillingMeteringConfigPath = dstFilePath
	}

	configFileName := "partner_integration_test_config.json"

	err = template.IntegrationTestConfig.write(dir, configFileName)
	if err != nil {
		return err
	}

	testImg := "gcr.io/marketplace-saas-tools/mp-saas-partner-integ-test"

	err = dockerPull(registry.GetExecutor(), testImg)
	if err != nil {
		return err
	}

	dockerArgs := []string{"-e", "GOOGLE_APPLICATION_CREDENTIALS=/input/cred.json"}
	processArgs := []string{"/input/" + configFileName}

	file, err := os.Open(filepath.Join(dir, configFileName))
	if err != nil {
		return err
	}

	configContent, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	fmt.Printf("config file that will be passed to the Test Framework Container\n%s\n", configContent)

	cp := newContainerProcess(
		registry.GetExecutor(),
		dockerArgs,
		testImg,
		processArgs,
		[]mount{
			&bindMount{src: dir, dst: "/input"},
		},
	)
	cmd := cp.getCommand()
	cmd.SetStderr(os.Stderr)
	cmd.SetStdout(os.Stdout)

	fmt.Printf("Executing partner integration test container %s\n", testImg)
	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to execute partner integration test container with docker")
	}

	return nil
}

func dockerPull(executor exec.Interface, imageURL string) error {
	cmd := executor.Command("docker", "pull", imageURL)
	cmd.SetStderr(os.Stderr)
	cmd.SetStderr(os.Stdout)

	return cmd.Run()
}
func copyFile(executor exec.Interface, source string, dest string) error {
	cmd := executor.Command("cp", source, dest)

	cmd.SetStderr(os.Stderr)
	cmd.SetStderr(os.Stdout)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (tc *PartnerIntegrationTestConfig) write(dir string, filename string) error {
	encodedJSON, err := json.Marshal(tc)

	if err != nil {
		return err
	}

	fmt.Println(string(encodedJSON[:]))

	err = ioutil.WriteFile(filepath.Join(dir, filename), encodedJSON, 0755)
	if err != nil {
		return err
	}

	return nil
}
