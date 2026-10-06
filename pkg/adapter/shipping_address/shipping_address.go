package shipping_address

import (
	"context"

	pbshipping "github.com/MamangRust/microservice-ecommerce-grpc-pb/shipping_address"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-pkg/resilience"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	shippingaddress_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/shipping_address_errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

// QueryRepository is the contract consumers depend on for shipping-address reads.
type QueryRepository interface {
	FindByID(ctx context.Context, shippingID int) (*models.ShippingAddress, error)
	FindByOrder(ctx context.Context, orderID int) (*models.ShippingAddress, error)
}

// CommandRepository is the contract consumers depend on for shipping-address writes.
type CommandRepository interface {
	Create(ctx context.Context, req *requests.CreateShippingAddressRequest) (*models.ShippingAddress, error)
	Update(ctx context.Context, req *requests.UpdateShippingAddressRequest) (*models.ShippingAddress, error)
	DeleteByOrderIDPermanent(ctx context.Context, orderID int) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

// Repository implements both QueryRepository and CommandRepository on top of the
// generated shipping query/command clients, guarding every call.
type Repository struct {
	query   pbshipping.ShippingQueryServiceClient
	command pbshipping.ShippingCommandServiceClient
	guard   *resilience.DependencyGuard
}

// New builds the shipping-address adapter from its query and command clients.
func New(query pbshipping.ShippingQueryServiceClient, command pbshipping.ShippingCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{query: query, command: command}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// NewAdapter wraps the generated shipping clients.
func NewAdapter(query pbshipping.ShippingQueryServiceClient, command pbshipping.ShippingCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	return New(query, command, opts...)
}

// NewQueryAdapter returns the adapter restricted to the QueryRepository surface.
func NewQueryAdapter(query pbshipping.ShippingQueryServiceClient, opts ...adapter.GuardOption) *Repository {
	return New(query, nil, opts...)
}

// NewCommandAdapter returns the adapter restricted to the CommandRepository surface.
func NewCommandAdapter(command pbshipping.ShippingCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	return New(nil, command, opts...)
}

// SetGuard attaches the dependency guard used to wrap every outbound call.
func (r *Repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

// Reads propagate the dependency's gRPC status unchanged (NotFound -> 404, and
// so on). Command failures keep the domain-error wrapping the order service
// already relied on; AppError.Unwrap still exposes the underlying gRPC status.
func (r *Repository) FindByID(ctx context.Context, shippingID int) (*models.ShippingAddress, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbshipping.ApiResponseShipping, error) {
		return r.query.FindById(callCtx, &pbshipping.FindByIdShippingRequest{Id: int32(shippingID)})
	})
	if err != nil {
		return nil, err
	}
	return toModel(res.Data), nil
}

func (r *Repository) FindByOrder(ctx context.Context, orderID int) (*models.ShippingAddress, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbshipping.ApiResponseShipping, error) {
		return r.query.FindByOrder(callCtx, &pbshipping.FindByIdShippingRequest{Id: int32(orderID)})
	})
	if err != nil {
		return nil, err
	}
	return toModel(res.Data), nil
}

func (r *Repository) Create(ctx context.Context, req *requests.CreateShippingAddressRequest) (*models.ShippingAddress, error) {
	var orderID int32
	if req.OrderID != nil {
		orderID = int32(*req.OrderID)
	}

	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbshipping.ApiResponseShipping, error) {
		return r.command.CreateShipping(callCtx, &pbshipping.CreateShippingAddressRequest{
			OrderId:        orderID,
			Alamat:         req.Alamat,
			Provinsi:       req.Provinsi,
			Kota:           req.Kota,
			Negara:         req.Negara,
			Courier:        req.Courier,
			ShippingMethod: req.ShippingMethod,
			ShippingCost:   int32(req.ShippingCost),
		})
	})
	if err != nil {
		return nil, shippingaddress_errors.ErrCreateShippingAddress.WithInternal(err)
	}
	return toModel(res.Data), nil
}

func (r *Repository) Update(ctx context.Context, req *requests.UpdateShippingAddressRequest) (*models.ShippingAddress, error) {
	var shippingID int32
	if req.ShippingID != nil {
		shippingID = int32(*req.ShippingID)
	}

	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbshipping.ApiResponseShipping, error) {
		return r.command.UpdateShipping(callCtx, &pbshipping.UpdateShippingAddressRequest{
			ShippingId:     shippingID,
			Alamat:         req.Alamat,
			Provinsi:       req.Provinsi,
			Kota:           req.Kota,
			Negara:         req.Negara,
			Courier:        req.Courier,
			ShippingMethod: req.ShippingMethod,
			ShippingCost:   int32(req.ShippingCost),
		})
	})
	if err != nil {
		return nil, shippingaddress_errors.ErrUpdateShippingAddress.WithInternal(err)
	}
	return toModel(res.Data), nil
}

func (r *Repository) DeleteByOrderIDPermanent(ctx context.Context, orderID int) (bool, error) {
	_, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbshipping.ApiResponseShippingDelete, error) {
		return r.command.DeleteShippingByOrderPermanent(callCtx, &pbshipping.FindByIdShippingRequest{Id: int32(orderID)})
	})
	if err != nil {
		return false, shippingaddress_errors.ErrDeleteShippingAddressPermanent.WithInternal(err)
	}
	return true, nil
}

func (r *Repository) DeleteAll(ctx context.Context) (bool, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbshipping.ApiResponseShippingAll, error) {
		return r.command.DeleteAllShippingPermanent(callCtx, &emptypb.Empty{})
	})
	if err != nil {
		return false, shippingaddress_errors.ErrDeleteAllPermanentShippingAddress.WithInternal(err)
	}
	return res.Status == "success", nil
}

func toModel(data *pbshipping.ShippingResponse) *models.ShippingAddress {
	if data == nil {
		return nil
	}
	return &models.ShippingAddress{
		ShippingAddressID: data.Id,
		OrderID:           data.OrderId,
		Alamat:            data.Alamat,
		Provinsi:          data.Provinsi,
		Negara:            data.Negara,
		Kota:              data.Kota,
		Courier:           data.Courier,
		ShippingMethod:    data.ShippingMethod,
		ShippingCost:      float64(data.ShippingCost),
	}
}

var (
	_ QueryRepository   = (*Repository)(nil)
	_ CommandRepository = (*Repository)(nil)
)
