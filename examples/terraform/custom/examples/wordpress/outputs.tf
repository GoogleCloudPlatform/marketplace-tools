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
output "instance_self_link" {
  description = "Self-link for the Wordpress compute instance"
  value       = google_compute_instance.default.self_link
}

output "instance_zone" {
  description = "Zone for the wordpress compute instance"
  value       = var.zone
}

output "instance_machine_type" {
  description = "Machine type for the wordpress compute instance"
  value       = var.machine_type
}

output "instance_nat_ip" {
  description = "Machine type for the wordpress compute instance"
  value       = google_compute_instance.default.network_interface[0].access_config[0].nat_ip
}

output "has_external_ip" {
  description = "Flag to indicate if the wordpress machine has an external IP"
  value       = google_compute_instance.default.network_interface[0].access_config[0].nat_ip != null ? true : false
}

output "instance_network" {
  description = "Machine type for the wordpress compute instance"
  value       = google_compute_instance.default.network_interface[0]
}

output "mysql_password" {
  description = "Password for the MySql user"
  value       = random_password.mysql.result
  sensitive   = true
}

output "root_password" {
  description = "Password for the root user"
  value       = random_password.mysql_root.result
  sensitive   = true
}

output "admin_password" {
  description = "Password for the admin user"
  value       = random_password.mysql_admin.result
  sensitive   = true
}

output "admin_user" {
  description = "Admin username for Wordpress"
  value       = var.wp_admin_email
}
