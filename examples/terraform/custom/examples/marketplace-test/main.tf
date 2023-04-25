provider "google" {
  project = var.project
}

# The test module references the module in the root directory
module "test" {
  source = "../.."
}

# The module must declare a project variable that Marketplace can
# set for validation
variable "project" {
  type = string
}

