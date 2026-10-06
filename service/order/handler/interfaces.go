package handler

import (
	"context"
	pborder "github.com/MamangRust/microservice-ecommerce-grpc-pb/order"

	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderQueryHandler interface {
	pborder.OrderQueryServiceServer
}

type OrderCommandHandler interface {
	pborder.OrderCommandServiceServer
}

type OrderHandleGrpc interface {
	FindAll(ctx context.Context, request *pborder.FindAllOrderRequest) (*pborder.ApiResponsePaginationOrder, error)
	FindById(ctx context.Context, request *pborder.FindByIdOrderRequest) (*pborder.ApiResponseOrder, error)

	FindByActive(ctx context.Context, request *pborder.FindAllOrderRequest) (*pborder.ApiResponsePaginationOrderDeleteAt, error)
	FindByTrashed(ctx context.Context, request *pborder.FindAllOrderRequest) (*pborder.ApiResponsePaginationOrderDeleteAt, error)

	Create(ctx context.Context, request *pborder.CreateOrderRequest) (*pborder.ApiResponseOrder, error)
	Update(ctx context.Context, request *pborder.UpdateOrderRequest) (*pborder.ApiResponseOrder, error)
	TrashedOrder(ctx context.Context, request *pborder.FindByIdOrderRequest) (*pborder.ApiResponseOrderDeleteAt, error)
	RestoreOrder(ctx context.Context, request *pborder.FindByIdOrderRequest) (*pborder.ApiResponseOrderDeleteAt, error)
	DeleteOrderPermanent(ctx context.Context, request *pborder.FindByIdOrderRequest) (*pborder.ApiResponseOrderDelete, error)
	RestoreAllOrder(ctx context.Context, _ *emptypb.Empty) (*pborder.ApiResponseOrderAll, error)
	DeleteAllOrderPermanent(ctx context.Context, _ *emptypb.Empty) (*pborder.ApiResponseOrderAll, error)
}
