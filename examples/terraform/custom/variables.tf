# The variable "image" is declared in Producer Portal
variable "image" {
  # Set the default value to your image. Marketplace will overwrite this value
  # to a Marketplace owned image on publishing the product
  default = "projects/<partner-project>/global/images/<image-name>"
}
