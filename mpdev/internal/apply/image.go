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

// PackerGceImageBuilder uses Packer to create a GCEImage when applied
type PackerGceImageBuilder struct {
	BaseResource
	Builder struct {
		Script struct {
			File string
		}
	}

	Tests []struct {
		Name   string
		Script struct {
			File string
		}
	}
}

// Apply build a GCE image using Packer
func (p *PackerGceImageBuilder) Apply(registry *registry) error {
	return nil
}

// GceImage represents a Google Compute Engine image. One of BuilderRef or
// ImageRef must be specified
type GceImage struct {
	BaseResource

	// References another GCE Image resource
	ImageRef Reference

	// References a builder resource which handles the actual creation
	// of the GCE Image
	BuilderRef Reference
	Image      Image
}

// Apply publishes an image to with the project and name specified in
// the GceImage
func (g *GceImage) Apply(registry *registry) error {
	return nil
}

// Image defines the location of the GCE Image when published
type Image struct {
	ProjectID          string `json:"projectId"`
	NamePartsSeparator string
	NameParts          []string
}
