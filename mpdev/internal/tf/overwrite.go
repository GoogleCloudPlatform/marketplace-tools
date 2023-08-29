// Copyright 2023 Google LLC
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

package tf

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"sigs.k8s.io/yaml"

	"os"
	"path"
)

const metadataFile = "metadata.yaml"
const metadataDisplayFile = "metadata.display.yaml"

type overwriteConfig struct {
	Variables    []string
	Replacements map[string]string
}

type EnumValueLabel struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// OverwriteTf replaces default variable values in Terraform modules
func OverwriteTf(config *overwriteConfig, dir string) error {
	fmt.Printf("Replacing the default values of the variables: %s\n", config.Variables)
	fmt.Printf("Mapping of values to replace: %s\n", config.Replacements)

	for _, varname := range config.Variables {
		varInfo, err := getVarInfo(varname, dir)
		if err != nil {
			return err
		}

		if varInfo.Default == nil {
			return fmt.Errorf("image variable: %s must have default value", varname)
		}

		defaultVal, ok := varInfo.Default.(string)
		if !ok {
			return fmt.Errorf("image variable: %s must be type string", varname)
		}

		replaceVal, ok := config.Replacements[defaultVal]
		if !ok {
			return fmt.Errorf("default value: %s of variable: %s not found in replacements",
				defaultVal, varname)
		}

		err = overwriteFile(varInfo.Pos.Filename, varname, replaceVal)
		if err != nil {
			return err
		}
	}

	fmt.Println("Successfully replaced default values in tf files")
	return nil
}

// GetOverwriteConfig parses overwriteConfig from a byte array
func GetOverwriteConfig(b []byte) (*overwriteConfig, error) {
	var config overwriteConfig
	err := json.Unmarshal(b, &config)
	if err != nil {
		return nil, fmt.Errorf("failure parsing overwrite config: %s error: %w", string(b), err)
	}

	return &config, nil
}

func getVarInfo(varname string, dir string) (*tfconfig.Variable, error) {
	module, diag := tfconfig.LoadModule(dir)
	if diag.HasErrors() {
		return nil, fmt.Errorf("failure parsing terraform module: %w", diag)
	}

	variable, ok := module.Variables[varname]
	if !ok {
		return nil, fmt.Errorf("variable: %s not found in module", varname)
	}

	return variable, nil
}

func getAttributeValueTokens(value string) hclwrite.Tokens {
	// Use logic similar to https://github.com/hashicorp/hcl/blob/4679383728fe331fc8a6b46036a27b8f818d9bc0/hclwrite/generate.go#L217-L234
	// for writing string values
	return hclwrite.Tokens{
		{
			Type:         hclsyntax.TokenOQuote,
			Bytes:        []byte(`"`),
			SpacesBefore: 1,
		},
		{
			Type:  hclsyntax.TokenQuotedLit,
			Bytes: []byte(value),
		},
		{
			Type:  hclsyntax.TokenCQuote,
			Bytes: []byte(`"`),
		},
	}

}

func overwriteFile(filename string, varname string, value string) error {
	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	file, diag := hclwrite.ParseConfig(b, "", hcl.Pos{Line: 1, Column: 1})
	if diag.HasErrors() {
		return diag
	}

	block := file.Body().FirstMatchingBlock("variable", []string{varname})
	if block == nil {
		return fmt.Errorf("did not find block with variable: %s", varname)
	}

	// SetAttributeValue() is cleaner to overwrite values, however SetAttributeRaw gives more
	// control over formatting. SetAttributeValue() and File.WriteTo() would overwrite all
	// formatting. See: https://github.com/hashicorp/hcl/issues/316
	block.Body().SetAttributeRaw("default", getAttributeValueTokens(value))

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0000)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(file.BuildTokens(nil).Bytes())
	return err
}

// OverwriteMetadata replaces default values for variables in Blueprints Metadata
//
// This overwrite rearranges the order of keys in the YAML, since we are not using the
// definition in
// https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/blob/master/cli/bpmetadata/types.go
//
// We are not using this definition because
// 1. There are version compatiblilty issues with Kpt and cloud-foundation-toolkit to resolve
// 2. We will avoid dropping fields if mpdev is using an out-of-date version of cloud-foundation-toolkit
func OverwriteMetadata(config *overwriteConfig, dir string) error {
	fmt.Printf("Replacing the default values of the variables: %s in %s\n",
		config.Variables, metadataFile)

	data, err := os.ReadFile(path.Join(dir, metadataFile))
	if err != nil {
		// CLI only modules will not have a metadata file. Ignore file not found errors
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	json, err := yaml.YAMLToJSON(data)
	if err != nil {
		return fmt.Errorf("failure parsing %s error: %w", metadataFile, err)
	}

	for _, variable := range config.Variables {
		query := fmt.Sprintf(`spec.interfaces.variables.#(name=="%s").defaultValue`, variable)
		defaultVal := gjson.GetBytes(json, query).String()
		if defaultVal == "" {
			return fmt.Errorf("Missing valid default value for variable: %s in %s",
				variable, metadataFile)
		}
		replaceVal, ok := config.Replacements[defaultVal]
		if !ok {
			return fmt.Errorf("default value: %s of variable: %s in %s not found"+
				" in replacements", defaultVal, variable, metadataFile)
		}

		json, err = sjson.SetBytes(json, query, replaceVal)
		if err != nil {
			return fmt.Errorf("Error setting default value of variable: %s. error: %w",
				variable, err)
		}
	}

	modifiedYaml, err := yaml.JSONToYAML([]byte(json))
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(dir, metadataFile), modifiedYaml, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully replaced default values in %s\n", metadataFile)
	return nil
}

// OverwriteDisplay replaces variable values in Blueprint metadata display file.
func OverwriteDisplay(config *overwriteConfig, dir string) error {
	fmt.Printf("Replacing the values of the display variables: %s in %s\n",
		config.Variables, metadataDisplayFile)

	data, err := os.ReadFile(path.Join(dir, metadataDisplayFile))
	if err != nil {
		// CLI only modules will not have a metadata display file. Ignore file not found errors
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	json, err := yaml.YAMLToJSON(data)
	if err != nil {
		return fmt.Errorf("failure parsing %s error: %w", metadataDisplayFile, err)
	}

	for _, variable := range config.Variables {
		variableQuery := fmt.Sprintf(`spec.ui.input.variables.%s`, variable)
		variableInfo := gjson.GetBytes(json, variableQuery).String()
		if variableInfo == "" {
			return fmt.Errorf("missing valid display info for variable: %s in %s",
				variable, metadataDisplayFile)
		}

		enumValueLabels := gjson.Get(variableInfo, "enumValueLabels").Array()
		if len(enumValueLabels) == 0 {
			fmt.Printf("No enum value labels for display variable: %s in %s\n",
				variable, metadataDisplayFile)
			continue
		}

		var replacementEnumValueLabels []EnumValueLabel
		for _, enumValueLabel := range enumValueLabels {
			currValue := enumValueLabel.Get("value").String()
			currLabel := enumValueLabel.Get("label").String()
			replaceVal, ok := config.Replacements[currValue]
			if !ok {
				return fmt.Errorf("enum value: %s of variable: %s in %s not found"+
					" in replacements", currValue, variable, metadataDisplayFile)
			}
			replacementEnumValueLabels = append(replacementEnumValueLabels, EnumValueLabel{Label: currLabel, Value: replaceVal})
		}

		enumQuery := fmt.Sprintf(`spec.ui.input.variables.%s.enumValueLabels`, variable)
		json, err = sjson.SetBytes(json, enumQuery, replacementEnumValueLabels)
		if err != nil {
			return fmt.Errorf("error setting default value of variable: %s. error: %w",
				variable, err)
		}

	}

	modifiedYaml, err := yaml.JSONToYAML([]byte(json))
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(dir, metadataDisplayFile), modifiedYaml, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully replaced display values in %s\n", metadataDisplayFile)
	return nil
}
