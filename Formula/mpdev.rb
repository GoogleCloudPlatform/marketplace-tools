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

# Homebrew formula for installing mpdev
class Mpdev < Formula
  desc "Toolkit to both configure and construct artifacts needed for publishing to the Google Cloud Marketplace."
  homepage "https://github.com/GoogleCloudPlatform/marketplace-tools"
  url "https://github.com/GoogleCloudPlatform/marketplace-tools/releases/download/v0.0.0/mpdev_darwin_amd64_v0.0.0.tar.gz"
  sha256 "2dc3bf838b10b0f386e3b0151fd3b9868a31de16d695f42b6ae97887c1c38646"

  bottle :unneeded

  def install
    bin.install "mpdev"
  end

  test do
    system "#{bin}/mpdev", version
  end
end
