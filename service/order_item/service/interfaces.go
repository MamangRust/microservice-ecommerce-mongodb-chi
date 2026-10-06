package service

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-order-item/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type OrderItemQueryService interface {
	FindAll(
		ctx context.Context,
		req *requests.FindAllOrderItems,
	) ([]*repository.OrderItemResult, *int, error)

	FindActive(
		ctx context.Context,
		req *requests.FindAllOrderItems,
	) ([]*repository.OrderItemResult, *int, error)

	FindTrashed(
		ctx context.Context,
		req *requests.FindAllOrderItems,
	) ([]*repository.OrderItemResult, *int, error)

	FindByOrder(
		ctx context.Context,
		order_id int,
	) ([]*repository.OrderItemResult, error)
}

type OrderItemCommandService interface {
	Create(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*models.OrderItem, error)
	Update(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*models.OrderItem, error)

	Trash(ctx context.Context, orderItemID int) (*models.OrderItem, error)
	Restore(ctx context.Context, orderItemID int) (*models.OrderItem, error)
	DeletePermanent(ctx context.Context, orderItemID int) (bool, error)
	DeleteByOrderPermanent(ctx context.Context, orderID int) (bool, error)

	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)

	CalculateTotalPrice(ctx context.Context, orderID int) (int, error)
}
