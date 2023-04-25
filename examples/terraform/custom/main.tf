resource "google_compute_instance" "default" {
  name         = "example"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
       # The boot disk must be set to the variable declared in Producer Portal
       image = var.image
    }
  }

  network_interface {
    network = "default"
  }
}

