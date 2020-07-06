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
	"strings"
)

// Resource represents a KRM resource that can be applied
type Resource interface {
	Apply(registry *registry) error
	GetReference() Reference
	GetDependencies() []Reference
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

// BaseResource contains fields should be present in all Resources. This
// struct should be embedded in types implementing the resource interface.
type BaseResource struct {
	TypeMeta
	Metadata Metadata
}

// GetReference computes the reference to the Resource.
func (rs *BaseResource) GetReference() Reference {
	groupAndVersion := strings.Split(rs.APIVersion, "/")
	return Reference{
		Group: groupAndVersion[0],
		Kind:  rs.Kind,
		Name:  rs.Metadata.Name,
	}
}

// GetDependencies returns the dependencies for the BaseResource
func (rs *BaseResource) GetDependencies() (r []Reference) {
	return
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
