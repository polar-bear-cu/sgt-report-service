package usecases

import (
	"context"
	"math"
	"time"
)

const (
	defaultMonths = 12
	monthFormat   = "2006-01"
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

type MonthTotal struct {
	Month string
	Total float64
}

type MonthlySpending struct {
	Months          []MonthTotal
	Total           float64
	AveragePerMonth float64
}

// MonthlySpending returns the last n months up to and including the current month (still in progress).
func (u *ReportUsecase) MonthlySpending(ctx context.Context, userID string, n int) (MonthlySpending, error) {
	if n <= 0 {
		n = defaultMonths
	}
	current := u.monthStart(time.Now())
	from := current.AddDate(0, -(n - 1), 0)

	charges, err := u.charges.ListCharges(ctx, userID, from, current.AddDate(0, 1, 0))
	if err != nil {
		return MonthlySpending{}, err
	}

	byMonth := make(map[string]float64, n)
	for _, c := range charges {
		byMonth[c.ChargedAt.In(u.loc).Format(monthFormat)] += c.Amount
	}

	res := MonthlySpending{Months: make([]MonthTotal, 0, n)}
	for m := from; !m.After(current); m = m.AddDate(0, 1, 0) {
		key := m.Format(monthFormat)
		res.Months = append(res.Months, MonthTotal{Month: key, Total: round2(byMonth[key])})
		res.Total += byMonth[key]
	}
	res.Total = round2(res.Total)
	res.AveragePerMonth = round2(res.Total / float64(n))
	return res, nil
}

func (u *ReportUsecase) monthStart(t time.Time) time.Time {
	t = t.In(u.loc)
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, u.loc)
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
