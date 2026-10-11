package dtos

type MonthlySpendingQuery struct {
	Months int `form:"months" binding:"omitempty,min=1,max=24"`
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
