variable "billing_project" {
  description = "The project that will host the resources (service account, pubsub, cloud run)"
  type        = string
}

variable "runtime_project" {
  description = "The project that will host the notification channel email (for sending emails)"
  type        = string
}

variable "billing_account" {
  description = "The billing account on which to create the budget alerts"
  type        = string
}

variable "members" {
  description = "The list of authorized account (user, service or group) to access to Cloud Run and PubSub topic to publish message. The fully qualified account format is required. For example user:user@email.com or group:group@email.com or serviceAccount:sa@email.com"
  type        = list(string)
  default     = []
}

variable "image_tag" {
  type    = string
  default = "us-central1-docker.pkg.dev/gblaquiere-dev/public/multi-org-billing-alert:1.4"
}

variable "region" {
  type    = string
  default = "us-central1"
}

variable "service_account_name" {
  type    = string
  default = "multi(org-billing"
}

variable "cloud_run_service_name" {
  type    = string
  default = "multi-org-billing-alert"
}

variable "pubsub_topic_name" {
  type    = string
  default = "multi-org-billing-alert-topic"
}

variable "pubsub_subscription_name" {
  type    = string
  default = "multi-org-billing-alert-subscription"
}
