package order_item

import (
	"context"

	pborder_item "github.com/MamangRust/microservice-ecommerce-grpc-pb/order_item"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-pkg/resilience"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	order_item_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/order_item_errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

// QueryRepository is the contract consumers depend on for order-item reads.
type QueryRepository interface {
	FindOrderItemByOrder(ctx context.Context, orderID int) ([]*models.OrderItem, error)
	// CalculateTotalPrice is served by the command service (it is the RPC the
	// owning service exposes), but callers treat it as a read.
	CalculateTotalPrice(ctx context.Context, orderID int) (*int32, error)
}

// CommandRepository is the contract consumers depend on for order-item writes.
type CommandRepository interface {
	Create(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*models.OrderItem, error)
	Update(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*models.OrderItem, error)
	Trash(ctx context.Context, orderID int) (*models.OrderItem, error)
	Restore(ctx context.Context, orderID int) (*models.OrderItem, error)
	DeletePermanent(ctx context.Context, orderID int) (bool, error)
	DeleteByOrderIDPermanent(ctx context.Context, orderID int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

// Repository implements both QueryRepository and CommandRepository on top of the
// generated order-item query/command clients, guarding every call.
type Repository struct {
	query   pborder_item.OrderItemQueryServiceClient
	command pborder_item.OrderItemCommandServiceClient
	guard   *resilience.DependencyGuard
}

// New builds the order-item adapter from its query and command clients.
func New(query pborder_item.OrderItemQueryServiceClient, command pborder_item.OrderItemCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{query: query, command: command}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// NewAdapter wraps the generated order-item clients.
func NewAdapter(query pborder_item.OrderItemQueryServiceClient, command pborder_item.OrderItemCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	return New(query, command, opts...)
}

// NewQueryAdapter returns the adapter restricted to the QueryRepository surface.
func NewQueryAdapter(query pborder_item.OrderItemQueryServiceClient, opts ...adapter.GuardOption) *Repository {
	return New(query, nil, opts...)
}

// NewCommandAdapter returns the adapter restricted to the CommandRepository surface.
func NewCommandAdapter(command pborder_item.OrderItemCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	return New(nil, command, opts...)
}

// SetGuard attaches the dependency guard used to wrap every outbound call.
func (r *Repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func (r *Repository) FindOrderItemByOrder(ctx context.Context, orderID int) ([]*models.OrderItem, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pborder_item.ApiResponsesOrderItem, error) {
		return r.query.FindOrderItemByOrder(callCtx, &pborder_item.FindByIdOrderItemRequest{Id: int32(orderID)})
	})
	if err != nil {
		return nil, order_item_errors.ErrFindOrderItemByOrder.WithInternal(err)
	}

	items := make([]*models.OrderItem, 0, len(res.Data))
	for _, item := range res.Data {
		items = append(items, toModel(item))
	}

	return items, nil
}

func (r *Repository) CalculateTotalPrice(ctx context.Context, orderID int) (*int32, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pborder_item.CalculateTotalPriceResponse, error) {
		return r.command.CalculateTotalPrice(callCtx, &pborder_item.CalculateTotalPriceRequest{OrderId: int32(orderID)})
	})
	if err != nil {
		return nil, order_item_errors.ErrCalculateTotalPrice.WithInternal(err)
	}

	total := int32(res.TotalPrice)
	return &total, nil
}

func (r *Repository) Create(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*models.OrderItem, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pborder_item.ApiResponseOrderItem, error) {
		return r.command.CreateOrderItem(callCtx, &pborder_item.CreateOrderItemRecordRequest{
			OrderId:   int32(req.OrderID),
			ProductId: int32(req.ProductID),
			Quantity:  int32(req.Quantity),
			Price:     int32(req.Price),
		})
	})
	if err != nil {
		return nil, order_item_errors.ErrCreateOrderItem.WithInternal(err)
	}
	return toModel(res.Data), nil
}

func (r *Repository) Update(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*models.OrderItem, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pborder_item.ApiResponseOrderItem, error) {
		return r.command.UpdateOrderItem(callCtx, &pborder_item.UpdateOrderItemRecordRequest{
			OrderItemId: int32(req.OrderItemID),
			Quantity:    int32(req.Quantity),
			Price:       int32(req.Price),
		})
	})
	if err != nil {
		return nil, order_item_errors.ErrUpdateOrderItem.WithInternal(err)
	}
	return toModel(res.Data), nil
}

func (r *Repository) Trash(ctx context.Context, orderID int) (*models.OrderItem, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pborder_item.ApiResponseOrderItem, error) {
		return r.command.TrashOrderItem(callCtx, &pborder_item.FindByIdOrderItemRequest{Id: int32(orderID)})
	})
	if err != nil {
		return nil, order_item_errors.ErrTrashedOrderItem.WithInternal(err)
	}
	return toModel(res.Data), nil
}

func (r *Repository) Restore(ctx context.Context, orderID int) (*models.OrderItem, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pborder_item.ApiResponseOrderItem, error) {
		return r.command.RestoreOrderItem(callCtx, &pborder_item.FindByIdOrderItemRequest{Id: int32(orderID)})
	})
	if err != nil {
		return nil, order_item_errors.ErrRestoreOrderItem.WithInternal(err)
	}
	return toModel(res.Data), nil
}

func (r *Repository) DeletePermanent(ctx context.Context, orderID int) (bool, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pborder_item.ApiResponseOrderItemDelete, error) {
		return r.command.DeleteOrderItemPermanent(callCtx, &pborder_item.FindByIdOrderItemRequest{Id: int32(orderID)})
	})
	if err != nil {
		return false, order_item_errors.ErrDeleteOrderItemPermanent.WithInternal(err)
	}
	return res.Status == "success", nil
}

func (r *Repository) DeleteByOrderIDPermanent(ctx context.Context, orderID int) (bool, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pborder_item.ApiResponseOrderItemDelete, error) {
		return r.command.DeleteOrderItemByOrderPermanent(callCtx, &pborder_item.FindByIdOrderItemRequest{Id: int32(orderID)})
	})
	if err != nil {
		return false, order_item_errors.ErrDeleteOrderItemPermanent.WithInternal(err)
	}
	return res.Status == "success", nil
}

func (r *Repository) RestoreAll(ctx context.Context) (bool, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pborder_item.ApiResponseOrderItemAll, error) {
		return r.command.RestoreAllOrdersItem(callCtx, &emptypb.Empty{})
	})
	if err != nil {
		return false, order_item_errors.ErrRestoreAllOrderItem.WithInternal(err)
	}
	return res.Status == "success", nil
}

func (r *Repository) DeleteAll(ctx context.Context) (bool, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pborder_item.ApiResponseOrderItemAll, error) {
		return r.command.DeleteAllPermanentOrdersItem(callCtx, &emptypb.Empty{})
	})
	if err != nil {
		return false, order_item_errors.ErrDeleteAllOrderPermanent.WithInternal(err)
	}
	return res.Status == "success", nil
}

func toModel(item *pborder_item.OrderItemResponse) *models.OrderItem {
	return &models.OrderItem{
		OrderItemID: item.Id,
		OrderID:     item.OrderId,
		ProductID:   item.ProductId,
		Quantity:    item.Quantity,
		Price:       item.Price,
	}
}

var (
	_ QueryRepository   = (*Repository)(nil)
	_ CommandRepository = (*Repository)(nil)
)
