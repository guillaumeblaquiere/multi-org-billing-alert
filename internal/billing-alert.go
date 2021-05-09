package internal

import (
	budgetApi "cloud.google.com/go/billing/budgets/apiv1"
	"context"
	"errors"
	"fmt"
	"gblaquiere.dev/multi-org-billing-alert/model"
	"google.golang.org/api/iterator"
	budgetModel "google.golang.org/genproto/googleapis/cloud/billing/budgets/v1"
	"google.golang.org/genproto/googleapis/type/money"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"log"
	"os"
)

func billingAlert(ctx context.Context, message *model.BillingAlert) (err error) {
	client, err := budgetApi.NewBudgetClient(ctx)
	if err != nil {
		log.Printf("budgetApi.NewBudgetClient: %+v\n", err)
		return err
	}

	//Check if the budgetApi exists
	req := &budgetModel.ListBudgetsRequest{
		Parent: getBillingParent(),
	}
	budgets := client.ListBudgets(ctx, req)

	var b *budgetModel.Budget
	displayName := getDisplayName(message)
	for {
		budget, err := budgets.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("budgets.Next: %+v\n", err)
			return err
		}
		if budget.DisplayName == displayName {
			if b == nil {
				b = budget
			} else {
				err = errors.New("2 or more budget exists")
				log.Printf("impossible to get the budget, %+v, for this budget Name %s \n", err, displayName)
				return err
			}
		}
	}

	// Create or update accordingly
	if b == nil {
		//Create a new budget
		err = createNewBudget(ctx, client, message)
		if err != nil {
			return err
		}
		log.Printf("Budget creation successful for the project %s \n", message.ProjectID)
	} else {
		//Update the retrieved budget
		err = updateBudget(ctx, client, message, b)
		if err != nil {
			return err
		}
		log.Printf("Budget update successful for the project %s \n", message.ProjectID)
	}

	return
}

func updateBudget(ctx context.Context, client *budgetApi.BudgetClient, message *model.BillingAlert, b *budgetModel.Budget) (err error) {
	updateBudgetAlert(b, message)
	req := &budgetModel.UpdateBudgetRequest{
		Budget: b,
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{ //Only these 2 fields to update. Can add more if required
				"amount.specified_amount",
				"notifications_rule",
			},
		},
	}
	_, err = client.UpdateBudget(ctx, req)
	if err != nil {
		log.Printf("client.UpdateBudget: %+v\n", err)
		return err
	}
	return
}

func createNewBudget(ctx context.Context, client *budgetApi.BudgetClient, message *model.BillingAlert) (err error) {
	b := &budgetModel.Budget{ //Initiate a new budget alert object
		DisplayName: getDisplayName(message),
		BudgetFilter: &budgetModel.Filter{
			Projects: []string{fmt.Sprintf("projects/%s", message.ProjectID)},
		},
		ThresholdRules: []*budgetModel.ThresholdRule{ //predefined threshold
			{
				ThresholdPercent: 0.5,
			},
			{
				ThresholdPercent: 0.9,
			},
			{
				ThresholdPercent: 1.0,
			},
		},
	}
	updateBudgetAlert(b, message)
	req := &budgetModel.CreateBudgetRequest{
		Parent: getBillingParent(),
		Budget: b,
	}
	_, err = client.CreateBudget(ctx, req)
	if err != nil {
		log.Printf("client.CreateBudget: %+v\n", err)
		return err
	}
	return
}

func updateBudgetAlert(b *budgetModel.Budget, message *model.BillingAlert) {
	b.Amount = &budgetModel.BudgetAmount{
		BudgetAmount: &budgetModel.BudgetAmount_SpecifiedAmount{
			SpecifiedAmount: &money.Money{
				CurrencyCode: "EUR",                                                                               //static currency
				Units:        int64(message.MonthlyBudget),                                                        //get only the int part
				Nanos:        int32((message.MonthlyBudget - float32(int32(message.MonthlyBudget))) * 1000000000), //remove the int part and set the floating part 10^9 to get a int
			},
		},
	}
	b.NotificationsRule = &budgetModel.NotificationsRule{
		MonitoringNotificationChannels: message.ChannelIds,
		DisableDefaultIamRecipients:    true, //to not disturb the Billing administrator
	}
}

func getDisplayName(message *model.BillingAlert) string {
	return fmt.Sprintf("billing-%s", message.ProjectID) //static naming
}

func getBillingParent() string {
	return fmt.Sprintf("billingAccounts/%s", os.Getenv("BILLING_ACCOUNT"))
}
