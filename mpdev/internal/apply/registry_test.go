package apply

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveFilePath(t *testing.T) {
	r := NewTestResource("testResource")

	testCases := map[string]string{
		// Absolute path
		"/tmp/foo/2.txt": "/tmp/foo/2.txt",
		// Relative path
		"foo/2.txt": "dir/foo/2.txt",
	}

	for path, expected := range testCases {
		registry := NewRegistry()
		registry.RegisterResource(r, "dir")
		resolvedPath := registry.ResolveFilePath(r, path)
		assert.Equal(t, expected, resolvedPath)
	}
}

func TestApplyOrder(t *testing.T) {
	order := 1
	applyFunc := func(expectedOrder int) func(r Registry) error {
		return func(r Registry) error {
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

	r1 := NewTestResourceFunc("r1", applyFunc(1), nil)
	r2 := NewTestResourceFunc("r2", applyFunc(2), depFunc(r1))
	r3 := NewTestResourceFunc("r3", applyFunc(3), depFunc(r1, r2))
	r4 := NewTestResourceFunc("r4", applyFunc(4), depFunc(r3))

	dir := "dirpath"
	registry := NewRegistry()
	registry.RegisterResource(r4, dir)
	registry.RegisterResource(r3, dir)
	registry.RegisterResource(r2, dir)
	registry.RegisterResource(r1, dir)

	err := registry.Apply()
	assert.Equal(t, 5, order)
	assert.NoError(t, err)
}