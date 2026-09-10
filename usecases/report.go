package usecases

import "context"

type Subscription struct {
	ID     string
	UserID string
}

type SubscriptionSource interface {
	ListForReport(ctx context.Context, userID string) ([]Subscription, error)
}

type Report struct {
	UserID            string
	SubscriptionCount int
}

type ReportUsecase struct {
	subs SubscriptionSource
}

func NewReport(subs SubscriptionSource) *ReportUsecase {
	return &ReportUsecase{subs: subs}
}

func (u *ReportUsecase) GenerateReport(ctx context.Context, userID string) (Report, error) {
	list, err := u.subs.ListForReport(ctx, userID)
	if err != nil {
		return Report{}, err
	}
	return Report{
		UserID:            userID,
		SubscriptionCount: len(list),
	}, nil
}
