package model

import billing_state "gblaquiere.dev/multi-org-billing-alert/model/billing-state"

type BillingAlert struct {
	ProjectID     string              `json:"project_id,omitempty"`
	GroupAlert    *GroupAlert         `json:"group_alert,omitempty"`
	MonthlyBudget float32             `json:"monthly_budget"`
	Emails        []string            `json:"emails"`
	Thresholds    []float64           `json:"thresholds"`
	ChannelIds    []string            `json:"-"`
	State         billing_state.State `json:"state"`
}

type GroupAlert struct {
	ProjectIds []string `json:"project_ids,omitempty"`
	AlertName  string   `json:"name,omitempty"`
}
