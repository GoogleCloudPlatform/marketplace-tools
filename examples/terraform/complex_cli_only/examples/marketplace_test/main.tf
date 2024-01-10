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

provider "google" {
  project = var.project_id
}

# The test module references the module in the root directory
module "test" {
  source = "../.."

  other_nic = google_compute_network.new_nic.name
}

# Create a network interface in the test module. The Marketplace validation project
# may not have a resource that your module in the root directory expects to exist.
resource "google_compute_network" "new_nic" {
  name = "new-test-nic"
}

# The module must declare a project variable that Marketplace can
# set for validation
variable "project_id" {
  type = string
}

