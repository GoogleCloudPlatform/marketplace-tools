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
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"k8s.io/utils/exec"
	testingexec "k8s.io/utils/exec/testing"
)

func TestSaasListingTestTemplate_Test(t *testing.T) {
	testCases := []struct {
		name            string
		dryrun          bool
		config          string
		expectedRunArgs [][]string
	}{{
		name:   "Happy path",
		dryrun: false,
		config: `
apiVersion: dev.marketplace.cloud.google.com/v1alpha1
kind: SaasListingTemplate
credFilePath: /docker-volume/cred.json
integrationTestConfig:
  provider: "providers/e2e-testing"
  productExternalName: "procurement-ingestion-enabled-tester-with-subscription-usage.endpoints.cloud-marketplace-testing.cloud.goog"
  billingAccount: "billingAccounts/017778-0B6CC1-FB92E9"
  plans: ["plan-a", "plan-d"]
  approveEntitlementTimeoutSeconds: 600
  approvePlanChangeTimeoutSeconds: 600
  billingMeteringTestConfig:
    - driver:
        driverCommand: "docker run -v /var/run/docker.sock:/var/run/docker.sock -v /tmp/saas:/keys -e GOOGLE_APPLICATION_CREDENTIALS=/keys/cred.json -e GOOGLE_CLOUD_PROJECT=testing-producer-saas gcr.io/marketplace-saas-tools/billing-metering-driver --service=procurement-ingestion-enabled-tester-with-subscription-usage.endpoints.cloud-marketplace-testing.cloud.goog --metric=procurement-ingestion-enabled-tester-with-subscription-usage.endpoints.cloud-marketplace-testing.cloud.goog/plan_a_metric_1 --consumer_id=project_number:{@USAGE_REPORTING_ID}"
        planId: "plan-a"
        connectionInfo:
          project: "projects/testing-producer-saas-322600"
          tableName: "testing-producer-saas-322600.testing_producer_saas_billing_export.gcp_billing_export_resource_v1_017778_0B6CC1_FB92E9"
        expectation:
          skuId: "0D8F-CA36-DD7A"
          usageExpectation:
            min: 150
            max: 200
            baseUnits: "requests"
          costExpectation:
            min: 800
            max: 1000
            currency: "USD"`,
		expectedRunArgs: [][]string{
			{"cp", "/docker-volume/cred.json", "/tmp/mp-saas-integ-test\\d+/cred.json"},
			{"docker", "pull", "gcr.io/marketplace-saas-tools/mp-saas-test-framework"},
			{"docker", "run", "--rm", "-i", "--mount", "type=bind,src=/tmp/mp-saas-integ-test\\d+,dst=/input", "--mount", "type=bind,src=/var/run/docker.sock,dst=/var/run/docker.sock", "-e", "GOOGLE_APPLICATION_CREDENTIALS=/input/cred.json", "gcr.io/marketplace-saas-tools/mp-saas-test-framework", "/input/partner_integration_test_config.json"}},
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var saasListingTemplate SaasListingTestTemplate

			err := yaml.Unmarshal([]byte(tc.config), &saasListingTemplate)

			assert.NoError(t, err)

			fcmd := testingexec.FakeCmd{
				RunScript: []testingexec.FakeRunAction{
					func() ([]byte, []byte, error) { return nil, nil, nil },
					func() ([]byte, []byte, error) { return nil, nil, nil },
					func() ([]byte, []byte, error) { return nil, nil, nil },
				},
			}

			executor := &testingexec.FakeExec{
				CommandScript: []testingexec.FakeCommandAction{
					func(cmd string, args ...string) exec.Cmd { return testingexec.InitFakeCmd(&fcmd, cmd, args...) },
					func(cmd string, args ...string) exec.Cmd { return testingexec.InitFakeCmd(&fcmd, cmd, args...) },
					func(cmd string, args ...string) exec.Cmd { return testingexec.InitFakeCmd(&fcmd, cmd, args...) },
				},
			}
			r := NewRegistry(executor)

			err = saasListingTemplate.Test(r, false)
			assert.NoError(t, err)

			assert.Equal(t, len(tc.expectedRunArgs), fcmd.RunCalls)

			for i, arr := range fcmd.RunLog {
				for j, s := range arr {
					res, _ := regexp.MatchString(tc.expectedRunArgs[i][j], s)
					fmt.Println(res, s, tc.expectedRunArgs[i][j])
					assert.True(t, res, "Command has incorrect value")
				}
			}
		})
	}
}
