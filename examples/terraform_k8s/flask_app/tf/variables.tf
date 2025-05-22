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

variable "project_id" {
  type        = string
  description = "GCP project id"
}

variable "create_cluster_service_account" {
  type    = bool
  default = false
}

variable "cluster_service_account" {
  type    = string
  default = ""
}

# Helm

variable "helm_release_name" {
  type = string
}

variable "helm_chart_repo" {
  type = string
}

variable "helm_chart_name" {
  type = string
}

variable "helm_chart_version" {
  type = string
}

# GKE

variable "kubernetes_version" {
  type    = string
  default = "1.30"
}

variable "cluster_name" {
  type = string
}

variable "cluster_location" {
  type = string
}

variable "create_cluster" {
  type    = bool
  default = false
}

variable "cpu_pools" {
  type = list(map(any))
  default = [{
    name         = "cpu-pool"
    machine_type = "n1-standard-16"
    autoscaling  = true
    min_count    = 1
    max_count    = 3
    disk_size_gb = 100
    disk_type    = "pd-standard"
  }]
}

variable "enable_gpu" {
  type        = bool
  description = "Set to true to create GPU node pool"
  default     = false
}

variable "gpu_pools" {
  type    = list(map(any))
  default = []
}

variable "enable_tpu" {
  type        = bool
  description = "Set to true to create TPU node pool"
  default     = false
}

variable "tpu_pools" {
  type    = list(map(any))
  default = []
}

variable "ip_range_pods" {
  type    = string
  default = ""
}

variable "ip_range_services" {
  type    = string
  default = ""
}

variable "network_name" {
  type    = string
  default = "default"
}

variable "subnetwork_name" {
  type    = string
  default = "default"
}

variable "subnetwork_region" {
  type    = string
  default = "us-central1"
}

# TODO: Add variables for your app, including images and other customizable values
variable "image_repo" {
  type = string
}

variable "image_tag" {
  type = string
}

variable "replica_count" {
  type = number
}
