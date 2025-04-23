provider "google" {
  project = var.project_id
}

provider "google-beta" {
  project = var.project_id
}

data "google_client_config" "default" {}

data "google_project" "project" {
  project_id = var.project_id
}

## Enable Required GCP Project Services APIs
module "project-services" {
  source  = "terraform-google-modules/project-factory/google//modules/project_services"
  version = "~> 17.0"

  project_id                  = var.project_id
  disable_services_on_destroy = false
  disable_dependent_services  = false
  activate_apis = flatten([
    "cloudresourcemanager.googleapis.com",
    "compute.googleapis.com",
    "config.googleapis.com",
    "container.googleapis.com",
    # TODO: Add your app's required APIs here
  ])
}
