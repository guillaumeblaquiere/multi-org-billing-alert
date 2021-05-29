resource "google_pubsub_topic" "multi_org_billing" {
  name = var.pubsub_topic_name
  project = var.runtime_project
}



resource "google_pubsub_subscription" "multi_org_billing" {
  name  = var.pubsub_subscription_name
  project = var.runtime_project
  topic = google_pubsub_topic.multi_org_billing.name

  push_config {
    push_endpoint = "${google_cloud_run_service.multi_org.status[0].url}/pubsub"
    oidc_token {
      service_account_email = google_service_account.pubsub.email
      audience = google_cloud_run_service.multi_org.status[0].url
    }
  }
}

resource "google_pubsub_topic_iam_binding" "binding" {
  project = google_pubsub_topic.multi_org_billing.project
  topic = google_pubsub_topic.multi_org_billing.name
  role = "roles/pubsub.publisher"
  members  = var.members
}

