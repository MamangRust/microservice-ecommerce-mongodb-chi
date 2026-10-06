package cache

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-order/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type OrderQueryCache interface {
	GetOrderAllCache(ctx context.Context, req *requests.FindAllOrder) ([]*repository.OrderResult, *int, bool)
	SetOrderAllCache(ctx context.Context, req *requests.FindAllOrder, data []*repository.OrderResult, total *int)

	GetOrderActiveCache(ctx context.Context, req *requests.FindAllOrder) ([]*repository.OrderResult, *int, bool)
	SetOrderActiveCache(ctx context.Context, req *requests.FindAllOrder, data []*repository.OrderResult, total *int)

	GetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrder) ([]*repository.OrderResult, *int, bool)
	SetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrder, data []*repository.OrderResult, total *int)

	GetCachedOrderCache(ctx context.Context, orderID int) (*models.Order, bool)
	SetCachedOrderCache(ctx context.Context, data *models.Order)

	GetOrderByMerchantCache(ctx context.Context, req *requests.FindAllOrderByMerchant) ([]*repository.OrderResult, *int, bool)
	SetOrderByMerchantCache(ctx context.Context, req *requests.FindAllOrderByMerchant, data []*repository.OrderResult, total *int)
}

type OrderCommandCache interface {
	DeleteOrderCache(ctx context.Context, orderID int)
	InvalidateOrderCache(ctx context.Context)
}
