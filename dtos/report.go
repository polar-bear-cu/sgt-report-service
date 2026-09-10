package dtos

type ReportResponse struct {
	UserID            string `json:"userId"`
	SubscriptionCount int    `json:"subscriptionCount"`
}
