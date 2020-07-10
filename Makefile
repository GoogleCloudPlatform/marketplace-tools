# Copyright 2020 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

.PHONY: fix vet fmt license license-check lint bazel-build-gen tidy build test

GOBIN := $(shell go env GOPATH)/bin
PKG := github.com/GoogleCloudPlatform/marketplace-tools/mpdev

build:
	bazel build //...:all
	# Tags docker image
	bazel run //mpdev/autogen:docker_image -- --norun

all: fix vet fmt license license-check lint bazel-build-gen tidy build test

fix:
	go fix ./...

fmt:
	go fmt ./...

tidy:
	go mod tidy

lint:
	( [ -f $(GOBIN)/golangci-lint ] || go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.27.0)
	$(GOBIN)/golangci-lint run ./...

license:
	( [ -f $(GOBIN)/addlicense ] || go get github.com/google/addlicense)
	$(GOBIN)/addlicense -y 2020 -l apache *

license-check:
	( [ -f $(GOBIN)/go-licenses ] || go get github.com/google/go-licenses)
	$(GOBIN)/go-licenses check $(PKG)

test:
	bazel test //... --test_output=errors

bazel-build-gen:
	bazel run :gazelle -- update-repos -from_file=go.mod -build_file_proto_mode disable --to_macro=repos.bzl%go_repositories --prune
	bazel run :gazelle

vet:
	go vet ./...
