package dtos

import "time"

type MonthlySpendingQuery struct {
	Months *int `form:"months" binding:"omitempty,min=1,max=24"`
}

type MonthTotalResponse struct {
	Month string  `json:"month"`
	Total float64 `json:"total"`
}

type MonthlySpendingResponse struct {
	Months          []MonthTotalResponse `json:"months"`
	Total           float64              `json:"total"`
	AveragePerMonth float64              `json:"averagePerMonth"`
}

type SpendingBySubscriptionQuery struct {
	From string `form:"from"`
	To   string `form:"to"`
}

type SubscriptionSpendingResponse struct {
	SubscriptionID string  `json:"subscriptionId"`
	Name           string  `json:"name"`
	Category       string  `json:"category"`
	Type           string  `json:"type"`
	Status         string  `json:"status"`
	Count          int     `json:"count"`
	Total          float64 `json:"total"`
}

type SpendingBySubscriptionResponse struct {
	From  string                         `json:"from"`
	To    string                         `json:"to"`
	Total float64                        `json:"total"`
	Items []SubscriptionSpendingResponse `json:"items"`
}

type MonthChargesQuery struct {
	Month string `form:"month"`
}

type CategoryTotalResponse struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
}

type ChargeResponse struct {
	SubscriptionID string    `json:"subscriptionId"`
	Name           string    `json:"name"`
	Category       string    `json:"category"`
	Type           string    `json:"type"`
	Amount         float64   `json:"amount"`
	ChargedAt      time.Time `json:"chargedAt"`
}

type MonthChargesResponse struct {
	Month      string                  `json:"month"`
	Total      float64                 `json:"total"`
	Count      int                     `json:"count"`
	Categories []CategoryTotalResponse `json:"categories"`
	Items      []ChargeResponse        `json:"items"`
}

type ExportChargesQuery struct {
	From string `form:"from"`
	To   string `form:"to"`
}
