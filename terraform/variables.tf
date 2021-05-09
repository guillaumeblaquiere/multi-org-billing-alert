variable "billing_project" {
  description = "The project that will host the resources (service account, pubsub, cloud run)"
}
variable "runtime_project" {
  description = "The project that will host the notification channel email (for sending emails)"
}
variable "billing_account" {
  description = "The billing account on which to create the budget alerts"
}


variable "image_tag" {
  default= "us-central1-docker.pkg.dev/gblaquiere-dev/public/multi-org-billing-alert"
}
variable "region" {
  default = "us-central1"
}
variable "service_account_name" {
  default = "multi(org-billing"
}
variable "cloud_run_service_name" {
  default = "multi-org-billing-alert"
}
variable "pubsub_topic_name" {
  default = "multi-org-billing-alert-topic"
}
variable "pubsub_subscription_name" {
  default = "multi-org-billing-alert-subscription"
}