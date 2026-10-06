package cache

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-cart/repository"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type CartQueryCache interface {
	GetCachedCartsCache(ctx context.Context, request *requests.FindAllCarts) ([]*repository.CartResult, *int, bool)
	SetCartsCache(ctx context.Context, request *requests.FindAllCarts, response []*repository.CartResult, total *int)

	// DeleteCartsCache invalidates every cached cart listing for a user so
	// mutations (create/delete/delete-all) never serve stale data.
	DeleteCartsCache(ctx context.Context, userID int)
}

type CartMencache interface {
	CartQueryCache
}
