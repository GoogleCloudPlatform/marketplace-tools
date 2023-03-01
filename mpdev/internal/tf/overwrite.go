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
	"os"
)

type overwriteConfig struct {
	Variables    []string
	Replacements map[string]string
}

// Overwrite replaces variable values in Terraform modules
func Overwrite(configBytes []byte, dir string) error {
	config, err := getOverwriteConfig(configBytes)

	if err != nil {
		return fmt.Errorf("failure parsing overwrite config: %s error: %w", string(configBytes), err)
	}
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
	return nil
}

func getOverwriteConfig(b []byte) (*overwriteConfig, error) {
	var config overwriteConfig
	err := json.Unmarshal(b, &config)
	if err != nil {
		return nil, err
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
			Type:  hclsyntax.TokenOQuote,
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
