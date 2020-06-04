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

FROM golang:1.13-stretch
ENV CGO_ENABLED=0
WORKDIR /go/src/
RUN apt-get update && apt-get install unzip

COPY go.mod go.sum ./
RUN go mod download

COPY scripts .
COPY deploymentmanager-autogen deploymentmanager-autogen
RUN ./install_protogen.sh && ./run_protogen.sh

COPY . .
RUN go build -v -o /usr/local/bin/autogen

FROM alpine:latest
COPY --from=0 /usr/local/bin/autogen /usr/local/bin/autogen
CMD ["autogen"]