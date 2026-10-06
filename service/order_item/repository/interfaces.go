package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type OrderItemResult struct {
	OrderItemID int32
	OrderID     int32
	ProductID   int32
	Quantity    int32
	Price       int32
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
	TotalCount  int64
}

type OrderItemQueryRepository interface {
	FindAll(
		ctx context.Context,
		req *requests.FindAllOrderItems,
	) ([]*OrderItemResult, error)

	FindActive(
		ctx context.Context,
		req *requests.FindAllOrderItems,
	) ([]*OrderItemResult, error)

	FindTrashed(
		ctx context.Context,
		req *requests.FindAllOrderItems,
	) ([]*OrderItemResult, error)

	FindOrderItemByOrder(
		ctx context.Context,
		order_id int,
	) ([]*OrderItemResult, error)
}

type OrderItemCommandRepository interface {
	Create(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*models.OrderItem, error)
	Update(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*models.OrderItem, error)

	Trash(ctx context.Context, orderItemID int) (*models.OrderItem, error)
	Restore(ctx context.Context, orderItemID int) (*models.OrderItem, error)
	DeletePermanent(ctx context.Context, orderItemID int) (bool, error)
	DeleteOrderItemByOrderPermanent(ctx context.Context, orderID int) (bool, error)

	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)

	CalculateTotalPrice(ctx context.Context, orderID int) (int, error)
}
