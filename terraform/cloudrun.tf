resource "google_cloud_run_service" "multi_org" {
  name = var.cloud_run_service_name
  project = var.runtime_project
  location = var.region

  template {
    spec {
      containers {
        image = var.image_tag
        env {
          name = "BILLING_ACCOUNT"
          value = var.billing_account
        }
        env {
          name = "BILLING_PROJECT"
          value = var.billing_project
        }
      }
      service_account_name = google_service_account.cloud_run.email
    }
  }

  traffic {
    percent = 100
    latest_revision = true
  }
  depends_on = [google_project_service.cloudrun]
}

resource "google_cloud_run_service_iam_binding" "binding" {
  location = google_cloud_run_service.multi_org.location
  project = google_cloud_run_service.multi_org.project
  service = google_cloud_run_service.multi_org.name
  role = "roles/run.invoker"
  members  = var.members
}