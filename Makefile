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

build:
	go build -o $(GOBIN)/mpdev -v ./mpdev

all: fix vet fmt license license-check lint test build buildall tidy

buildall:
	GOOS=windows go build -o /dev/null
	GOOS=linux go build -o /dev/null
	GOOS=darwin go build -o /dev/null

fix:
	go fix ./...

fmt:
	go fmt ./...

tidy:
	go mod tidy

lint:
	(which golangci-lint || go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.22.2)
	$(GOBIN)/golangci-lint run ./...

license:
	(which addlicense || go get github.com/google/addlicense)
	$(GOBIN)/addlicense -y 2020 -l apache *

license-check:
	(which go-licenses || go get github.com/google/go-licenses)
	$(GOBIN)/go-licenses check github.com/GoogleCloudPlatform/marketplace-tools

test:
	go test -cover ./...

vet:
	go vet ./...
