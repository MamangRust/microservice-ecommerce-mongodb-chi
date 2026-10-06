package cache

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-order-item/repository"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type OrderItemQueryCache interface {
	GetCachedOrderItemsAll(ctx context.Context, req *requests.FindAllOrderItems) ([]*repository.OrderItemResult, *int, bool)
	SetCachedOrderItemsAll(ctx context.Context, req *requests.FindAllOrderItems, data []*repository.OrderItemResult, total *int)

	GetCachedOrderItemActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*repository.OrderItemResult, *int, bool)
	SetCachedOrderItemActive(ctx context.Context, req *requests.FindAllOrderItems, data []*repository.OrderItemResult, total *int)

	GetCachedOrderItemTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*repository.OrderItemResult, *int, bool)
	SetCachedOrderItemTrashed(ctx context.Context, req *requests.FindAllOrderItems, data []*repository.OrderItemResult, total *int)

	GetCachedOrderItems(ctx context.Context, orderID int) ([]*repository.OrderItemResult, bool)
	SetCachedOrderItems(ctx context.Context, data []*repository.OrderItemResult)
}

type OrderItemCommandCache interface {
	InvalidateOrderItemCache(ctx context.Context) error
}
