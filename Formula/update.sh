#!/bin/bash
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

# Updates mpdev.rb formula sha256 and url
# Usage:
# Update formula to latest remote tag
# ./update.sh
# Update formula to specific tag
# ./update.sh v0.1.0
set -eux;

cleanup() {
  rm -f "${RELEASE}"
}

cd "$(dirname "$0")"
REPO=https://github.com/GoogleCloudPlatform/marketplace-tools;

TAG=$(git ls-remote --tags --sort="v:refname" "$REPO" | tail -n1 | sed -e 's|^.*refs/tags/||' -e 's|\^{}||');
TAG=${1-$TAG}
RELEASE=mpdev_darwin_amd64_${TAG}.tar.gz;
URL="${REPO}/releases/download/${TAG}/${RELEASE}"
wget "$URL";
trap cleanup EXIT;

SHA=$(sha256sum ${RELEASE} | awk '{print $1}');
sed -i -e "s|sha256.*|sha256 \"$SHA\"|" -e "s|url.*|url \"$URL\"|" ./mpdev.rb


