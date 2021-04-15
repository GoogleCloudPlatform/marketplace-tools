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

package resources

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

const apiVersion = "dev.marketplace.cloud.google.com/v1alpha1"

var typeMapper = map[TypeMeta]func() Resource{
	{APIVersion: apiVersion, Kind: "GceImage"}:                         func() Resource { return &GceImage{} },
	{APIVersion: apiVersion, Kind: "PackerGceImageBuilder"}:            func() Resource { return &PackerGceImageBuilder{} },
	{APIVersion: apiVersion, Kind: "DeploymentManagerAutogenTemplate"}: func() Resource { return &DeploymentManagerAutogenTemplate{} },
	{APIVersion: apiVersion, Kind: "DeploymentManagerTemplate"}:        func() Resource { return &DeploymentManagerTemplate{} },
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
		return nil, errors.Wrapf(err, "unable to marshal resource with kind: %s", typeMeta.Kind)
	}

	err = json.Unmarshal(b, resource)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to unmarshal resource with kind: %s", typeMeta.Kind)
	}

	return resource, nil
}

// Unstructured is used to unmarshal json/yaml KRM Resources and extract
// the TypeMeta, such that the KRM Resource can be unmarshalled to a
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
