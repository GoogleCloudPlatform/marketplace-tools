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
	"path/filepath"

	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

// Registry stores references to all resources and can apply
// all resources in the registry
type Registry interface {
	RegisterResource(resource Resource, workingDirectory string)
	GetResource(reference Reference) Resource
	ResolveFilePath(rs Resource, path string) string
	Apply() error
}

type registry struct {
	refMap map[Reference]Resource
	dirMap map[Reference]string
}

// NewRegistry creates a registry that stores references to all resources
func NewRegistry() Registry {
	return &registry{
		refMap: map[Reference]Resource{},
		dirMap: map[Reference]string{},
	}
}

func (r *registry) GetResource(reference Reference) Resource {
	return r.refMap[reference]
}

// RegisterResource adds a resource to the registry
func (r *registry) RegisterResource(rs Resource, workingDirectory string) {
	ref := rs.GetReference()
	r.refMap[ref] = rs
	r.dirMap[ref] = workingDirectory
}

// Apply invokes `Apply` on all resources in the registry.
func (r *registry) Apply() error {
	resources, err := r.topologicalSort()
	if err != nil {
		return err
	}

	for _, resource := range resources {
		err = resource.Apply(r)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *registry) ResolveFilePath(rs Resource, path string) string {
	if !filepath.IsAbs(path) {
		return filepath.Join(r.dirMap[rs.GetReference()], path)
	}

	return path
}

// topologicalSort returns a list of resources such that each
// resource is after its dependencies in the list.
func (r *registry) topologicalSort() ([]Resource, error) {
	dag := simple.NewDirectedGraph()

	// Add resource references as nodes to graph
	var refToID = map[Reference]int64{}
	var idToRef = map[int64]Reference{}
	for ref := range r.refMap {
		n := dag.NewNode()
		refToID[ref] = n.ID()
		idToRef[n.ID()] = ref
		dag.AddNode(n)
	}

	// Add dependencies as edges to graph
	for ref, resource := range r.refMap {
		to := dag.Node(refToID[ref])
		for _, depRef := range resource.GetDependencies() {
			from := dag.Node(refToID[depRef])
			e := dag.NewEdge(from, to)
			dag.SetEdge(e)
		}
	}

	// Execute topological sort
	nodes, err := topo.Sort(dag)
	if err != nil {
		return nil, err
	}

	// Convert node id to resource
	resources := make([]Resource, 0, len(nodes))
	for _, node := range nodes {
		ref := idToRef[node.ID()]
		resources = append(resources, r.refMap[ref])
	}

	return resources, err
}
