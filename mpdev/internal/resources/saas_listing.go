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

// ConnectionInfo has the needed information to connect to BigQuery billing export
type ConnectionInfo struct {
	Project   string `json:"project"`
	TableName string `json:"tableName"`
}

// Expectation outlines the results from running the user provided driver
type Expectation struct {
	SkuID            string           `json:"skuId"`
	UsageExpectation UsageExpectation `json:"usageExpectation"`
	CostExpectation  CostExpectation  `json:"costExpectation"`
}

// UsageExpectation has usage metrics that are compared by the metering validator
type UsageExpectation struct {
	Min       float64 `json:"min"`
	Max       float64 `json:"max"`
	BaseUnits string  `json:"baseUnits"`
}

// CostExpectation has metrics that are compared by the metering validator
type CostExpectation struct {
	Min      float64 `json:"min"`
	Max      float64 `json:"max"`
	Currency string  `json:"currency"`
}

// BillingMeteringDriver is the billing metering driver is executed by the test framework
type BillingMeteringDriver struct {
	DriverCommand  string         `json:"driverCommand"`
	PlanID         string         `json:"planId"`
	Expectation    Expectation    `json:"expectation"`
	ConnectionInfo ConnectionInfo `json:"connectionInfo"`
}

// BillingMeteringTestConfig has the parameters that are needed to run the billing metering test
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

	configFileName := "partner_integration_test_config.json"

	err = template.IntegrationTestConfig.write(dir, configFileName)
	if err != nil {
		return err
	}

	testImg := "gcr.io/marketplace-saas-tools/mp-saas-test-framework"

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
			&bindMount{src: "/var/run/docker.sock", dst: "/var/run/docker.sock"},
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
