provider "helm" {
  alias = "app"
  kubernetes {
    host                   = local.host
    token                  = data.google_client_config.default.access_token
    cluster_ca_certificate = local.ca_certificate
  }
}

resource "helm_release" "flask_demo_app" {
  provider = helm.app

  name = "flask_demo_app"

  repository = var.helm_chart_repo
  chart      = var.helm_chart_name
  version    = var.helm_chart_version

  set {
    name = "image.repository"
    value = var.image_repo
  }

  set {
    name = "image.tag"
    value = var.image_tag
  }

  set {
    name = "replicaCount"
    value = var.replica_count
  }
}
