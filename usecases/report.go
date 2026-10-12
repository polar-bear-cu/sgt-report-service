package usecases

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"math"
	"slices"
	"time"
)

var ErrValidation = errors.New("validation failed")

const (
	defaultMonths = 12
	maxMonths     = 24
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

type SubscriptionSpending struct {
	SubscriptionID string
	Name           string
	Category       string
	Type           string
	Status         string
	Count          int
	Total          float64
}

type SpendingBySubscription struct {
	From  string
	To    string
	Total float64
	Items []SubscriptionSpending
}

// SpendingBySubscription sums the charges of each subscription from the start of month from
// to the end of month to (default: the last 12 months), highest total first.
func (u *ReportUsecase) SpendingBySubscription(ctx context.Context, userID, from, to string) (SpendingBySubscription, error) {
	start, end, err := u.monthRange(from, to)
	if err != nil {
		return SpendingBySubscription{}, err
	}

	charges, err := u.charges.ListCharges(ctx, userID, start, end.AddDate(0, 1, 0))
	if err != nil {
		return SpendingBySubscription{}, err
	}

	res := SpendingBySubscription{From: start.Format(monthFormat), To: end.Format(monthFormat)}
	bySub := make(map[string]*SubscriptionSpending)
	for _, c := range charges {
		s, ok := bySub[c.SubscriptionID]
		if !ok {
			s = &SubscriptionSpending{SubscriptionID: c.SubscriptionID, Name: c.Name, Category: c.Category, Status: c.SubscriptionStatus}
			bySub[c.SubscriptionID] = s
		}
		s.Type = c.Type // charges are oldest first, so this ends as the latest type
		s.Count++
		s.Total += c.Amount
		res.Total += c.Amount
	}

	res.Items = make([]SubscriptionSpending, 0, len(bySub))
	for _, s := range bySub {
		s.Total = round2(s.Total)
		res.Items = append(res.Items, *s)
	}
	slices.SortFunc(res.Items, func(a, b SubscriptionSpending) int {
		return cmp.Or(cmp.Compare(b.Total, a.Total), cmp.Compare(a.Name, b.Name))
	})
	res.Total = round2(res.Total)
	return res, nil
}

type CategoryTotal struct {
	Category string
	Total    float64
}

type MonthCharges struct {
	Month      string
	Total      float64
	Count      int
	Categories []CategoryTotal
	Items      []Charge
}

// MonthCharges lists the charges of one month (default: the current month), oldest first,
// with the total of each category, highest first.
func (u *ReportUsecase) MonthCharges(ctx context.Context, userID, month string) (MonthCharges, error) {
	start := u.monthStart(time.Now())
	if month != "" {
		t, err := u.parseMonth(month, "month")
		if err != nil {
			return MonthCharges{}, err
		}
		start = t
	}

	charges, err := u.charges.ListCharges(ctx, userID, start, start.AddDate(0, 1, 0))
	if err != nil {
		return MonthCharges{}, err
	}

	res := MonthCharges{Month: start.Format(monthFormat), Count: len(charges), Items: make([]Charge, 0, len(charges))}
	byCategory := make(map[string]float64)
	for _, c := range charges {
		c.ChargedAt = c.ChargedAt.In(u.loc)
		res.Items = append(res.Items, c)
		byCategory[c.Category] += c.Amount
		res.Total += c.Amount
	}

	res.Categories = make([]CategoryTotal, 0, len(byCategory))
	for cat, total := range byCategory {
		res.Categories = append(res.Categories, CategoryTotal{Category: cat, Total: round2(total)})
	}
	slices.SortFunc(res.Categories, func(a, b CategoryTotal) int {
		return cmp.Or(cmp.Compare(b.Total, a.Total), cmp.Compare(a.Category, b.Category))
	})
	res.Total = round2(res.Total)
	return res, nil
}

type ChargesExport struct {
	From  string
	To    string
	Items []Charge
}

// ExportCharges lists every charge in a month range (default: the last 12 months), oldest first.
func (u *ReportUsecase) ExportCharges(ctx context.Context, userID, from, to string) (ChargesExport, error) {
	start, end, err := u.monthRange(from, to)
	if err != nil {
		return ChargesExport{}, err
	}

	charges, err := u.charges.ListCharges(ctx, userID, start, end.AddDate(0, 1, 0))
	if err != nil {
		return ChargesExport{}, err
	}
	for i := range charges {
		charges[i].ChargedAt = charges[i].ChargedAt.In(u.loc)
	}
	return ChargesExport{From: start.Format(monthFormat), To: end.Format(monthFormat), Items: charges}, nil
}

// monthRange returns the first day of month from and of month to, both included.
func (u *ReportUsecase) monthRange(from, to string) (time.Time, time.Time, error) {
	end := u.monthStart(time.Now())
	if to != "" {
		t, err := u.parseMonth(to, "to")
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		end = t
	}

	start := end.AddDate(0, -(defaultMonths - 1), 0)
	if from != "" {
		t, err := u.parseMonth(from, "from")
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		start = t
	}

	if start.After(end) {
		return time.Time{}, time.Time{}, fmt.Errorf("%w: from must not be after to", ErrValidation)
	}
	if months := (end.Year()-start.Year())*12 + int(end.Month()-start.Month()) + 1; months > maxMonths {
		return time.Time{}, time.Time{}, fmt.Errorf("%w: range must be at most %d months", ErrValidation, maxMonths)
	}
	return start, end, nil
}

func (u *ReportUsecase) parseMonth(s, field string) (time.Time, error) {
	t, err := time.ParseInLocation(monthFormat, s, u.loc)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %s must be YYYY-MM", ErrValidation, field)
	}
	return t, nil
}

func (u *ReportUsecase) monthStart(t time.Time) time.Time {
	t = t.In(u.loc)
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, u.loc)
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
