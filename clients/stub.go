package clients

import (
	"context"

	"github.com/polar-bear-cu/sgt-report-service/usecases"
)

type StubSubscriptionSource struct{}

func (StubSubscriptionSource) ListForReport(_ context.Context, userID string) ([]usecases.Subscription, error) {
	return []usecases.Subscription{
		{ID: "sub-1", UserID: userID},
		{ID: "sub-2", UserID: userID},
		{ID: "sub-3", UserID: userID},
	}, nil
}
