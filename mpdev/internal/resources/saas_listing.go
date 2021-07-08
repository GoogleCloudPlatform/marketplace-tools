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

type SaasListingTestTemplate struct {
	BaseResource
	IntegrationTestConfig PartnerIntegrationTestConfig `yaml:"integrationTestConfig"`
	CredFilePath          string                       `yaml:"credFilePath"`
}

func (template *SaasListingTestTemplate) Apply(registry Registry, dryRun bool) error {
	return nil
}

type PartnerIntegrationTestConfig struct {
	Provider                         string   `json:"provider"`
	ProductExternalName              string   `json:"productExternalName"`
	BillingAccount                   string   `json:"billingAccount"`
	Plans                            []string `json:"plans"`
	ApproveEntitlementTimeoutSeconds int32    `json:"approveEntitlementTimeoutSeconds"`
	ApprovePlanChangeTimeoutSeconds  int32    `json:"approvePlanChangeTimeoutSeconds"`
}

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

	testImg := "gcr.io/marketplace-saas-tools/mp-saas-partner-integ-test"

	dockerPull(registry.GetExecutor(), testImg)

	dockerArgs := []string{"-e", "GOOGLE_APPLICATION_CREDENTIALS=/input/cred.json"}
	processArgs := []string{"/input/" + configFileName}

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

func dockerPull(executor exec.Interface, imageUrl string) error {
	cmd := executor.Command("docker", "pull", imageUrl)
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
	encodedJson, err := json.Marshal(tc)

	if err != nil {
		return err
	}

	fmt.Println(string(encodedJson[:]))

	err = ioutil.WriteFile(filepath.Join(dir, filename), encodedJson, 0755)
	if err != nil {
		return err
	}

	return nil
}
