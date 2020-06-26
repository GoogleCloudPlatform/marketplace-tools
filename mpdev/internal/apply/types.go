// Copyright 2020 Google LLC
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

package apply

import (
	"encoding/json"
	"fmt"
)

const apiVersion = "dev.marketplace.cloud.google.com/v1alpha1"

var typeMapper = map[TypeMeta]func() Resource{
	{APIVersion: apiVersion, Kind: "GceImage"}:              func() Resource { return &GceImage{} },
	{APIVersion: apiVersion, Kind: "PackerGceImageBuilder"}: func() Resource { return &PackerGceImageBuilder{} },
}

// UnstructuredToResource converts Unstructured to a specific type implementing the
// Resource interface, using the TypeMeta from the Unstructured object.
func UnstructuredToResource(obj Unstructured) (Resource, error) {
	typeMeta := obj.getTypeMeta()
	fn := typeMapper[typeMeta]
	if fn == nil {
		return nil, fmt.Errorf("unknown Kind: %s. APIVersion: %s", typeMeta.Kind, typeMeta.APIVersion)
	}
	resource := fn()
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// TypeMeta describes an individual KRM resource with strings representing the
// type of the object and its API schema version.
type TypeMeta struct {
	Kind       string
	APIVersion string `yaml:"apiVersion,omitempty"`
}

// ObjectMeta is metadata that all KRM resources must have
type ObjectMeta struct {
	Name        string
	Annotations map[string]string
}

// Unstructured is used to unmarshal json/yaml KRM resources and extract
// the TypeMeta, such that the KRM resource can be unmarshalled to a
// specific type implementing the Resource interface.
type Unstructured map[string]interface{}

func (u *Unstructured) getTypeMeta() TypeMeta {
	kind, ok := (*u)["kind"].(string)
	if !ok {
		return TypeMeta{}
	}

	apiVersion, ok := (*u)["apiVersion"].(string)
	if !ok {
		return TypeMeta{}
	}
	return TypeMeta{
		Kind:       kind,
		APIVersion: apiVersion,
	}
}

// Resource represents a KRM resource that can be applied
type Resource interface {
	Apply()
}

// Reference allows a Resource to reference another Resource as part of its
// specification. The combination of Group, Kind, Name MUST be unique for all
// applied resources.
type Reference struct {
	Group string
	Kind  string
	Name  string
}
