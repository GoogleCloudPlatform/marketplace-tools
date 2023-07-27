# Wordpress

This module deploys a [WordPress Google Click to Deploy Solution](https://console.cloud.google.com/marketplace/product/click-to-deploy-images/wordpress) from Marketplace.

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| boot\_disk\_size | The boot disk size for the VM instance in GBs | `string` | `"10"` | no |
| boot\_disk\_type | The boot disk type for the VM instance. | `string` | `"pd-standard"` | no |
| enable\_https | Enabled HTTPS communication for Wordpress | `bool` | `false` | no |
| enable\_logging | Enable cloud logging for the VM instance. | `bool` | n/a | yes |
| enable\_monitoring | Enable cloud monitoring for the VM instance. | `bool` | n/a | yes |
| enable\_tcp\_443 | Enable network traffic over port 443 for Wordpress | `bool` | `false` | no |
| enable\_tcp\_80 | Enable network traffic over port 80 for Wordpress | `bool` | `false` | no |
| external\_ips | The external IPs assigned to the VM for public access. | `list(string)` | <pre>[<br>  "EPHEMERAL"<br>]</pre> | no |
| goog\_cm\_deployment\_name | The name of the deployment and VM instance. | `string` | n/a | yes |
| install\_phpmyadmin | Install phpMyAdmin on the VM instance | `bool` | `true` | no |
| machine\_type | The machine type to create, e.g. e2-small | `string` | `"n2-standard-4"` | no |
| networks | The network name to attach the VM instance. | `list(string)` | <pre>[<br>  "default"<br>]</pre> | no |
| project\_id | The ID of the project in which to provision resources. | `string` | n/a | yes |
| source\_image | The image name for the disk for the VM instance. | `string` | `"projects/click-to-deploy-images/global/images/wordpress-v20220821"` | no |
| sub\_networks | The sub network name to attach the VM instance. | `list(string)` | <pre>[<br>  "default"<br>]</pre> | no |
| tcp\_443\_ip\_source\_ranges | A comma separated string of source IP ranges for accessing the VM instance over HTTPS port 443. | `string` | n/a | yes |
| tcp\_80\_ip\_source\_ranges | A comma separated string of source IP ranges for accessing the VM instance over HTTP port 80. | `string` | n/a | yes |
| wp\_admin\_email | The email address for Wordpress admin. | `string` | n/a | yes |
| zone | The zone for the solution to be deployed. | `string` | `"us-west1-a"` | no |

## Outputs

| Name | Description |
|------|-------------|
| admin\_password | Password for the admin user |
| admin\_user | Admin username for Wordpress |
| has\_external\_ip | Flag to indicate if the wordpress machine has an external IP |
| instance\_machine\_type | Machine type for the wordpress compute instance |
| instance\_nat\_ip | Machine type for the wordpress compute instance |
| instance\_network | Machine type for the wordpress compute instance |
| instance\_self\_link | Self-link for the Wordpress compute instance |
| instance\_zone | Zone for the wordpress compute instance |
| mysql\_password | Password for the MySql user |
| root\_password | Password for the root user |

