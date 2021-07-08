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
  provider: providers/e2e-testing
  productExternalName: procurement-ingestion-enabled-tester-with-subscription-usage.endpoints.cloud-marketplace-testing.cloud.goog
  billingAccount: billingAccounts/014E03-1FB93C-054E24
  plans: [plan-a, plan-d]
  approveEntitlementTimeoutSeconds: 1500
  approvePlanChangeTimeoutSeconds: 1500`,
		expectedRunArgs: [][]string{
			{"cp", "/docker-volume/cred.json", "/tmp/mp-saas-integ-test\\d+/cred.json"},
			{"docker", "pull", "gcr.io/marketplace-saas-tools/mp-saas-partner-integ-test"},
			{"docker", "run", "--rm", "-i", "--mount", "type=bind,src=/tmp/mp-saas-integ-test\\d+,dst=/input", "-e", "GOOGLE_APPLICATION_CREDENTIALS=/input/cred.json", "gcr.io/marketplace-saas-tools/mp-saas-partner-integ-test", "/input/partner_integration_test_config.json"}},
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
