// REQUIRED variables

// The project that will host the resources (service account, pubsub, cloud run)
runtime_project="test-terraform-313212"
// The project that will host the notification channel email (for sending emails)
billing_project="test-terraform-313212"
// The billing account on which to create the budget alerts
billing_account="014FBB-2A1336-855F7A"

// OPTIONAL variables
//region="us-central1"
//service_account_name="multi-org-billing"
//cloud_run_service_name="multi-org-billing-alert"
//pubsub_topic_name="multi-org-billing-alert-topic"
//pubsub_subscription_name="multi-org-billing-alert-subscription"
//image_tag="us-central1-docker.pkg.dev/gblaquiere-dev/public/multi-org-billing-alert"