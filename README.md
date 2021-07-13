# Overview
This project creates alert billing in the billing account when it receives a message mentioning

* ProjectID(s)
* Monthly Budget
* The email(s) to alert
* *Optionally the custom thresholds*

You can also `GET` the existing alert set on a `ProjectId`, and `DELETE` an existing alert on a `ProjectId`

You can easily extend and customize this sample project. 

This workaround exists because of Google Cloud Billing limitations. See my [articles](https://medium.com/google-cloud/billing-alert-with-cloud-monitoring-notification-channel-c4cfa3588feb) to learn more.

# How to use

The latest built image is present at this location: `us-central1-docker.pkg.dev/gblaquiere-dev/public/multi-org-billing-alert:1.2`

## JSON format

```
{
   "project_id":"<PROJECT_ID>",
   "monthly_budget": FLOAT,
   "emails":[
      "<YOUR_EMAILS>"
   ],
   "thresholds":[
      FLOATS
   ],
   "group_alert":{
      "name":"<ALERT_NAME>",
      "project_ids":[
         "<YOUR_PROJECT_IDS>"
      ]
   }
}
```

## Environment variables

- **BILLING_ACCOUNT**: Billing account number on which to set the alert. *Required*
- **BILLING_PROJECT**: Project with the Cloud Monitoring workspace. If missing or empty, the value is get from the metadata
  server (the current runtime project_id is used)
  
## Deployment

### Service account permission

You need to create a service account which is Billing Admin on your billing account and Notification Channel Editor 
on the billing project which host the Cloud Operation workspace

```
gcloud iam service-accounts create multi-org-billing

gcloud beta billing accounts add-iam-policy-binding <YOUR_BILLING_ACCOUNT> \
   --member=serviceAccount:multi-org-billing@<PROJECT_ID>.iam.gserviceaccount.com \
   --role=roles/billing.admin 
   
gcloud projects add-iam-policy-binding <BILLING_PROJECT_ID> \ 
  --member=serviceAccount:multi-org@<PROJECT_ID>.iam.gserviceaccount.com \ 
  --role=roles/monitoring.notificationChannelEditor
```

### Cloud Run deployment

You it directly in your Cloud Run deployment

```
gcloud run deploy multi-org-billing-alert \
  --image=us-central1-docker.pkg.dev/gblaquiere-dev/public/multi-org-billing-alert:1.1 \
  --region=us-central1 \
  --service-account=serviceAccount:multi-org-billing@<PROJECT_ID>.iam.gserviceaccount.com \
  --platform=managed
  --set-env-vars=BILLING_ACCOUNT=<YOUR_BILLING_ACCOUNT>[,BILLING_PROJECT=<BILLING_PROJECT_ID>]
```

### Command Sample

**Create and Update**

Create and update are the same. The app detects the previous existence of the alert, based on a defined naming convention.
```
# Minimal example
curl -X POST -H "content-type: application/json" -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
  -d '{"project_id": "<PROJECT_ID>","monthly_budget": 10,"emails":["<YOUR_EMAILS>"]}' \
   https://<CLOUD RUN ENDPOINT>/http

# With optional configuable thresholds
curl -X POST -H "content-type: application/json" -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
  -d '{"project_id": "<PROJECT_ID>","monthly_budget": 10,"emails":["<YOUR_EMAILS>"], "thresholds":[0.1,0.5,0.85,1.0]}' \
   https://<CLOUD RUN ENDPOINT>/http
```
*Thresholds are in percent, so `1.0` = 100%*


**Getting and Deleting**

```
# Get an existing budget on a project
curl http://localhost:8080/http/alertname/<PROJECT_ID>

# Delete an existing budget on a project
curl -X DELETE http://localhost:8080/http/alertname/<PROJECT_ID>
```

**Multi projects alert case**

In some situation, you need to create multi-projects alert. You can use the `group_alert` JSON object to fill in your 
project ids list and a name for this alert.

```
# 
curl -X POST -H "content-type: application/json" -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
  -d '{"monthly_budget": 10,"emails":["<YOUR_EMAILS>"], "group_alert":{"name":"<ALERT_NAME>","project_ids":["<YOUR_PROJECT_IDS>"]}}' \
   https://<CLOUD RUN ENDPOINT>/http

```
* *If `project_id` and `group_alert.project_ids` are provided, the list are merged.*
* *If `group_alert.project_ids` is provided without `group_alert.name`, an error is raised*
* *If `group_alert.name` is provided without `group_alert.project_ids`, the name is ignored*


## Automated deployment

This deployment is based on Cloud Build to build the image and terraform to deploy the service and configure the
required elements.

### Terraform configuration

Start by updating terraform configuration in the `terraform.tfvars` file

- **runtime_project**: The project that will host the resources (service account, pubsub, cloud run)
- **billing_project**: The project that will host the notification channel email (for sending emails)
- **billing_account**: The billing account on which to create the budget alerts 
- **members**: The list (`[array]`) of authorized account (user, service or group) to access to Cloud Run and PubSub topic to publish message. The fully qualified account format is required. For example
  * user:user@email.com
  * group:group@email.com
  * serviceAccount:sa@email.com

In the `terraform.tfvars` file, you have several commented variables. Values are set by default on them, but you can
override them if you want. 

### Terraform commands

Use standard commands

```
terraform init
terratorm apply
```

Your current credential need to have access to the runtime project to create resources (service account, pubsub, cloud run,...)

# Run locally

Authenticate yourselves ***(don't use service account key file when you can use your user credentials)***

```
gcloud auth application-default login
```

Be sure to have the following role binding on your user account

* Billing Admin role `roles/billing.admin` on the Billing account that you use 
* Notification Channel Editor `roles/monitoring.notificationChannelEditor` on the billing project

Then run this command
```
BILLING_ACCOUNT=<YOUR_BILLING_ACCOUNT> BILLING_PROJECT=<BILLING_PROJECT_ID> go run .
```

Test the HTTP entry point

```
curl -X POST -H "content-type: application/json" -d '{"project_id": "<PROJECT_ID>","monthly_budget": 10,"emails":["<YOUR_EMAIL>"]}' localhost:8080/http
curl http://localhost:8080/http/projectid/<PROJECT_ID>
curl -X DELETE http://localhost:8080/http/projectid/<PROJECT_ID>
```

# Build

You can also build the image by yourselves (see bellow) with Docker or Cloud Build.

## Docker

Use standard Docker command. Build the image and push it to Google Cloud registry

```
docker build -t <Registry>/multi-org-billing-alert .
docker push <Registry>/multi-org-billing-alert
```

*Registry can be [Container Registry](https://cloud.google.com/container-registry) or [Artifact Registry](https://cloud.google.com/artifact-registry)*

## Cloud Build

```
gcloud builds submit --tag=<Registry>/multi-org-billing-alert
```

*Registry can be [Container Registry](https://cloud.google.com/container-registry) or [Artifact Registry](https://cloud.google.com/artifact-registry)*

## Buildpack

You can use Build Pack to build the image and deploy it immediately on Cloud Run

```
gcloud beta run deploy multi-org-billing-alert \
  --source=. \
  --region=us-central1 \
  --service-account=serviceAccount:multi-org-billing@<PROJECT_ID>.iam.gserviceaccount.com \
  --platform=managed
  --set-env-vars=BILLING_ACCOUNT=<YOUR_BILLING_ACCOUNT>[,BILLING_PROJECT=<BILLING_PROJECT_ID>]
```

# License

This library is licensed under Apache 2.0. Full license text is available in
[LICENSE](https://github.com/guillaumeblaquiere/multi-org-billing-alert/tree/master/LICENSE).