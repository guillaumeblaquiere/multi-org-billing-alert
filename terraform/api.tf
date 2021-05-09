resource "google_project_service" "billing" {
  service = "billingbudgets.googleapis.com"
  project = var.runtime_project
  disable_on_destroy = true
}
resource "google_project_service" "iam" {
  service = "iam.googleapis.com"
  project = var.runtime_project
  disable_on_destroy = true
}
resource "google_project_service" "cloudrun" {
  service = "run.googleapis.com"
  project = var.runtime_project
  disable_on_destroy = true
}
resource "google_project_service" "monitoring" {
  service = "monitoring.googleapis.com"
  project = var.runtime_project
  disable_on_destroy = true
}
resource "google_project_service" "pubsub" {
  service = "pubsub.googleapis.com"
  project = var.runtime_project
  disable_on_destroy = true
}