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
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
	"gopkg.in/yaml.v3"
	"k8s.io/utils/exec"
)

// Registry stores references to all resources and can apply
// all resources in the registry
type Registry interface {
	RegisterResource(resource Resource, workingDirectory string)
	GetExecutor() exec.Interface
	GetResource(reference Reference) Resource
	ResolveFilePath(rs Resource, path string) (string, error)
	Apply(dryRun bool) error
	Test(dryRun bool) error
}

type registry struct {
	refMap   map[Reference]Resource
	dirMap   map[Reference]string
	executor exec.Interface
}

// NewRegistry creates a registry that stores references to all resources
func NewRegistry(executor exec.Interface) Registry {
	return &registry{
		refMap:   map[Reference]Resource{},
		dirMap:   map[Reference]string{},
		executor: executor,
	}
}

func (r *registry) GetResource(reference Reference) Resource {
	return r.refMap[reference]
}

func (r *registry) GetExecutor() exec.Interface {
	return r.executor
}

// RegisterResource adds a resource to the registry
func (r *registry) RegisterResource(rs Resource, workingDirectory string) {
	ref := rs.GetReference()
	r.refMap[ref] = rs
	r.dirMap[ref] = workingDirectory
}

// Apply invokes `Apply` on all resources in the registry.
func (r *registry) Apply(dryRun bool) error {
	resources, err := r.topologicalSort()
	if err != nil {
		return err
	}

	for _, resource := range resources {
		fmt.Printf("Starting to validate/create resource %+v\n", resource.GetReference())
		applyErr := resource.Apply(r, dryRun)
		if applyErr != nil {
			applyErr := errors.Wrapf(applyErr, "Error in resource %+v\n", resource.GetReference())
			// Accumulate errors if dryRun
			if dryRun {
				err = multierror.Append(applyErr, err)
			} else {
				return applyErr
			}
		}
	}
	fmt.Printf("all resources have been validated/created\n")

	return err
}

// Apply invokes `Test` on all resources in the registry.
func (r *registry) Test(dryRun bool) error {
	resources, err := r.topologicalSort()
	if err != nil {
		return err
	}

	for _, resource := range resources {
		fmt.Printf("Starting to test resource %+v\n", resource.GetReference())
		testErr := resource.Test(r, dryRun)
		if testErr != nil {
			testErr := errors.Wrapf(testErr, "Error in resource %+v\n", resource.GetReference())
			// Accumulate errors if dryRun
			if dryRun {
				err = multierror.Append(testErr, err)
			} else {
				return testErr
			}
		}
	}
	fmt.Printf("all resources have been tested\n")

	return err
}

func (r *registry) ResolveFilePath(rs Resource, path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}

	return filepath.Abs(filepath.Join(r.dirMap[rs.GetReference()], path))
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
			depRes := r.GetResource(depRef)
			if depRes == nil {
				return []Resource{}, fmt.Errorf("resource not found with reference %+v", depRef)
			}

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

// PopulateRegistryFromFiles registers resources declared in list of files
func PopulateRegistryFromFiles(r Registry, files []string) error {
	for _, file := range files {
		objs, err := decodeFile(file)
		if err != nil {
			return err
		}

		dir := filepath.Dir(file)

		for _, obj := range objs {
			resource, err := UnstructuredToResource(obj)
			if err != nil {
				return err
			}
			r.RegisterResource(resource, dir)
		}
	}
	return nil
}

func decodeFile(file string) ([]Unstructured, error) {
	var objs []Unstructured

	var f *os.File
	var err error
	if file == "-" {
		f = os.Stdin
	} else {
		f, err = os.Open(file)
		if err != nil {
			return objs, err
		}
		defer f.Close()
	}

	dec := yaml.NewDecoder(f)
	for err == nil {
		var m Unstructured
		err = dec.Decode(&m)
		if err == nil {
			objs = append(objs, m)
		}
	}

	if err != io.EOF {
		return objs, errors.Wrap(err, "failed to parse yaml")
	}

	return objs, nil
}
