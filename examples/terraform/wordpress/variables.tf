/**
 * Copyright 2021 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

variable "project_id" {
  description = "The ID of the project in which to provision resources."
  type        = string
}

// Marketplace requires this variable name to be declared
variable "goog_cm_deployment_name" {
  description = "The name of the deployment and VM instance."
  type        = string
}

variable "wp_admin_email" {
  description = "The email address for Wordpress admin."
  type        = string
}

variable "zone" {
  description = "The zone for the solution to be deployed."
  type        = string
  default     = "us-west1-a"
}

variable "source_image" {
  description = "The image name for the disk for the VM instance."
  type        = string
  default     = "projects/click-to-deploy-images/global/images/wordpress-v20220821"
}

variable "machine_type" {
  description = "The machine type to create, e.g. e2-small"
  type        = string
  default     = "n2-standard-4"
}

variable "boot_disk_type" {
  description = "The boot disk type for the VM instance."
  type        = string
  default     = "pd-standard"
}

variable "boot_disk_size" {
  description = "The boot disk size for the VM instance in GBs"
  type        = string
  default     = "20"
}

variable "enable_logging" {
  description = "Enable cloud logging for the VM instance."
  type        = bool
}

variable "enable_monitoring" {
  description = "Enable cloud monitoring for the VM instance."
  type        = bool
}

variable "enable_tcp_80" {
  description = "Enable network traffic over port 80 for Wordpress"
  type        = bool
  default     = false
}

variable "tcp_80_ip_source_ranges" {
  description = "A comma separated string of source IP ranges for accessing the VM instance over HTTP port 80."
  type        = string
  nullable    = true
}

variable "enable_tcp_443" {
  description = "Enable network traffic over port 443 for Wordpress"
  type        = bool
  default     = false
}

variable "tcp_443_ip_source_ranges" {
  description = "A comma separated string of source IP ranges for accessing the VM instance over HTTPS port 443."
  type        = string
  nullable    = true
}

variable "install_phpmyadmin" {
  description = "Install phpMyAdmin on the VM instance"
  type        = bool
  default     = true
}

variable "enable_https" {
  description = "Enabled HTTPS communication for Wordpress"
  type        = bool
  default     = false
}

variable "networks" {
  description = "The network name to attach the VM instance."
  type        = list(string)
  default     = ["default"]
}

variable "sub_networks" {
  description = "The sub network name to attach the VM instance."
  type        = list(string)
  default     = ["default"]
}

variable "external_ips" {
  description = "The external IPs assigned to the VM for public access."
  type        = list(string)
  default     = ["EPHEMERAL"]
}
