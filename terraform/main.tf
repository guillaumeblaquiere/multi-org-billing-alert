provider "google" {
  project = var.runtime_project
  region = var.region
}

provider "google-beta" {
  project = var.runtime_project
  region = var.region
}