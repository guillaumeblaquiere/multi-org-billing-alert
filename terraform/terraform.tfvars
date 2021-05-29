// REQUIRED variables

// The project that will host the resources (service account, pubsub, cloud run)
runtime_project=""
// The project that will host the notification channel email (for sending emails)
billing_project=""
// The billing account on which to create the budget alerts
billing_account=""
// The list of authorized account (user, service or group) to access to Cloud Run and PubSub topic to publish message.
// The fully qualified account format is required. For example
//  * user:user@email.com
//  * group:group@email.com
//  * serviceAccount:sa@email.com
members=[]


// OPTIONAL variables
//region="us-central1"
//service_account_name="multi-org-billing"
//cloud_run_service_name="multi-org-billing-alert"
//pubsub_topic_name="multi-org-billing-alert-topic"
//pubsub_subscription_name="multi-org-billing-alert-subscription"
//image_tag="us-central1-docker.pkg.dev/gblaquiere-dev/public/multi-org-billing-alert"