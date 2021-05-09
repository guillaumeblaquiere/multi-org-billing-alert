resource "google_service_account" "cloud_run" {
  account_id   = var.service_account_name
  project = var.runtime_project
}

resource "google_project_iam_binding" "project" {
  project     = var.billing_project
  role = "roles/monitoring.notificationChannelEditor"
  members = [
    "serviceAccount:${google_service_account.cloud_run.email}",
  ]
}

resource "google_billing_account_iam_binding" "billing_account_admin" {
  provider = google-beta
  billing_account_id = var.billing_account
  role = "roles/billing.admin"
  members = [
    "serviceAccount:${google_service_account.cloud_run.email}",
  ]
}

resource "google_service_account" "pubsub" {
  account_id   = "pubsub-call-cloud-run"
  project = var.runtime_project
}

resource "google_cloud_run_service_iam_binding" "pubsub_call_cloud_run" {
  project     = var.billing_project
  service = google_cloud_run_service.multi_org.name
  location = var.region
  role = "roles/run.invoker"
  members = [
    "serviceAccount:${google_service_account.pubsub.email}",
  ]
}