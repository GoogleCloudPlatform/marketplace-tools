# Google Cloud Marketplace Terraform Module

This example meets the Marketplace requirements for a Terraform module with a UI
deployment. The module was generated using the
[Marketplace Autogen](https://github.com/GoogleCloudPlatform/marketplace-autogen)
tool, with [autogen.yaml](./autogen.yaml) used as an input.

> [!IMPORTANT]
> Consider using the [Guided Configuration](https://cloud.google.com/marketplace/docs/partners/vm/configure-terraform-deployment#simple-deployment)
> in Producer Portal instead if you have a simple VM product. This Terraform
> module offers a comprehensive example with advanced configurations, if you're
> interested in creating a customized Terraform module via
> [Manual Configuration](https://cloud.google.com/marketplace/docs/partners/vm/configure-terraform-deployment#complex-deployment)

Marketplace Autogen generates fully functional Terraform modules from simplified
YAML configuration files. It currently supports single-VM and multi-VM
configurations, and can be a good starting point for partners since it also
generates the UI metadata used to actuate the Marketplace deployment UI. The UI
metadata is structured as [`BlueprintMetadata`](https://pkg.go.dev/github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/bpmetadata#BlueprintMetadata).

## Usage
The provided test configuration can be used by executing:

```
terraform plan --var-file marketplace_test.tfvars --var project_id=<YOUR_PROJECT>
```

## Inputs
| Name | Description | Type | Default | Required |
|------|-------------|------|---------|----------|
| project_id | The ID of the project in which to provision resources. | `string` | `null` | yes |
| goog_cm_deployment_name | The name of the deployment and VM instance. | `string` | `null` | yes |
| source_image | The image name for the disk for the VM instance. | `string` | `"projects/click-to-deploy-images/global/images/wordpress-v20240121"` | no |
| zone | The zone for the solution to be deployed. | `string` | `null` | yes |
| machine_type | The machine type to create, e.g. e2-small | `string` | `"e2-small"` | no |
| boot_disk_type | The boot disk type for the VM instance. | `string` | `"pd-standard"` | no |
| boot_disk_size | The boot disk size for the VM instance in GBs | `number` | `20` | no |
| networks | The network name to attach the VM instance. | `list(string)` | `["default"]` | no |
| sub_networks | The sub network name to attach the VM instance. | `list(string)` | `[]` | no |
| external_ips | The external IPs assigned to the VM for public access. | `list(string)` | `["EPHEMERAL"]` | no |
| enable_tcp_80 | Allow HTTP traffic from the Internet | `bool` | `true` | no |
| tcp_80_source_ranges | Source IP ranges for HTTP traffic | `string` | `""` | no |
| enable_tcp_443 | Allow HTTPS traffic from the Internet | `bool` | `true` | no |
| tcp_443_source_ranges | Source IP ranges for HTTPS traffic | `string` | `""` | no |
| adminEmailAddress | The e-mail address used to create the administrator account for WordPress. | `string` | `null` | yes |
| installphpmyadmin | phpMyAdmin is an open source tool to administer MySQL databases with the use of a web browser. | `bool` | `true` | no |
| httpsEnabled | Enabled HTTPS communication. | `bool` | `true` | no |
| enable_cloud_logging | Enables Cloud Logging. | `bool` | `false` | no |
| enable_cloud_monitoring | Enables Cloud Monitoring. | `bool` | `false` | no |

## Outputs

| Name | Description |
|------|-------------|
| site_url | Site Url |
| admin_url | Admin Url |
| wordpress_mysql_user | Username for WordPress MySQL password. |
| wordpress_mysql_password | Password for WordPress MySQL. |
| mysql_root_user | Username for MySQL root password. |
| mysql_root_password | Password for MySQL root. |
| wordpress_admin_user | Username for WordPress Admin password. |
| wordpress_admin_password | Password for WordPress Admin. |
| instance_self_link | Self-link for the compute instance. |
| instance_zone | Zone for the compute instance. |
| instance_machine_type | Machine type for the compute instance. |
| instance_nat_ip | External IP of the compute instance. |
| instance_network | Self-link for the network of the compute instance. |

## Requirements
### Terraform

Be sure you have the correct Terraform version (1.2.0+), you can choose the binary here:

https://releases.hashicorp.com/terraform/

### Configure a Service Account
In order to execute this module you must have a Service Account with the following project roles:

- `roles/compute.admin`
- `roles/iam.serviceAccountUser`

If you are using a shared VPC:

- `roles/compute.networkAdmin` is required on the Shared VPC host project.

### Enable API
In order to operate with the Service Account you must activate the following APIs on the project where the Service Account was created:

- Compute Engine API - `compute.googleapis.com`
