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

locals {
  network_interface = google_compute_instance.instance.network_interface[0]
  instance_nat_ip   = length(local.network_interface.access_config) > 0 ? local.network_interface.access_config[0].nat_ip : null
  instance_ip       = coalesce(local.instance_nat_ip, local.network_interface.network_ip)
}

output "site_url" {
  description = "Site Url"
  value       = "https://${local.instance_ip}/"
}

output "admin_url" {
  description = "Admin Url"
  value       = "https://${local.instance_ip}/wp-admin"
}

output "wordpress_mysql_user" {
  description = "Username for WordPress MySQL password."
  value       = "wordpress"
}

output "wordpress_mysql_password" {
  description = "Password for WordPress MySQL."
  value       = random_password.wordpress_mysql.result
  sensitive   = true
}

output "mysql_root_user" {
  description = "Username for MySQL root password."
  value       = "root"
}

output "mysql_root_password" {
  description = "Password for MySQL root."
  value       = random_password.mysql_root.result
  sensitive   = true
}

output "wordpress_admin_user" {
  description = "Username for WordPress Admin password."
  value       = var.adminEmailAddress
}

output "wordpress_admin_password" {
  description = "Password for WordPress Admin."
  value       = random_password.wordpress_admin.result
  sensitive   = true
}

output "instance_self_link" {
  description = "Self-link for the compute instance."
  value       = google_compute_instance.instance.self_link
}

output "instance_zone" {
  description = "Zone for the compute instance."
  value       = var.zone
}

output "instance_machine_type" {
  description = "Machine type for the compute instance."
  value       = var.machine_type
}

output "instance_nat_ip" {
  description = "External IP of the compute instance."
  value       = local.instance_nat_ip
}

output "instance_network" {
  description = "Self-link for the network of the compute instance."
  value       = var.networks[0]
}
