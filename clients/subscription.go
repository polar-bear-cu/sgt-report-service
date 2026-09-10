package clients

import (
	"context"

	subscriptionv1 "github.com/polar-bear-cu/sgt-proto/gen/go/subscription/v1"
	"google.golang.org/grpc"

	"github.com/polar-bear-cu/sgt-report-service/usecases"
)

type SubscriptionClient struct {
	rpc subscriptionv1.SubscriptionServiceClient
}

func NewSubscriptionClient(conn *grpc.ClientConn) *SubscriptionClient {
	return &SubscriptionClient{rpc: subscriptionv1.NewSubscriptionServiceClient(conn)}
}

func (s *SubscriptionClient) ListForReport(ctx context.Context, userID string) ([]usecases.Subscription, error) {
	resp, err := s.rpc.GetSubscriptionsForReport(ctx, &subscriptionv1.GetSubscriptionsForReportRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}

	out := make([]usecases.Subscription, 0, len(resp.GetSubscription()))
	for _, sub := range resp.GetSubscription() {
		out = append(out, usecases.Subscription{ID: sub.GetId(), UserID: sub.GetUserId()})
	}
	return out, nil
}
