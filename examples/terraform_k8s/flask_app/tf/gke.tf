locals {
  endpoint       = var.create_cluster ? "https://${module.gke[0].endpoint}" : "https://${data.google_container_cluster.default[0].endpoint}"
  ca_certificate = var.create_cluster ? base64decode(module.gke[0].ca_certificate) : base64decode(data.google_container_cluster.default[0].master_auth[0].cluster_ca_certificate)
  host           = local.endpoint
  cluster_name   = var.goog_cm_deployment_name != "" ? "${var.goog_cm_deployment_name}-${var.cluster_name}" : var.cluster_name
}

locals {
  gpu_l4_t4_location = {
    asia-east1      = "asia-east1-a,asia-east1-c"
    asia-northeast1 = "asia-northeast1-a,asia-northeast1-c"
    asia-northeast3 = "asia-northeast3-b"
    asia-south1     = "asia-south1-a,asia-south1-b"
    asia-southeast1 = "asia-southeast1-a,asia-southeast1-b,asia-southeast1-c"
    europe-west1    = "europe-west1-b,europe-west1-c"
    europe-west2    = "europe-west2-a,europe-west2-b"
    europe-west3    = "europe-west3-b"
    europe-west4    = "europe-west4-a,europe-west4-b,europe-west4-c"
    us-central1     = "us-central1-a,us-central1-b,us-central1-c"
    us-east1        = "us-east1-c,us-east1-d"
    us-east4        = "us-east4-a,us-east4-c"
    us-west1        = "us-west1-a,us-west1-b"
    us-west4        = "us-west4-a"
  }
}

locals {
  region     = length(split("-", var.cluster_location)) == 2 ? var.cluster_location : ""
  regional   = local.region != "" ? true : false
  zone       = length(split("-", var.cluster_location)) > 2 ? split(",", var.cluster_location) : split(",", local.gpu_l4_t4_location[local.region])
  gpu_pools  = [for elm in var.gpu_pools : (local.regional && contains(keys(local.gpu_l4_t4_location), local.region) && elm["node_locations"] == "") ? merge(elm, { "node_locations" : local.gpu_l4_t4_location[local.region] }) : elm]
  node_pools = concat((var.enable_gpu ? var.gpu_pools : []), (var.enable_tpu ? var.tpu_pools : []), var.cpu_pools)
}

provider "kubernetes" {
  alias                  = "app"
  host                   = local.host
  token                  = data.google_client_config.default.access_token
  cluster_ca_certificate = local.ca_certificate
}

data "google_container_cluster" "default" {
  count      = var.create_cluster ? 0 : 1
  name       = var.cluster_name
  location   = var.cluster_location
  depends_on = [module.project-services]
}

data "google_compute_network" "existing-network" {
  name    = var.network_name
  project = var.project_id
}

data "google_compute_subnetwork" "subnetwork" {
  name    = var.subnetwork_name
  region  = var.subnetwork_region
  project = var.project_id
}

module "gke" {
  count   = var.create_cluster ? 1 : 0
  source  = "terraform-google-modules/kubernetes-engine/google"
  version = "~> 27.0"

  project_id             = var.project_id
  create_service_account = var.create_cluster_service_account
  service_account        = var.cluster_service_account

  name               = var.cluster_name
  region             = local.region
  regional           = local.regional
  zones              = local.zone
  kubernetes_version = var.kubernetes_version

  remove_default_node_pool = true
  initial_node_count       = 1

  network           = var.network_name
  subnetwork        = var.subnetwork_name
  ip_range_pods     = var.ip_range_pods
  ip_range_services = var.ip_range_services

  issue_client_certificate = true

  node_pools = local.node_pools

  node_pools_oauth_scopes = {
    all = ["https://www.googleapis.com/auth/cloud-platform"]
  }
}
