// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package tf

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestOverwrite(t *testing.T) {

	testcases := []struct {
		name            string
		tfFiles         map[string]string
		expectedTfFiles map[string]string
		overwriteConfig string
		errorContains   string
	}{{
		name: "Overwrite multiple variables and files",
		tfFiles: map[string]string{
			"main.tf":        mainTf,
			"anyfilename.tf": otherTf,
		},
		expectedTfFiles: map[string]string{
			"main.tf":        mainTfReplaced,
			"anyfilename.tf": otherTfReplaced,
		},
		overwriteConfig: `
{
  "variables": ["value_to_replace", "other_value_to_replace", "another_variable"],
  "replacements": {
	  "original-value": "new-value",
	  "old-value": "newer-value",
	  "oldest-value": "newest-value"
  }
}
`,
	}, {
		name: "Invalid HCL shows parsing error",
		tfFiles: map[string]string{
			"main.tf": "this is broken",
		},
		overwriteConfig: `
{
  "variables": ["value_to_replace"],
  "replacements": { "original-value": "new-value" }
}
`,
		errorContains: "failure parsing terraform module",
	}, {
		name: "Invalid overwrite config shows parsing error",
		tfFiles: map[string]string{
			"main.tf": mainTf,
		},
		overwriteConfig: "not valid json",
		errorContains:   "failure parsing overwrite config",
	}, {
		name: "Fail when variable not present in Terraform module",
		tfFiles: map[string]string{
			"main.tf": mainTf,
		},
		overwriteConfig: `
{
  "variables": ["missing_variable"],
  "replacements": { "original-value": "new-value" }
}
`,
		errorContains: "variable: missing_variable not found",
	}, {
		name: "Fail when variable has no default value set",
		tfFiles: map[string]string{
			"main.tf": tfNoDefault,
		},
		overwriteConfig: `
{
  "variables": ["value_to_replace"],
  "replacements": { "original-value": "new-value" }
}
`,
		errorContains: "variable: value_to_replace must have default value",
	}, {
		name: "Fail when variable default value is not a string",
		tfFiles: map[string]string{
			"main.tf": tfDefaultWrongType,
		},
		overwriteConfig: `
{
  "variables": ["value_to_replace"],
  "replacements": { "original-value": "new-value" }
}
`,
		errorContains: "image variable: value_to_replace must be type string",
	}, {
		name: "Fail when variable default value is not in replacements",
		tfFiles: map[string]string{
			"main.tf": mainTf,
		},
		overwriteConfig: `
{
  "variables": ["value_to_replace"],
  "replacements": { "non-existent": "new-value" }
}
`,
		errorContains: "default value: original-value of variable: value_to_replace not found in replacements",
	}}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "tftest")
			assert.NoError(t, err)
			defer os.RemoveAll(tmpDir)

			for file, content := range tc.tfFiles {
				err = os.WriteFile(path.Join(tmpDir, file), []byte(content), 0600)
				assert.NoError(t, err)
			}

			err = Overwrite([]byte(tc.overwriteConfig), tmpDir)

			if tc.errorContains == "" {
				assert.NoError(t, err)

				actualContents, err := getDirContents(tmpDir)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedTfFiles, actualContents)
			} else {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.errorContains)
			}
		})
	}
}

func getDirContents(dir string) (map[string]string, error) {
	fileContents := make(map[string]string)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		fileContents[path[len(dir)+1:]] = string(b)

		return nil
	})
	return fileContents, err
}

var mainTf string = `
resource "google_compute_instance_template" "template" {
  name = "template"
}

variable "value_to_replace" {
  type = string
  default = "original-value"
}

variable "other_value_to_replace" {
  type = string
  default = "old-value"
}
`

var mainTfReplaced string = `
resource "google_compute_instance_template" "template" {
  name = "template"
}

variable "value_to_replace" {
  type = string
  default = "new-value"
}

variable "other_value_to_replace" {
  type = string
  default = "newer-value"
}
`

var otherTf string = `
variable "another_variable" {
  type = string
  default = "oldest-value"
}
`

var otherTfReplaced string = `
variable "another_variable" {
  type = string
  default = "newest-value"
}
`

var tfNoDefault string = `
variable "value_to_replace" {
  type = string
}
`

var tfDefaultWrongType string = `
variable "value_to_replace" {
  type = map(number)
  default = {
    foo = 2
    bar = 4
  }
}
`
