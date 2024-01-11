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

# The variable "image" is declared in Producer Portal
variable "image" {
  # Set the default value to your image. Marketplace will overwrite this value
  # to a Marketplace owned image on publishing the product
  default = "projects/<partner-project>/global/images/<image-name>"
}

# Specifies a network interface that this module assumes already exists
variable "other_nic" {
  type = string
}

