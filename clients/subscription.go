package clients

import (
	"context"
	"time"

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

func (s *SubscriptionClient) ListCharges(ctx context.Context, userID string, from, to time.Time) ([]usecases.Charge, error) {
	resp, err := s.rpc.ListCharges(ctx, &subscriptionv1.ListChargesRequest{
		UserId: userID,
		From:   from.Format(time.DateOnly),
		To:     to.Format(time.DateOnly),
	})
	if err != nil {
		return nil, err
	}

	out := make([]usecases.Charge, 0, len(resp.GetCharges()))
	for _, c := range resp.GetCharges() {
		out = append(out, usecases.Charge{
			ID:                 c.GetId(),
			SubscriptionID:     c.GetSubscriptionId(),
			Name:               c.GetName(),
			Category:           c.GetCategory(),
			Type:               c.GetType(),
			Amount:             c.GetAmount(),
			ChargedAt:          c.GetChargedAt().AsTime(),
			SubscriptionStatus: c.GetSubscriptionStatus(),
		})
	}
	return out, nil
}
