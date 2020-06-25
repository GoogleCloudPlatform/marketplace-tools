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

.PHONY: fix vet fmt license license-check lint test build buildall tidy

GOBIN := $(shell go env GOPATH)/bin
PKG := github.com/GoogleCloudPlatform/marketplace-tools/mpdev

build:
	go build -o $(GOBIN)/mpdev $(PKG)

all: fix vet fmt license license-check lint tidy test build buildall

buildall:
	GOOS=windows go build -o /dev/null $(PKG)
	GOOS=linux go build -o /dev/null $(PKG)
	GOOS=darwin go build -o /dev/null $(PKG)

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
	go test -cover ./...

vet:
	go vet ./...