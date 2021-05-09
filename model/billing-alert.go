package model

type BillingAlert struct {
	ProjectID     string   `json:"project_id"`
	MonthlyBudget float32  `json:"monthly_budget"`
	Emails        []string `json:"emails"`
	ChannelIds    []string
}

/*
{
	"project_id": "gdglyon-cloudrun",
	"monthly_budget": 10,
	"emails":["guillaume.blaquiere@gmail.com"]
	}
'{"project_id": "gdglyon-cloudrun","monthly_budget": 10,"emails":["guillaume.blaquiere@gmail.com"]}'

ewoJInByb2plY3RfaWQiOiAiZ2RnbHlvbi1jbG91ZHJ1biIsCgkibW9udGhseV9idWRnZXQiOiAxMCwKCSJlbWFpbHMiOlsiZ3VpbGxhdW1lLmJsYXF1aWVyZUBnbWFpbC5jb20iXQoJfQ==
*/
