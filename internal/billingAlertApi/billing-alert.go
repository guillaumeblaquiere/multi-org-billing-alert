package billingAlertApi

import (
	budgetApi "cloud.google.com/go/billing/budgets/apiv1"
	"context"
	"errors"
	"fmt"
	"gblaquiere.dev/multi-org-billing-alert/internal/httperrors"
	"gblaquiere.dev/multi-org-billing-alert/internal/notificationChannelApi"
	"gblaquiere.dev/multi-org-billing-alert/model"
	ressourceManager "google.golang.org/api/cloudresourcemanager/v3"
	"google.golang.org/api/iterator"
	budgetModel "google.golang.org/genproto/googleapis/cloud/billing/budgets/v1"
	"google.golang.org/genproto/googleapis/type/money"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"log"
	"net/http"
	"os"
)

const projectPrefix = "projects/"

var defaultThresholds = []*budgetModel.ThresholdRule{ //predefined threshold
	{
		ThresholdPercent: 0.5,
	},
	{
		ThresholdPercent: 0.9,
	},
	{
		ThresholdPercent: 1.0,
	}}

var client *budgetApi.BudgetClient = nil
var clientResourceManager *ressourceManager.Service = nil

//Initialize the client at startup
func init() {
	var err error
	ctx := context.Background()
	client, err = budgetApi.NewBudgetClient(ctx)
	if err != nil {
		log.Panicf("budgetApi.NewBudgetClient: %+v\n", err)
	}
	clientResourceManager, err = ressourceManager.NewService(ctx)
	if err != nil {
		log.Panicf("budgetApi.NewResourceManagerClient: %+v\n", err)
	}
}

func UpsertBillingAlert(ctx context.Context, message *model.BillingAlert) (err error) {

	//Check if the budgetApi exists
	b, err := getExistingBillingAlert(ctx, getMessageName(message))

	if err != nil {
		return err
	}

	// Create or update accordingly
	if b == nil {
		//Create a new budget
		err = createNewBudget(ctx, message)
		if err != nil {
			return err
		}
		log.Printf("Budget creation successful for the project %s \n", message.ProjectID)
	} else {
		//Update the retrieved budget
		err = updateBudget(ctx, message, b)
		if err != nil {
			return err
		}
		log.Printf("Budget update successful for the project %s \n", message.ProjectID)
	}

	return
}

func getExistingBillingAlert(ctx context.Context, alertName string) (b *budgetModel.Budget, err error) {
	req := &budgetModel.ListBudgetsRequest{
		Parent: getBillingParent(),
	}
	budgets := client.ListBudgets(ctx, req)

	displayName := getDisplayName(alertName)
	for {
		budget, err := budgets.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("budgets.Next: %+v\n", err)
			return nil, err
		}
		if budget.DisplayName == displayName {
			if b == nil {
				b = budget
			} else {
				err = errors.New("2 or more budget exists") //internal error, should not exist.
				log.Printf("impossible to get the budget, %+v, for this budget Name %s \n", err, displayName)
				return nil, err
			}
		}
	}
	return
}

func updateBudget(ctx context.Context, message *model.BillingAlert, b *budgetModel.Budget) (err error) {
	updatedPath := updateBudgetAlert(b, message)
	req := &budgetModel.UpdateBudgetRequest{
		Budget: b,
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: updatedPath,
		},
	}
	_, err = client.UpdateBudget(ctx, req)
	if err != nil {
		log.Printf("client.UpdateBudget: %+v\n", err)
		return err
	}
	return
}

func createNewBudget(ctx context.Context, message *model.BillingAlert) (err error) {
	b := &budgetModel.Budget{ //Initiate a new budget alert object
		DisplayName:    getDisplayName(getMessageName(message)),
		ThresholdRules: defaultThresholds,
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

func getMessageName(message *model.BillingAlert) string {
	name := message.ProjectID
	if message.GroupAlert != nil {
		name = message.GroupAlert.AlertName
	}
	return name
}

//Merge all the project IDs from message and from the projectList in group alert
func getProjectList(message *model.BillingAlert) (projectList []string) {
	if message.ProjectID != "" {
		projectList = append(projectList, fmt.Sprintf("%s%s", projectPrefix, message.ProjectID))
	}
	if message.GroupAlert != nil && len(message.GroupAlert.ProjectIds) > 0 {
		for _, projectId := range message.GroupAlert.ProjectIds {
			projectList = append(projectList, fmt.Sprintf("%s%s", projectPrefix, projectId))
		}
	}
	return
}

func updateBudgetAlert(b *budgetModel.Budget, message *model.BillingAlert) (updatedPath []string) {
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
	b.BudgetFilter = &budgetModel.Filter{
		Projects: getProjectList(message),
	}

	updatedPath = []string{ //Only these 2 fields to update. Can add more if required
		"amount.specified_amount",
		"notifications_rule",
		"budget_filter.projects",
	}

	//Optional thresholds
	if len(message.Thresholds) > 0 {
		updatedPath = append(updatedPath, "threshold_rules")

		newThreshold := []*budgetModel.ThresholdRule{}
		for _, threshold := range message.Thresholds {
			newThreshold = append(newThreshold, &budgetModel.ThresholdRule{
				ThresholdPercent: threshold,
			})
		}
		b.ThresholdRules = newThreshold
	}

	return
}

func getDisplayName(alertName string) string {
	return fmt.Sprintf("billing-%s", alertName) //static naming
}

func getBillingParent() string {
	return fmt.Sprintf("billingAccounts/%s", os.Getenv("BILLING_ACCOUNT"))
}

func GetBillingAlert(ctx context.Context, alertName string) (billingAlert *model.BillingAlert, err error) {

	b, err := getExistingBillingAlert(ctx, alertName)
	if err != nil {
		return
	}

	if b == nil {
		err = httperrors.New(errors.New(fmt.Sprintf("projectid %q does not exist", alertName)), http.StatusNotFound)
		return
	}

	//Prepare response
	emails, err := getEmailList(ctx, b)
	if err != nil {
		return
	}

	billingAlert = &model.BillingAlert{
		MonthlyBudget: float32(b.Amount.GetSpecifiedAmount().Units) + (float32(b.Amount.GetSpecifiedAmount().Nanos) / 1000000000),
		Emails:        emails,
		Thresholds:    getThresholds(b),
	}

	if len(b.BudgetFilter.GetProjects()) > 1 {
		projectList, err := getProjectIds(ctx, b.BudgetFilter.Projects)
		if err != nil {
			return billingAlert, err
		}
		g := &model.GroupAlert{
			ProjectIds: projectList,
			AlertName:  alertName,
		}
		billingAlert.GroupAlert = g
	} else {
		billingAlert.ProjectID = alertName
	}

	return
}

func getProjectIds(ctx context.Context, projects []string) (projectList []string, err error) {
	for _, project := range projects {
		p, err := clientResourceManager.Projects.Get(project).Context(ctx).Do()
		if err != nil {
			log.Printf("resourceManager.GetProject: %+v\n", err)
			return projectList, err
		}
		projectList = append(projectList, p.ProjectId)
	}
	return
}

func getThresholds(b *budgetModel.Budget) (thresholds []float64) {
	for _, t := range b.ThresholdRules {
		thresholds = append(thresholds, t.ThresholdPercent)
	}
	return
}

func getEmailList(ctx context.Context, b *budgetModel.Budget) (emailList []string, err error) {
	for _, channelId := range b.NotificationsRule.MonitoringNotificationChannels {
		channel, err := notificationChannelApi.GetChannelID(ctx, channelId)
		if err != nil {
			return nil, httperrors.New(err, http.StatusInternalServerError)
		}
		email := channel.Labels[notificationChannelApi.EmailAddressLabelKey]
		emailList = append(emailList, email)
	}
	return
}

func DeleteBillingAlert(ctx context.Context, alertName string) (err error) {
	b, err := getExistingBillingAlert(ctx, alertName)
	if err != nil {
		return httperrors.New(err, http.StatusInternalServerError)
	}

	if b == nil {
		return httperrors.New(errors.New(fmt.Sprintf("projectid %q does not exist", alertName)), http.StatusNotFound)
	}

	req := &budgetModel.DeleteBudgetRequest{
		Name: b.Name,
	}
	err = client.DeleteBudget(ctx, req)
	return
}
