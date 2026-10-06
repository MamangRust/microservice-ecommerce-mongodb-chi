package transaction

import (
	"context"
	"fmt"

	pbtransaction "github.com/MamangRust/microservice-ecommerce-grpc-pb/transaction"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	"github.com/MamangRust/microservice-ecommerce-pkg/resilience"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CommandRepository is the contract consumers depend on for transaction writes.
type CommandRepository interface {
	DeleteByOrderIDPermanent(ctx context.Context, orderID int) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

// Repository implements CommandRepository on top of the generated transaction
// command client, guarding every call.
type Repository struct {
	client pbtransaction.TransactionCommandServiceClient
	guard  *resilience.DependencyGuard
}

// New builds the transaction adapter from its command client.
func New(client pbtransaction.TransactionCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{client: client}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// NewCommandAdapter wraps the generated transaction command client.
func NewCommandAdapter(client pbtransaction.TransactionCommandServiceClient, opts ...adapter.GuardOption) CommandRepository {
	return New(client, opts...)
}

// SetGuard attaches the dependency guard used to wrap every outbound call.
func (r *Repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func (r *Repository) DeleteByOrderIDPermanent(ctx context.Context, orderID int) (bool, error) {
	if r.client == nil {
		return false, fmt.Errorf("transaction command client is not initialized")
	}
	if _, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbtransaction.ApiResponseTransactionDelete, error) {
		return r.client.DeleteTransactionByOrderPermanent(callCtx, &pbtransaction.FindByIdTransactionRequest{
			Id: int32(orderID),
		})
	}); err != nil {
		return false, err
	}

	return true, nil
}

func (r *Repository) DeleteAll(ctx context.Context) (bool, error) {
	if r.client == nil {
		return false, fmt.Errorf("transaction command client is not initialized")
	}
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbtransaction.ApiResponseTransactionAll, error) {
		return r.client.DeleteAllTransactionPermanent(callCtx, &emptypb.Empty{})
	})
	if err != nil {
		return false, err
	}

	return res.Status == "success", nil
}

var _ CommandRepository = (*Repository)(nil)
