terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 3.5.0"
    }
  }
}

provider "google" {
  project     = {{.ProjectID}}
  credentials = {{.CredentialsFile}}
  region      = "us-central1"
  zone        = "us-central1-c"
}
module "gke_auth" {
  source = "terraform-google-modules/kubernetes-engine/google//modules/auth"
  depends_on   = [module.gke]
  project_id   = {{.ProjectID}}
  location     = module.gke.location
  cluster_name = module.gke.name
}
resource "local_file" "kubeconfig" {
  content  = module.gke_auth.kubeconfig_raw
  filename = "kubeconfig-test"
}
module "gcp-network" {
  source       = "terraform-google-modules/network/google"
  version      = "~> 2.5"
  project_id   = {{.ProjectID}}
  network_name = "katana-net"
  subnets = [
    {
      subnet_name   = "katana-net1"
      subnet_ip     = "10.10.0.0/20"
      subnet_region = "us-central1"
    },
  ]
  secondary_ranges = {
    "katana-net1" = [
      {
        range_name    = "ip-range-pods"
        ip_cidr_range = "10.20.0.0/20"
      },
      {
        range_name    = "ip-range-services"
        ip_cidr_range = "10.30.0.0/20"
      },
    ]
  }
}

module "gke" {
  source                 = "terraform-google-modules/kubernetes-engine/google//modules/private-cluster"
  project_id             = {{.ProjectID}}
  name                   = {{.ClusterName}}
  regional               = true
  region                 = "us-central1"
  zones                  = ["us-central1-a", "us-central1-b", "us-central1-f"]
  network                = module.gcp-network.network_name
  subnetwork             = module.gcp-network.subnets_names[0]
  ip_range_pods          = "ip-range-pods"
  ip_range_services      = "ip-range-services"
  node_pools = [
    {
      name                      = "node-pool"
      machine_type              = "e2-medium"
      node_locations            = "us-central1-a"
      min_count                 = 1
      max_count                 = 2
      disk_size_gb              = 100
    },
  ]
}