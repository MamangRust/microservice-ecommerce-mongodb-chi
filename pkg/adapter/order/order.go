package order

import (
	"context"

	pborder "github.com/MamangRust/microservice-ecommerce-grpc-pb/order"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-pkg/resilience"
)

// QueryRepository is the contract consumers depend on for order reads.
type QueryRepository interface {
	FindByID(ctx context.Context, orderID int) (*models.Order, error)
}

// Repository implements QueryRepository on top of the generated order query
// client, guarding every call.
type Repository struct {
	client pborder.OrderQueryServiceClient
	guard  *resilience.DependencyGuard
}

// New builds the order adapter from its query client.
func New(client pborder.OrderQueryServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{client: client}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// NewQueryAdapter wraps the generated order query client.
func NewQueryAdapter(client pborder.OrderQueryServiceClient, opts ...adapter.GuardOption) QueryRepository {
	return New(client, opts...)
}

// SetGuard attaches the dependency guard used to wrap every outbound call.
func (r *Repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

// FindByID propagates the dependency's gRPC status unchanged (NotFound -> 404,
// and so on) instead of wrapping it in a domain error.
func (r *Repository) FindByID(ctx context.Context, orderID int) (*models.Order, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pborder.ApiResponseOrder, error) {
		return r.client.FindById(callCtx, &pborder.FindByIdOrderRequest{Id: int32(orderID)})
	})
	if err != nil {
		return nil, err
	}

	return &models.Order{
		OrderID:    res.Data.Id,
		UserID:     res.Data.UserId,
		MerchantID: res.Data.MerchantId,
		TotalPrice: res.Data.TotalPrice,
	}, nil
}

var _ QueryRepository = (*Repository)(nil)
