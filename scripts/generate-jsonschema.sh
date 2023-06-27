#!/bin/bash
# Copyright 2023 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -eu

cd "$(dirname "${BASH_SOURCE[0]}")"

GO111MODULE=on go get github.com/chrusty/protoc-gen-jsonschema/cmd/protoc-gen-jsonschema

json_schema_dir=../jsonschema
autogen_dir="$(bazel info output_base)"/external/deploymentmanager-autogen

protoc --jsonschema_out=json_fieldnames,disallow_additional_properties:$json_schema_dir --proto_path=$autogen_dir --proto_path=$PWD resource.proto

# This jq command does the following
# 1. Fixes the references to definitions in the json schema. This appears to
#    be a compatibility issue in the protoc-gen-jsonschema tool. We change
#    {"$ref": "cloud.deploymentmanager.<Resource>"} to
#    {"$ref": "#/definitions/cloud.deploymentmanager.<Resource>"}
# 2. Restricts apiVersion, kind, and group fields to enum values. This ensures
#    that auto-complete will only suggest values supported by mpdev.
jq 'walk(if type == "object" and .["$ref"] then .["$ref"] ="#/definitions/\(.["$ref"])" else . end) |
    .properties.apiVersion.enum = ["dev.marketplace.cloud.google.com/v1alpha1"] |
    .properties.deploymentManagerRef.properties.group.enum = ["dev.marketplace.cloud.google.com"] |
    .properties.kind.enum = ["DeploymentManagerTemplate", "DeploymentManagerAutogenTemplate"] |
    .properties.deploymentManagerRef.properties.kind.enum = ["DeploymentManagerAutogenTemplate"]
    '   $json_schema_dir/Resource.jsonschema > $json_schema_dir/mpdev.jsonschema

rm $json_schema_dir/Resource.jsonschema
