package internal

import (
	"context"
	"gblaquiere.dev/multi-org-billing-alert/model"
)

func CreateBillingAlert(ctx context.Context, message *model.BillingAlert) (err error) {

	//Get the channelIDs list for all the email
	err = getChannelIDs(ctx, message)
	if err != nil {
		return err
	}

	//Create or update the alert budgetApi
	err = billingAlert(ctx, message)
	return
}
