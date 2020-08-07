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

# Homebrew formula for installing mpdev. See ./update.sh for upgrading to newest release
class Mpdev < Formula
  desc "Toolkit to both configure and construct artifacts needed for publishing to the Google Cloud Marketplace."
  homepage "https://github.com/GoogleCloudPlatform/marketplace-tools"
  url "https://github.com/GoogleCloudPlatform/marketplace-tools/releases/download/v0.1.0/mpdev_darwin_amd64_v0.1.0.tar.gz"
  sha256 "3cf70d46ad0264a786126f043104e132cd1f498a25413a030c14c290c6003ec9"

  bottle :unneeded

  def install
    bin.install "mpdev"
  end

  test do
    system "#{bin}/mpdev", version
  end
end
