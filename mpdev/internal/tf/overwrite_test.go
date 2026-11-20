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
	"errors"
	"os"
	"path"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/GoogleCloudPlatform/marketplace-tools/mpdev/internal/fsutil"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestOverwriteTf(t *testing.T) {

	testcases := []struct {
		name            string
		tfFiles         map[string]string
		expectedTfFiles map[string]string
		overwriteConfig overwriteConfig
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
		overwriteConfig: overwriteConfig{
			Variables: []string{"value_to_replace", "other_value_to_replace", "another_variable"},
			Replacements: map[string]string{
				"original-value": "new-value",
				"old-value":      "newer-value",
				"oldest-value":   "newest-value",
			},
		},
	}, {
		name: "Invalid HCL shows parsing error",
		tfFiles: map[string]string{
			"main.tf": "this is broken",
		},
		overwriteConfig: overwriteConfig{
			Variables: []string{"value_to_replace"},
			Replacements: map[string]string{
				"original-value": "new-value",
			},
		},
		errorContains: "failure parsing terraform module",
	}, {
		name: "Fail when variable not present in Terraform module",
		tfFiles: map[string]string{
			"main.tf": mainTf,
		},
		overwriteConfig: overwriteConfig{
			Variables: []string{"missing_variable"},
			Replacements: map[string]string{
				"original-value": "new-value",
			},
		},
		errorContains: "variable: missing_variable not found",
	}, {
		name: "Fail when variable has no default value set",
		tfFiles: map[string]string{
			"main.tf": tfNoDefault,
		},
		overwriteConfig: overwriteConfig{
			Variables: []string{"value_to_replace"},
			Replacements: map[string]string{
				"original-value": "new-value",
			},
		},
		errorContains: "variable: value_to_replace must have default value",
	}, {
		name: "Fail when variable default value is not a string",
		tfFiles: map[string]string{
			"main.tf": tfDefaultWrongType,
		},
		overwriteConfig: overwriteConfig{
			Variables: []string{"value_to_replace"},
			Replacements: map[string]string{
				"original-value": "new-value",
			},
		},
		errorContains: "image variable: value_to_replace must be type string",
	}, {
		name: "Fail when variable default value is not in replacements",
		tfFiles: map[string]string{
			"main.tf": mainTf,
		},
		overwriteConfig: overwriteConfig{
			Variables: []string{"value_to_replace"},
			Replacements: map[string]string{
				"non-existent": "new-value",
			},
		},
		errorContains: "default value: original-value of variable: value_to_replace not found in replacements",
	}, {
		name: "With NewValues, overwrite multiple variables and files",
		tfFiles: map[string]string{
			"main.tf":        mainTf,
			"anyfilename.tf": otherTf,
		},
		expectedTfFiles: map[string]string{
			"main.tf":        mainTfReplaced,
			"anyfilename.tf": otherTfReplaced,
		},
		overwriteConfig: overwriteConfig{
			NewValues: map[string]string{
				"value_to_replace":       "new-value",
				"other_value_to_replace": "newer-value",
				"another_variable":       "newest-value",
			},
		},
	}, {
		name: "With NewValues, overwrite multiple variables and files. Add consumer label",
		tfFiles: map[string]string{
			"main.tf": mainTfNoLabel,
		},
		expectedTfFiles: map[string]string{
			"main.tf": mainTfLabelUpserted,
		},
		overwriteConfig: overwriteConfig{
			ConsumerLabel: "new-consumer-label",
			NewValues: map[string]string{
				"value_to_replace": "new-value",
			},
		},
	}, {
		name: "With NewValues, overwrite multiple variables and files. Do not replace existing consumer label",
		tfFiles: map[string]string{
			"main.tf": mainTfProvidedLabel,
		},
		expectedTfFiles: map[string]string{
			"main.tf": mainTfLabelUpserted,
		},
		overwriteConfig: overwriteConfig{
			ConsumerLabel: "even-newer-consumer-label",
			NewValues: map[string]string{
				"value_to_replace": "new-value",
			},
		},
	},
		{
			name: "With NewValues, ignores Variables and Replacements",
			tfFiles: map[string]string{
				"main.tf":        mainTf,
				"anyfilename.tf": otherTf,
			},
			expectedTfFiles: map[string]string{
				"main.tf":        mainTfReplaced,
				"anyfilename.tf": otherTfReplaced,
			},
			overwriteConfig: overwriteConfig{
				NewValues: map[string]string{
					"value_to_replace":       "new-value",
					"other_value_to_replace": "newer-value",
					"another_variable":       "newest-value",
				},
				Variables: []string{"value_to_replace", "other_value_to_replace", "another_variable"},
				Replacements: map[string]string{
					"original-value": "new-value-unused",
					"old-value":      "newer-value-unused",
					"oldest-value":   "newest-value-unused",
				},
			},
		}, {
			name: "With NewValues, adds default value when variable has no default value set",
			tfFiles: map[string]string{
				"main.tf": tfNoDefault,
			},
			expectedTfFiles: map[string]string{
				"main.tf": tfDefaultAdded,
			},
			overwriteConfig: overwriteConfig{
				NewValues: map[string]string{
					"value_to_replace": "new-value",
				},
			},
		}, {
			name: "With NewValues, fail when variable default value is not a string",
			tfFiles: map[string]string{
				"main.tf": tfNoDefaultWrongType,
			},
			overwriteConfig: overwriteConfig{
				NewValues: map[string]string{
					"value_to_replace": "new-value",
				},
			},
			errorContains: "image variable: value_to_replace must be type string",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "tftest")
			assert.NoError(t, err)
			defer os.RemoveAll(tmpDir)

			for file, content := range tc.tfFiles {
				err = os.WriteFile(path.Join(tmpDir, file), []byte(content), 0600)
				assert.NoError(t, err)
			}

			err = OverwriteTf(&tc.overwriteConfig, tmpDir)

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

func TestGetOverwriteConfig(t *testing.T) {
	testcases := []struct {
		name           string
		configBytes    []byte
		expectedConfig overwriteConfig
		errorContains  string
	}{{
		name:          "Invalid overwrite config shows parsing error",
		configBytes:   []byte("not valid json"),
		errorContains: "failure parsing overwrite config",
	}, {
		name: "Parses overwrite config",
		configBytes: []byte(`
{
	"variables": ["source_image"],
	"replacements": {"old_image": "new_image" },
	"newValues": {
		"source_image": "new_image"
	}
}
`),
		expectedConfig: overwriteConfig{
			Variables: []string{"source_image"},
			Replacements: map[string]string{
				"old_image": "new_image",
			},
			NewValues: map[string]string{
				"source_image": "new_image",
			},
		},
	},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			config, err := GetOverwriteConfig([]byte(tc.configBytes))
			if tc.errorContains == "" {
				assert.NoError(t, err)
				assert.Equal(t, &tc.expectedConfig, config)
			} else {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.errorContains)
			}
		})
	}
}

func TestOverwriteMetadata(t *testing.T) {
	testcases := []struct {
		name             string
		originalMetadata string
		expectedMetadata string
		overwriteConfig  overwriteConfig
		errorContains    string
	}{{
		name:             "Overwrite multiple variables",
		originalMetadata: metadata,
		expectedMetadata: metadataReplaced,
		overwriteConfig: overwriteConfig{
			Variables: []string{"source_image", "another_image"},
			Replacements: map[string]string{
				"old-image":   "new-image",
				"older-image": "newer-image",
			},
		},
	}, {
		name:             "Fail when metadata is invalid yaml",
		originalMetadata: "- not validyaml\ninvalid-",
		errorContains:    "failure parsing metadata.yaml",
	}, {
		name:             "Fail when variable not present in Metadata",
		originalMetadata: metadata,
		overwriteConfig: overwriteConfig{
			Variables: []string{"missing_variable"},
			Replacements: map[string]string{
				"original-value": "new-value",
			},
		},
		errorContains: "Missing valid default value for variable: missing_variable",
	}, {
		name:             "Fail when variable has no default value set",
		originalMetadata: metadataNoDefault,
		overwriteConfig: overwriteConfig{
			Variables: []string{"source_image"},
			Replacements: map[string]string{
				"original-value": "new-value",
			},
		},
		errorContains: "Missing valid default value for variable: source_image",
	}, {
		name:             "Fail when variable default value is not in replacements",
		originalMetadata: metadata,
		overwriteConfig: overwriteConfig{
			Variables: []string{"source_image"},
			Replacements: map[string]string{
				"non-existent": "new-value",
			},
		},
		errorContains: "default value: old-image of variable: source_image in metadata.yaml not found in replacements",
	}, {
		name:             "With NewValues, overwrite multiple variables",
		originalMetadata: metadata,
		expectedMetadata: metadataReplaced,
		overwriteConfig: overwriteConfig{
			NewValues: map[string]string{
				"source_image":  "new-image",
				"another_image": "newer-image",
			},
		},
	}, {
		name:             "With NewValues, ignores Variables and Replacements",
		originalMetadata: metadata,
		expectedMetadata: metadataReplaced,
		overwriteConfig: overwriteConfig{
			NewValues: map[string]string{
				"source_image":  "new-image",
				"another_image": "newer-image",
			},
			Variables: []string{"source_image", "another_image"},
			Replacements: map[string]string{
				"old-image":   "new-image-unused",
				"older-image": "newer-image-unused",
			},
		},
	}, {
		name:             "With NewValues, fail when variable is not present in Metadata",
		originalMetadata: metadata,
		overwriteConfig: overwriteConfig{
			NewValues: map[string]string{
				"missing_variable": "new-value",
			},
		},
		errorContains: "missing variable entry for variable: missing_variable",
	}, {
		name:             "With NewValues, adds default value when variable has no default value set",
		originalMetadata: metadataNoDefault,
		expectedMetadata: metadataDefaultAdded,
		overwriteConfig: overwriteConfig{
			NewValues: map[string]string{
				"source_image": "new-value",
			},
		},
	},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "tftest")
			assert.NoError(t, err)
			defer os.RemoveAll(tmpDir)

			err = os.WriteFile(path.Join(tmpDir, "metadata.yaml"),
				[]byte(tc.originalMetadata), 0600)
			assert.NoError(t, err)

			err = OverwriteMetadata(&tc.overwriteConfig, tmpDir)

			if tc.errorContains == "" {
				assert.NoError(t, err)

				metadataBytes, err := os.ReadFile(path.Join(tmpDir, "metadata.yaml"))
				assert.NoError(t, err)
				actualMetadata := make(map[interface{}]interface{})
				expectedMetadata := make(map[interface{}]interface{})

				err = yaml.Unmarshal(metadataBytes, &actualMetadata)
				assert.NoError(t, err)

				err = yaml.Unmarshal([]byte(tc.expectedMetadata), expectedMetadata)
				assert.NoError(t, err)

				assert.Equal(t, expectedMetadata, actualMetadata)
			} else {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.errorContains)
			}
		})
	}
}

func TestOverwriteMetadataNoFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "tftest")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	err = OverwriteMetadata(&overwriteConfig{}, tmpDir)
	assert.NoError(t, err)
}

func TestOverwiteMetadataPermissionError(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "tftest")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	err = os.WriteFile(path.Join(tmpDir, "metadata.yaml"), []byte("file"), 0111)
	assert.NoError(t, err)

	err = OverwriteMetadata(&overwriteConfig{}, tmpDir)
	assert.Error(t, err)
	assert.True(t, os.IsPermission(err))
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

func TestOverwriteDisplay(t *testing.T) {
	testcases := []struct {
		name                    string
		originalMetadataDisplay string
		expectedMetadataDisplay string
		overwriteConfig         overwriteConfig
		errorContains           string
	}{
		{
			name:                    "Overwrite single display variable enum values",
			originalMetadataDisplay: metadataDisplayWithEnumsSingle,
			expectedMetadataDisplay: metadataDisplayWithEnumsSingleReplaced,
			overwriteConfig: overwriteConfig{
				Variables: []string{"source_image"},
				Replacements: map[string]string{
					"projects/click-to-deploy-images/global/images/wordpress-1": "projects/replacement/global/images/wordpress-1-new",
				},
			},
		},
		{
			name:                    "Overwrite multiple display variable enum values",
			originalMetadataDisplay: metadataDisplayWithEnumsDouble,
			expectedMetadataDisplay: metadataDisplayWithEnumsDoubleReplaced,
			overwriteConfig: overwriteConfig{
				Variables: []string{"source_image", "another_image"},
				Replacements: map[string]string{
					"projects/click-to-deploy-images/global/images/wordpress-1": "projects/replacement/global/images/wordpress-1-new",
					"projects/click-to-deploy-images/global/images/wordpress-2": "projects/replacement/global/images/wordpress-2-new",
					"projects/click-to-deploy-images/global/images/wordpress-3": "projects/replacement/global/images/wordpress-3-new",
				},
			},
		},
		{
			name:                    "No changes if no display variable enum value labels",
			originalMetadataDisplay: metadataDisplayNoEnums,
			expectedMetadataDisplay: metadataDisplayNoEnums,
			overwriteConfig: overwriteConfig{
				Variables:    []string{"source_image"},
				Replacements: map[string]string{},
			},
		}, {
			name:                    "Fail when metadata display is invalid yaml",
			originalMetadataDisplay: "- not validyaml\ninvalid-",
			errorContains:           "failure parsing metadata.display.yaml",
		}, {
			name:                    "Fail when display variable not present in Metadata display",
			originalMetadataDisplay: metadata,
			overwriteConfig: overwriteConfig{
				Variables: []string{"missing_variable"},
				Replacements: map[string]string{
					"original-value": "new-value",
				},
			},
			errorContains: "missing valid display info for variable: missing_variable",
		}, {
			name:                    "Fail when display variable enum value is not in replacements",
			originalMetadataDisplay: metadataDisplayWithEnumsSingle,
			overwriteConfig: overwriteConfig{
				Variables: []string{"source_image"},
				Replacements: map[string]string{
					"non-existent": "new-value",
				},
			},
			errorContains: "enum value: projects/click-to-deploy-images/global/images/wordpress-1 of variable: source_image in metadata.display.yaml not found in replacements",
		}}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "tftest")
			assert.NoError(t, err)
			defer os.RemoveAll(tmpDir)

			err = os.WriteFile(path.Join(tmpDir, "metadata.display.yaml"),
				[]byte(tc.originalMetadataDisplay), 0600)
			assert.NoError(t, err)

			err = OverwriteDisplay(&tc.overwriteConfig, tmpDir)

			if tc.errorContains == "" {
				assert.NoError(t, err)

				metadataDisplayBytes, err := os.ReadFile(path.Join(tmpDir, "metadata.display.yaml"))
				assert.NoError(t, err)
				actualMetadataDisplay := make(map[interface{}]interface{})
				expectedMetadataDisplay := make(map[interface{}]interface{})

				err = yaml.Unmarshal(metadataDisplayBytes, &actualMetadataDisplay)
				assert.NoError(t, err)

				err = yaml.Unmarshal([]byte(tc.expectedMetadataDisplay), expectedMetadataDisplay)
				assert.NoError(t, err)

				assert.Equal(t, expectedMetadataDisplay, actualMetadataDisplay)
			} else {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.errorContains)
			}
		})
	}
}

var mainTf string = `
resource "google_compute_instance_template" "template" {
  name = "template"
}

variable "value_to_replace" {
  type    = string
  default = "original-value"
}

variable "other_value_to_replace" {
  type    = string
  default = "old-value"
}
`

var mainTfReplaced string = `
resource "google_compute_instance_template" "template" {
  name = "template"
}

variable "value_to_replace" {
  type    = string
  default = "new-value"
}

variable "other_value_to_replace" {
  type    = string
  default = "newer-value"
}
`

var mainTfNoLabel string = `
provider "google" {
  project = var.project_id
}

resource "google_compute_instance_template" "template" {
  name = "template"
}

variable "value_to_replace" {
  type    = string
  default = "original-value"
}
`

var mainTfProvidedLabel string = `
provider "google" {
  project = var.project_id
  default_labels = {
    goog-partner-solution = "new-consumer-label"
  }
}

resource "google_compute_instance_template" "template" {
  name = "template"
}

variable "value_to_replace" {
  type    = string
  default = "original-value"
}
`

var mainTfLabelUpserted string = `
provider "google" {
  project = var.project_id
  default_labels = {
    goog-partner-solution = "new-consumer-label"
  }
}

resource "google_compute_instance_template" "template" {
  name = "template"
}

variable "value_to_replace" {
  type    = string
  default = "new-value"
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
  type    = string
  default = "newest-value"
}
`

var tfNoDefault string = `
variable "value_to_replace" {
  type = string
}
`

var tfDefaultAdded string = `
variable "value_to_replace" {
  type    = string
  default = "new-value"
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
var tfNoDefaultWrongType string = `
variable "value_to_replace" {
  type = map(number)
}
`

var metadata string = `
spec:
  interfaces:
    variables:
    - name: source_image
      description: The image name for the disk for the VM instance.
      varType: string
      defaultValue: old-image
    - name: another_image
      description: The image name for the disk for the VM instance.
      varType: string
      defaultValue: older-image
`

var metadataReplaced string = `
spec:
  interfaces:
    variables:
    - name: source_image
      description: The image name for the disk for the VM instance.
      varType: string
      defaultValue: new-image
    - name: another_image
      description: The image name for the disk for the VM instance.
      varType: string
      defaultValue: newer-image
`

var metadataNoDefault string = `
spec:
  interfaces:
    variables:
    - name: source_image
      description: The image name for the disk for the VM instance.
      varType: string
`

var metadataDefaultAdded string = `
spec:
  interfaces:
    variables:
    - name: source_image
      description: The image name for the disk for the VM instance.
      varType: string
      defaultValue: new-value
`

var metadataDisplayWithEnumsSingle string = `
spec:
  ui:
    input:
      variables:
        source_image:
          name: source_image
          title: Source Image
          enumValueLabels:
            - label: wordpress-1
              value: projects/click-to-deploy-images/global/images/wordpress-1
          xGoogleProperty:
            type: ET_GCE_DISK_IMAGE
`

var metadataDisplayWithEnumsSingleReplaced string = `
spec:
  ui:
    input:
      variables:
        source_image:
          name: source_image
          title: Source Image
          enumValueLabels:
            - label: wordpress-1
              value: projects/replacement/global/images/wordpress-1-new
          xGoogleProperty:
            type: ET_GCE_DISK_IMAGE
`

var metadataDisplayWithEnumsDouble string = `
spec:
  ui:
    input:
      variables:
        source_image:
          name: source_image
          title: Source Image
          enumValueLabels:
            - label: wordpress-1
              value: projects/click-to-deploy-images/global/images/wordpress-1
            - label: wordpress-2
              value: projects/click-to-deploy-images/global/images/wordpress-2
          xGoogleProperty:
            type: ET_GCE_DISK_IMAGE
        another_image:
          name: another_image
          title: Another Image
          enumValueLabels:
            - label: wordpress-3
              value: projects/click-to-deploy-images/global/images/wordpress-3
          xGoogleProperty:
            type: ET_GCE_DISK_IMAGE
`

var metadataDisplayWithEnumsDoubleReplaced string = `
spec:
  ui:
    input:
      variables:
        source_image:
          name: source_image
          title: Source Image
          enumValueLabels:
            - label: wordpress-1
              value: projects/replacement/global/images/wordpress-1-new
            - label: wordpress-2
              value: projects/replacement/global/images/wordpress-2-new
          xGoogleProperty:
            type: ET_GCE_DISK_IMAGE
        another_image:
          name: another_image
          title: Another Image
          enumValueLabels:
            - label: wordpress-3
              value: projects/replacement/global/images/wordpress-3-new
          xGoogleProperty:
            type: ET_GCE_DISK_IMAGE
`

var metadataDisplayNoEnums string = `
spec:
  ui:
    input:
      variables:
        source_image:
          name: source_image
          title: Source Image
          xGoogleProperty:
            type: ET_GCE_DISK_IMAGE
`

func TestOverwriteRejectsSymlinksAndPreservesCanaries(t *testing.T) {
	t.Run("variables_tf_absolute_and_relative_symlinks_rejected_for_new_values_and_replacements", func(t *testing.T) {
		outsideDir := t.TempDir()
		canaryTf := filepath.Join(outsideDir, "canary.tf")
		const canaryBody = `variable "value_to_replace" {
  type    = string
  default = "CANARY_SECRET_VALUE"
}
`
		assert.NoError(t, os.WriteFile(canaryTf, []byte(canaryBody), 0600))

		for _, useRel := range []bool{false, true} {
			modDir := filepath.Join(outsideDir, "mod")
			assert.NoError(t, os.RemoveAll(modDir))
			assert.NoError(t, os.MkdirAll(modDir, 0755))
			assert.NoError(t, os.WriteFile(filepath.Join(modDir, "main.tf"), []byte(`provider "google" {}`+"\n"), 0644))

			target := canaryTf
			if useRel {
				var err error
				target, err = filepath.Rel(modDir, canaryTf)
				assert.NoError(t, err)
			}
			assert.NoError(t, os.Symlink(target, filepath.Join(modDir, "variables.tf")))

			errNew := OverwriteTf(&overwriteConfig{
				NewValues: map[string]string{"value_to_replace": "pwned"},
			}, modDir)
			assert.ErrorIs(t, errNew, fsutil.ErrSymlink)
			assert.NotContains(t, errNew.Error(), "CANARY_SECRET_VALUE")

			errLegacy := OverwriteTf(&overwriteConfig{
				Variables:    []string{"value_to_replace"},
				Replacements: map[string]string{"CANARY_SECRET_VALUE": "pwned"},
			}, modDir)
			assert.ErrorIs(t, errLegacy, fsutil.ErrSymlink)

			gotCanary, err := os.ReadFile(canaryTf)
			assert.NoError(t, err)
			assert.Equal(t, canaryBody, string(gotCanary))
		}
	})

	t.Run("main_tf_consumer_label_symlink_rejected", func(t *testing.T) {
		outsideDir := t.TempDir()
		canaryMain := filepath.Join(outsideDir, "canary_main.tf")
		assert.NoError(t, os.WriteFile(canaryMain, []byte(mainTfNoLabel), 0600))

		modDir := t.TempDir()
		assert.NoError(t, os.Symlink(canaryMain, filepath.Join(modDir, "main.tf")))

		err := OverwriteTf(&overwriteConfig{ConsumerLabel: "injected-label"}, modDir)
		assert.ErrorIs(t, err, fsutil.ErrSymlink)

		gotCanary, readErr := os.ReadFile(canaryMain)
		assert.NoError(t, readErr)
		assert.Equal(t, mainTfNoLabel, string(gotCanary))
	})

	t.Run("metadata_and_display_symlinks_rejected_preserving_canaries", func(t *testing.T) {
		outsideDir := t.TempDir()
		canaryMeta := filepath.Join(outsideDir, "canary_metadata.yaml")
		canaryDisp := filepath.Join(outsideDir, "canary_display.yaml")
		assert.NoError(t, os.WriteFile(canaryMeta, []byte(metadata), 0600))
		assert.NoError(t, os.WriteFile(canaryDisp, []byte(metadataDisplayWithEnumsSingle), 0600))

		modDir := t.TempDir()
		assert.NoError(t, os.Symlink(canaryMeta, filepath.Join(modDir, "metadata.yaml")))
		err := OverwriteMetadata(&overwriteConfig{
			NewValues: map[string]string{"source_image": "pwned"},
		}, modDir)
		assert.ErrorIs(t, err, fsutil.ErrSymlink)

		dispDir := t.TempDir()
		assert.NoError(t, os.Symlink(canaryDisp, filepath.Join(dispDir, "metadata.display.yaml")))
		err = OverwriteDisplay(&overwriteConfig{
			NewValues: map[string]string{"source_image": "pwned"},
		}, dispDir)
		assert.ErrorIs(t, err, fsutil.ErrSymlink)

		gotMeta, err := os.ReadFile(canaryMeta)
		assert.NoError(t, err)
		assert.Equal(t, metadata, string(gotMeta))

		gotDisp, err := os.ReadFile(canaryDisp)
		assert.NoError(t, err)
		assert.Equal(t, metadataDisplayWithEnumsSingle, string(gotDisp))
	})

	t.Run("chained_symlinks_and_tf_json_symlinks_rejected", func(t *testing.T) {
		outsideDir := t.TempDir()
		canaryJSON := filepath.Join(outsideDir, "canary.tf.json")
		const tfJSON = `{"variable":{"value_to_replace":{"type":"string","default":"secret"}}}`
		assert.NoError(t, os.WriteFile(canaryJSON, []byte(tfJSON), 0600))
		hop1 := filepath.Join(outsideDir, "hop1.tf.json")
		assert.NoError(t, os.Symlink(canaryJSON, hop1))

		modDir := t.TempDir()
		assert.NoError(t, os.Symlink(hop1, filepath.Join(modDir, "vars.tf.json")))

		err := OverwriteTf(&overwriteConfig{
			NewValues: map[string]string{"value_to_replace": "pwned"},
		}, modDir)
		assert.ErrorIs(t, err, fsutil.ErrSymlink)

		gotJSON, readErr := os.ReadFile(canaryJSON)
		assert.NoError(t, readErr)
		assert.Equal(t, tfJSON, string(gotJSON))
	})

	t.Run("dangling_symlinks_fail_closed_and_do_not_create_targets", func(t *testing.T) {
		missingTarget := filepath.Join(t.TempDir(), "does_not_exist.yaml")

		metaDir := t.TempDir()
		assert.NoError(t, os.Symlink(missingTarget, filepath.Join(metaDir, "metadata.yaml")))
		err := OverwriteMetadata(&overwriteConfig{}, metaDir)
		assert.ErrorIs(t, err, fsutil.ErrSymlink)
		_, statErr := os.Lstat(missingTarget)
		assert.True(t, errors.Is(statErr, os.ErrNotExist))

		dispDir := t.TempDir()
		assert.NoError(t, os.Symlink(missingTarget, filepath.Join(dispDir, "metadata.display.yaml")))
		err = OverwriteDisplay(&overwriteConfig{}, dispDir)
		assert.ErrorIs(t, err, fsutil.ErrSymlink)

		tfDir := t.TempDir()
		assert.NoError(t, os.Symlink(missingTarget, filepath.Join(tfDir, "variables.tf")))
		err = OverwriteTf(&overwriteConfig{}, tfDir)
		assert.ErrorIs(t, err, fsutil.ErrSymlink)
	})

	t.Run("unreferenced_and_nested_symlinks_and_special_files_in_tree_rejected", func(t *testing.T) {
		canary := filepath.Join(t.TempDir(), "docker_config.json")
		assert.NoError(t, os.WriteFile(canary, []byte(`{"auths":{}}`), 0600))

		modDir := t.TempDir()
		assert.NoError(t, os.WriteFile(filepath.Join(modDir, "main.tf"), []byte(mainTf), 0644))
		assert.NoError(t, os.MkdirAll(filepath.Join(modDir, "subdir"), 0755))
		assert.NoError(t, os.Symlink(canary, filepath.Join(modDir, "subdir", "unused_leak.txt")))

		err := OverwriteTf(&overwriteConfig{
			NewValues: map[string]string{"value_to_replace": "new-value"},
		}, modDir)
		assert.ErrorIs(t, err, fsutil.ErrSymlink)

		fifoDir := t.TempDir()
		assert.NoError(t, os.WriteFile(filepath.Join(fifoDir, "main.tf"), []byte(mainTf), 0644))
		assert.NoError(t, syscall.Mkfifo(filepath.Join(fifoDir, "named.pipe"), 0600))
		err = OverwriteTf(&overwriteConfig{}, fifoDir)
		assert.ErrorIs(t, err, fsutil.ErrSpecialFile)
	})

	t.Run("module_inside_symlinked_parent_directory_succeeds_for_regular_files", func(t *testing.T) {
		base := t.TempDir()
		realMod := filepath.Join(base, "private_var", "workspace")
		assert.NoError(t, os.MkdirAll(realMod, 0755))
		assert.NoError(t, os.WriteFile(filepath.Join(realMod, "main.tf"), []byte(mainTfNoLabel), 0644))
		assert.NoError(t, os.WriteFile(filepath.Join(realMod, "metadata.yaml"), []byte(metadata), 0644))
		assert.NoError(t, os.WriteFile(filepath.Join(realMod, "metadata.display.yaml"), []byte(metadataDisplayWithEnumsSingle), 0644))

		symlinkVar := filepath.Join(base, "var")
		assert.NoError(t, os.Symlink(filepath.Join(base, "private_var"), symlinkVar))
		viaSymlinkParent := filepath.Join(symlinkVar, "workspace")

		cfg := &overwriteConfig{
			ConsumerLabel: "new-consumer-label",
			NewValues: map[string]string{
				"value_to_replace": "new-value",
			},
		}
		assert.NoError(t, OverwriteTf(cfg, viaSymlinkParent))

		metaCfg := &overwriteConfig{
			NewValues: map[string]string{
				"source_image":  "new-image",
				"another_image": "newer-image",
			},
		}
		assert.NoError(t, OverwriteMetadata(metaCfg, viaSymlinkParent))

		dispCfg := &overwriteConfig{
			NewValues: map[string]string{
				"source_image": "projects/replacement/global/images/wordpress-1-new",
			},
		}
		assert.NoError(t, OverwriteDisplay(dispCfg, viaSymlinkParent))
	})
}
