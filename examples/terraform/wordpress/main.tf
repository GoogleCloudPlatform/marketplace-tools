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
provider "google" {
  project = var.project_id
}

locals {
  network_interfaces_map = { for i, n in var.networks : n => {
    network     = n,
    subnetwork  = element(var.sub_networks, i)
    external_ip = element(var.external_ips, i)
    }
  }

  tcp_80_ip_source_ranges_map = var.enable_tcp_80 ? {
    "80" = split(var.tcp_80_ip_source_ranges, ",")
  } : {}

  tcp_443_ip_source_ranges_map = var.enable_tcp_443 ? {
    "443" = split(var.tcp_443_ip_source_ranges, ",")
  } : {}

  ip_source_ranges_map = merge(local.tcp_80_ip_source_ranges_map, local.tcp_443_ip_source_ranges_map)
}

resource "random_password" "mysql" {
  length  = 8
  special = false
}

resource "random_password" "mysql_root" {
  length  = 14
  special = false
}

resource "random_password" "mysql_admin" {
  length           = 8
  special          = true
  override_special = "-"
}

resource "google_compute_instance" "default" {
  name         = var.goog_cm_deployment_name
  machine_type = var.machine_type
  zone         = var.zone

  metadata = {
    enable-oslogin           = "TRUE"
    google-logging-enable    = var.enable_logging == true ? "1" : "0"
    google-monitoring-enable = var.enable_monitoring == true ? "1" : "0"
    installphpmyadmin        = var.install_phpmyadmin == true ? "True" : "False"
    wordpress-admin-email    = var.wp_admin_email
    wordpress-enable-https   = var.enable_https == true ? "True" : "False"
    wordpress-mysql-password = random_password.mysql.result
    mysql-root-password      = random_password.mysql_root.result
    wordpress-admin-password = random_password.mysql_admin.result
  }

  boot_disk {
    initialize_params {
      size  = var.boot_disk_size
      type  = var.boot_disk_type
      image = var.source_image
    }
  }

  dynamic "network_interface" {
    for_each = local.network_interfaces_map
    content {
      network    = network_interface.key
      subnetwork = network_interface.value.subnetwork

      dynamic "access_config" {
        for_each = compact([network_interface.value.external_ip == "None" ? null : 1])
        content {
          nat_ip = network_interface.value.external_ip == "Ephemeral" ? null : network_interface.value.external_ip
        }
      }
    }
  }

  tags = var.tcp_80_ip_source_ranges != "" || var.tcp_443_ip_source_ranges != "" ? ["${var.goog_cm_deployment_name}-deployment"] : []
  service_account {
    scopes = [
      "https://www.googleapis.com/auth/cloud.useraccounts.readonly",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring.write",
      "https://www.googleapis.com/auth/cloudruntimeconfig",
    ]
  }
}

resource "google_compute_firewall" "http" {
  for_each = local.ip_source_ranges_map
  project  = var.project_id
  name     = "${var.goog_cm_deployment_name}-tcp-${each.key}"
  network  = element(var.networks, 0)

  allow {
    protocol = "tcp"
    ports    = [each.key]
  }

  source_ranges = each.value
  target_tags   = ["${var.goog_cm_deployment_name}-deployment"]
}
