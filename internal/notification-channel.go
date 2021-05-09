package internal

import (
	"cloud.google.com/go/compute/metadata"
	monitoringApi "cloud.google.com/go/monitoring/apiv3"
	"context"
	"fmt"
	"gblaquiere.dev/multi-org-billing-alert/model"
	"google.golang.org/api/iterator"
	"google.golang.org/genproto/googleapis/monitoring/v3"
	"log"
	"os"
)

func getChannelIDs(ctx context.Context, message *model.BillingAlert) (err error) {

	billingProject, err := getBillingProject()
	if err != nil {
		log.Printf("getBillingProject(): %+v\n", err)
		return err
	}

	client, err := monitoringApi.NewNotificationChannelClient(ctx)
	if err != nil {
		log.Printf("monitoringApi.NewNotificationChannelClient: %+v\n", err)
		return err
	}

	//Create filter
	filter := ""
	for i, v := range message.Emails {
		if i > 0 {
			filter += " OR "
		}
		filter += fmt.Sprintf("labels.email_address='%s'", v)
	}

	req := &monitoring.ListNotificationChannelsRequest{
		Name:   billingProject,
		Filter: filter,
	}
	notificationChannels := client.ListNotificationChannels(ctx, req)

	var existingNotificationChannels []*monitoring.NotificationChannel
	for {
		notificationChannel, err := notificationChannels.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("notificationChannels.Next: %+v\n", err)
			return err
		}
		existingNotificationChannels = append(existingNotificationChannels, notificationChannel)
	}

	//Create missing Notification channel
	err = createMissingNotificationChannel(ctx, client, message, &existingNotificationChannels, billingProject)
	fmt.Printf("%+v\n", existingNotificationChannels)

	if err != nil {
		return err
	}

	//Add to message
	for _, notificationChannel := range existingNotificationChannels {
		message.ChannelIds = append(message.ChannelIds, notificationChannel.GetName())
	}

	return
}

func getBillingProject() (billingProject string, err error) {
	billingProject = os.Getenv("BILLING_PROJECT")

	if billingProject == "" { // get the value in the metadata
		billingProject, err = metadata.Get("/project/project-id")
	}
	billingProject = fmt.Sprintf("projects/%s", billingProject)
	return
}

func createMissingNotificationChannel(ctx context.Context, client *monitoringApi.NotificationChannelClient, message *model.BillingAlert, notificationChannels *[]*monitoring.NotificationChannel, project string) (err error) {
	for _, email := range message.Emails {
		found := false
		for _, notificationChannel := range *notificationChannels {
			if email == notificationChannel.Labels["email_address"] {
				found = true
			}
		}
		if !found {
			req := &monitoring.CreateNotificationChannelRequest{
				NotificationChannel: &monitoring.NotificationChannel{
					Type:        "email",
					DisplayName: fmt.Sprintf("alert-billing-%s", email),
					Labels: map[string]string{
						"email_address": email,
					},
				},
				Name: project,
			}
			nc, err := client.CreateNotificationChannel(ctx, req)
			if err != nil {
				log.Printf("client.CreateNotificationChannel: %+v\n", err)
				return err
			}
			*notificationChannels = append(*notificationChannels, nc)
		}
	}
	return
}
