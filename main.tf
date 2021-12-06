terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "3.5.0"
    }
  }
}

provider "google" {
  project     = "tensile-site-325304"
  credentials = "/home/rohithvb/Downloads/tensile-site-325304-d2945003c5df.json"
  region      = "us-central1"
  zone        = "us-central1-c"
}

resource "google_container_cluster" "cluster" {
  name               = "katana-cluster"
  location           = "us-central1-a"
  initial_node_count = 1
  node_config {
    machine_type = "e2-medium"
    disk_size_gb = 100
  }

}