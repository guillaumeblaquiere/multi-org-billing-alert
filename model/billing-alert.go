package model

type BillingAlert struct {
	ProjectID     string    `json:"project_id"`
	MonthlyBudget float32   `json:"monthly_budget"`
	Emails        []string  `json:"emails"`
	Thresholds    []float64 `json:"thresholds""`
	ChannelIds    []string  `json:"-""`
}
