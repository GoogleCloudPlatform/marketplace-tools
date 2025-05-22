# Copyright 2024 Google LLC
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

provider "helm" {
  alias = "app"
  kubernetes {
    host                   = local.host
    token                  = data.google_client_config.default.access_token
    cluster_ca_certificate = local.ca_certificate
  }
}

locals {
  helm_release_name = var.helm_release_name != "" ? var.helm_release_name : "flask-app-${random_string.suffix.result}"
}

resource "helm_release" "primary" {
  provider = helm.app

  name = local.helm_release_name

  repository = var.helm_chart_repo
  chart      = var.helm_chart_name
  version    = var.helm_chart_version

  set {
    name  = "image.repository"
    value = var.image_repo
  }

  set {
    name  = "image.tag"
    value = var.image_tag
  }

  set {
    name  = "replicaCount"
    value = var.replica_count
  }
}
