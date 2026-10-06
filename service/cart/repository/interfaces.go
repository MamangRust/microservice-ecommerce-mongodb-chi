package repository

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

// CartResult is the result type for paginated cart queries.
type CartResult struct {
	CartID     int32
	UserID     int32
	ProductID  int32
	Name       string
	Price      int32
	Image      string
	Quantity   int32
	Weight     int32
	CreatedAt  string
	UpdatedAt  string
	TotalCount int64
}

// CartCreateResult is the result type for cart creation.
type CartCreateResult struct {
	CartID    int32
	UserID    int32
	ProductID int32
	Name      string
	Price     int32
	Image     string
	Quantity  int32
	Weight    int32
	CreatedAt string
	UpdatedAt string
}

type CartQueryRepository interface {
	FindCarts(ctx context.Context, req *requests.FindAllCarts) ([]*CartResult, error)
}

type CartCommandRepository interface {
	CreateCart(ctx context.Context, req *requests.CartCreateRecord) (*CartCreateResult, error)
	DeletePermanent(ctx context.Context, req *requests.DeleteCartRequest) (bool, error)
	DeleteAllPermanently(ctx context.Context, req *requests.DeleteAllCartRequest) (bool, error)
}
