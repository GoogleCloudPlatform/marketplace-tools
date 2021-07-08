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
credFilePath: /usr/local/google/home/willtang/docker-volume/cred.json
productLifecycleTestConfig:
  provider: providers/e2e-testing
  productExternalName: procurement-ingestion-enabled-tester-with-subscription-usage.endpoints.cloud-marketplace-testing.cloud.goog
  billingAccount: billingAccounts/014E03-1FB93C-054E24
  plans: [plan-a, plan-d]
  approveEntitlementTimeoutSeconds: 1500
  approvePlanChangeTimeoutSeconds: 1500`,
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

			err = saasListingTemplate.Apply(r, false)

			assert.Equal(t, len(tc.expectedRunArgs), fcmd.RunCalls)
			assert.Equal(t, tc.expectedRunArgs, fcmd.RunLog)
		})
	}
}
