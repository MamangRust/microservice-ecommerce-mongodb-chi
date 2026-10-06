package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-grpc-order/dto"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type OrderResult struct {
	OrderID    int32
	UserID     int32
	MerchantID int32
	TotalPrice int32
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
	DeletedAt  *time.Time
	TotalCount int64
}

type StockReservationResult struct {
	ReservationID int32
	OrderID       int32
	ProductID     int32
	Quantity      int32
	Status        string
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

type MerchantQueryRepository interface {
	FindByID(ctx context.Context, user_id int) (*dto.GetMerchantByIDRow, error)
}

type OrderCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateOrderRecordRequest) (*models.Order, error)
	Update(ctx context.Context, request *requests.UpdateOrderRecordRequest) (*models.Order, error)
	Trash(ctx context.Context, order_id int) (*models.Order, error)
	Restore(ctx context.Context, order_id int) (*models.Order, error)
	DeletePermanent(ctx context.Context, order_id int) (bool, error)
	DeletePermanentWithChildren(ctx context.Context, order_id int) (bool, error)
	FindTrashedByID(ctx context.Context, order_id int) (*models.Order, error)
	FindTrashed(ctx context.Context) ([]*models.Order, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

type OrderQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllOrder) ([]*OrderResult, error)
	FindActive(ctx context.Context, req *requests.FindAllOrder) ([]*OrderResult, error)
	FindTrashed(ctx context.Context, req *requests.FindAllOrder) ([]*OrderResult, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllOrderByMerchant) ([]*OrderResult, error)
	FindByID(ctx context.Context, order_id int) (*models.Order, error)
}

type OutboxRepository interface {
	Create(ctx context.Context, topic, key string, payload []byte) (*models.OutboxEvent, error)
	GetPending(ctx context.Context, limit int) ([]*models.OutboxEvent, error)
	Claim(ctx context.Context, limit int, leaseUntil time.Time) ([]*models.OutboxEvent, error)
	MarkDelivered(ctx context.Context, outboxID int64) (*models.OutboxEvent, error)
	MarkFailed(ctx context.Context, outboxID int64, nextAttemptAt time.Time) (*models.OutboxEvent, error)
	MarkDead(ctx context.Context, outboxID int64) (*models.OutboxEvent, error)
	DeleteOld(ctx context.Context, cutoff time.Time) (int64, error)
}

type StockReservationRepository interface {
	GetByOrder(ctx context.Context, orderID int) ([]*models.OrderStockReservation, error)
	Upsert(ctx context.Context, orderID, productID, quantity int) (*models.OrderStockReservation, error)
	UpdateQuantity(ctx context.Context, orderID, productID, quantity int) (*models.OrderStockReservation, error)
	Release(ctx context.Context, orderID, productID int) (*models.OrderStockReservation, error)
	Reserve(ctx context.Context, orderID, productID int) (*models.OrderStockReservation, error)
	GetReservedForTrashedOrders(ctx context.Context) ([]*models.OrderStockReservation, error)
	GetReleasedForTrashedOrders(ctx context.Context) ([]*models.OrderStockReservation, error)
	DeleteByOrder(ctx context.Context, orderID int) error
	DeleteByOrderProduct(ctx context.Context, orderID, productID int) error
	DeleteAllForTrashedOrders(ctx context.Context) error
	GetReleasedForActiveOrders(ctx context.Context) ([]*models.OrderStockReservation, error)
	DeleteOldReleasedReservations(ctx context.Context, cutoff time.Time) (int64, error)
	DeleteOldProductStockAdjustments(ctx context.Context, cutoff time.Time) (int64, error)
}
