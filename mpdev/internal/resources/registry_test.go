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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/exec"
)

func TestResolveFilePath(t *testing.T) {
	r := newTestResource("testResource")

	wd, err := os.Getwd()
	assert.NoError(t, err)

	testCases := map[string]string{
		// Absolute path
		"/tmp/foo/2.txt": "/tmp/foo/2.txt",
		// Relative path
		"foo/2.txt":    filepath.Join(wd, "dir/foo/2.txt"),
		"../foo/2.txt": filepath.Join(wd, "foo/2.txt"),
	}

	for path, expected := range testCases {
		registry := NewRegistry(exec.New())
		registry.RegisterResource(r, "dir")
		resolvedPath, err := registry.ResolveFilePath(r, path)
		assert.NoError(t, err)
		assert.Equal(t, expected, resolvedPath)
	}
}

func TestApplyOrder(t *testing.T) {
	order := 1
	applyFunc := func(expectedOrder int) func(r Registry, dryRun bool) error {
		return func(r Registry, _ bool) error {
			assert.Equal(t, expectedOrder, order)
			order++
			return nil
		}
	}

	depFunc := func(resources ...Resource) func() []Reference {
		return func() (refs []Reference) {
			for _, r := range resources {
				refs = append(refs, r.GetReference())
			}
			return refs
		}
	}

	r1 := newTestResourceFunc("r1", applyFunc(1), nil, nil)
	r2 := newTestResourceFunc("r2", applyFunc(2), nil, depFunc(r1))
	r3 := newTestResourceFunc("r3", applyFunc(3), nil, depFunc(r1, r2))
	r4 := newTestResourceFunc("r4", applyFunc(4), nil, depFunc(r3))

	dir := "dirpath"
	registry := NewRegistry(exec.New())
	registry.RegisterResource(r4, dir)
	registry.RegisterResource(r3, dir)
	registry.RegisterResource(r2, dir)
	registry.RegisterResource(r1, dir)

	err := registry.Apply(true)
	assert.Equal(t, 5, order)
	assert.NoError(t, err)
}

func TestApplyInvalidRef(t *testing.T) {
	depFunc := func() []Reference {
		ref := Reference{
			Group: "fakeGroup",
			Kind:  "fakeKind",
			Name:  "fakeResource",
		}
		return []Reference{ref}
	}

	r := newTestResourceFunc("r", nil, nil, depFunc)
	dir := "dirpath"
	registry := NewRegistry(exec.New())
	registry.RegisterResource(r, dir)

	err := registry.Apply(true)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "resource not found with reference"))
}

func TestApplyError(t *testing.T) {
	testCases := []struct {
		name          string
		dryRun        bool
		expectedCalls int
	}{{
		"Apply error, dryRun:false",
		false,
		1,
	}, {
		"Apply error, dryRun: true",
		true,
		2,
	}}
	applyFuncErr := func(expectedDryRun bool) (func(r Registry, dryRun bool) error, *int) {
		applyCalls := 0
		return func(_ Registry, dryRun bool) error {
			applyCalls++
			assert.Equal(t, expectedDryRun, dryRun)
			return fmt.Errorf("Error during apply func")
		}, &applyCalls
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			applyFunc, callCtr := applyFuncErr(tc.dryRun)
			r1 := newTestResourceFunc("r1", applyFunc, nil, nil)
			r2 := newTestResourceFunc("r2", applyFunc, nil, nil)
			dir := "dirpath"
			registry := NewRegistry(exec.New())
			registry.RegisterResource(r1, dir)
			registry.RegisterResource(r2, dir)

			err := registry.Apply(tc.dryRun)
			assert.Error(t, err)
			// Check that dry run is called for both resources and errors are accumulated
			assert.Equal(t, tc.expectedCalls, *callCtr)

			if tc.dryRun {
				merr, ok := err.(*multierror.Error)
				assert.True(t, ok)
				assert.Equal(t, tc.expectedCalls, merr.Len())
			}
		})
	}

}
