package usecases

import (
	"context"
	"time"
)

type Charge struct {
	ID                 string
	SubscriptionID     string
	Name               string
	Category           string
	Type               string
	Amount             float64
	ChargedAt          time.Time
	SubscriptionStatus string
}

type ChargeSource interface {
	ListCharges(ctx context.Context, userID string, from, to time.Time) ([]Charge, error)
}

type ReportUsecase struct {
	charges ChargeSource
	loc     *time.Location
}

func NewReport(charges ChargeSource, loc *time.Location) *ReportUsecase {
	return &ReportUsecase{charges: charges, loc: loc}
}
