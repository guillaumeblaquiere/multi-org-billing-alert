package internal

import (
	"context"
	"errors"
	"gblaquiere.dev/multi-org-billing-alert/internal/billingAlertApi"
	"gblaquiere.dev/multi-org-billing-alert/internal/httperrors"
	"gblaquiere.dev/multi-org-billing-alert/internal/notificationChannelApi"
	"gblaquiere.dev/multi-org-billing-alert/model"
	"log"
	"net/http"
)

func CreateBillingAlert(ctx context.Context, message *model.BillingAlert) (billingAlert *model.BillingAlert, err error) {

	//clean up
	if message.GroupAlert != nil && len(message.GroupAlert.ProjectIds) == 0 {
		message.GroupAlert.AlertName = ""
	}

	errMessage := ""
	if message.GroupAlert != nil && len(message.GroupAlert.ProjectIds) > 0 && message.GroupAlert.AlertName == "" {
		errMessage += "List of project IDs provided without alert name\n"
	}

	// If no group of projectIds, clean the potential alert Name
	if message.ProjectID == "" && message.GroupAlert == nil {
		errMessage += "projectid: ProjectId can't be null or empty\n"
	}

	if len(message.Emails) == 0 {
		errMessage += "emails: Notification email list must contain at least one email to notify\n"

	}
	if len(message.Emails) > 5 {
		errMessage += "emails: Notification email list can contain up to 5 emails, no more. Use groups instead of individual emails.\n"

	}
	if message.MonthlyBudget <= 0 {
		errMessage += "monthly_budget: Monthly budget must be positive float.\n"

	}

	if errMessage != "" {
		return nil, httperrors.New(errors.New(errMessage), http.StatusBadRequest)
	}

	//Get the channelIDs list for all the email
	err = notificationChannelApi.GetChannelIDs(ctx, message)
	if err != nil {
		return
	}

	//Create or update the alert budgetApi
	billingAlert, err = billingAlertApi.UpsertBillingAlert(ctx, message)
	return
}

func GetBillingAlert(ctx context.Context, projectId string) (billingAlert *model.BillingAlert, err error) {
	if projectId == "" {
		log.Printf("no project ID parameter\n")
		err = httperrors.New(errors.New("projectId not provided"), http.StatusBadRequest)
		return
	}
	billingAlert, err = billingAlertApi.GetBillingAlert(ctx, projectId)
	return
}

func DeleteBillingAlert(ctx context.Context, projectId string) (billingAlert *model.BillingAlert, err error) {

	if projectId == "" {
		log.Printf("no project ID parameter\n")
		return nil, httperrors.New(errors.New("projectId not provided"), http.StatusBadRequest)
	}
	billingAlert, err = billingAlertApi.DeleteBillingAlert(ctx, projectId)
	return
}
