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
	"os/exec"
	"strings"
)

const apiVersion = "dev.marketplace.cloud.google.com/v1alpha1"

var typeMapper = map[TypeMeta]func() Resource{
	{APIVersion: apiVersion, Kind: "GceImage"}:                         func() Resource { return &GceImage{} },
	{APIVersion: apiVersion, Kind: "PackerGceImageBuilder"}:            func() Resource { return &PackerGceImageBuilder{} },
	{APIVersion: apiVersion, Kind: "DeploymentManagerAutogenTemplate"}: func() Resource { return &DeploymentManagerAutogenTemplate{} },
	{APIVersion: apiVersion, Kind: "DeploymentManagerTemplateOnGCS"}:   func() Resource { return &DeploymentManagerTemplateOnGCS{} },
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

// Metadata is metadata that all KRM resources must have
type Metadata struct {
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
	Apply() error
	GetReference() Reference
	SetReferenceMap(ReferenceMap)
}

// ResourceShared contains fields should be present in all Resources. This
// struct should be embedded in types implementing the resource interface.
type ResourceShared struct {
	TypeMeta
	Metadata Metadata

	referenceMap ReferenceMap
}

// GetReference computes the reference to the Resource.
func (rs *ResourceShared) GetReference() Reference {
	groupAndVersion := strings.Split(rs.APIVersion, "/")
	return Reference{
		Group: groupAndVersion[0],
		Kind:  rs.Kind,
		Name:  rs.Metadata.Name,
	}
}

// SetReferenceMap sets a reference to resource map.
func (rs *ResourceShared) SetReferenceMap(referenceMap ReferenceMap) {
	rs.referenceMap = referenceMap
}

// Reference allows a Resource to reference another Resource as part of its
// specification. The combination of Group, Kind, Name MUST be unique for all
// applied resources.
type Reference struct {
	Group string
	Kind  string
	Name  string
}

// ReferenceMap is a mapping between KRM references and the resource object.
type ReferenceMap map[Reference]Resource

// ContainerProcess construct a command to execute the container process
type ContainerProcess struct {
	containerImage string
	processArgs    []string
	mounts         []string
}

func (cp *ContainerProcess) getCommand() *exec.Cmd {
	args := []string{"docker", "run", "--rm", "-i"}
	for _, mount := range cp.mounts {
		args = append(args, "--mount", mount)
	}
	args = append(args, cp.containerImage)
	args = append(args, cp.processArgs...)
	return exec.Command(args[0], args[1:]...)
}
