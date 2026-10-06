package handler

import (
	"context"
	pborder "github.com/MamangRust/microservice-ecommerce-grpc-pb/order"

	"github.com/MamangRust/microservice-ecommerce-grpc-order/service"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/order_errors"
)

type orderQueryHandler struct {
	pborder.UnimplementedOrderQueryServiceServer
	orderQuery service.OrderQueryService
	logger     logger.LoggerInterface
}

func NewOrderQueryHandler(
	orderQuery service.OrderQueryService,
	logger logger.LoggerInterface,
) pborder.OrderQueryServiceServer {
	return &orderQueryHandler{
		orderQuery: orderQuery,
		logger:     logger,
	}
}

func (s *orderQueryHandler) FindAll(ctx context.Context, request *pborder.FindAllOrderRequest) (*pborder.ApiResponsePaginationOrder, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllOrder{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orders, totalRecords, err := s.orderQuery.FindAll(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbOrders := make([]*pborder.OrderResponse, len(orders))
	for i, order := range orders {
		pbOrders[i] = mapToProtoOrderResponseFromResult(order)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pborder.ApiResponsePaginationOrder{
		Status:     "success",
		Message:    "Successfully fetched order",
		Data:       pbOrders,
		Pagination: paginationMeta,
	}, nil
}

func (s *orderQueryHandler) FindById(ctx context.Context, request *pborder.FindByIdOrderRequest) (*pborder.ApiResponseOrder, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	order, err := s.orderQuery.FindByID(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pborder.ApiResponseOrder{
		Status:  "success",
		Message: "Successfully fetched order",
		Data:    mapToProtoOrderResponseFromModel(order),
	}, nil
}

func (s *orderQueryHandler) FindByActive(ctx context.Context, request *pborder.FindAllOrderRequest) (*pborder.ApiResponsePaginationOrderDeleteAt, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllOrder{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orders, totalRecords, err := s.orderQuery.FindActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbOrders := make([]*pborder.OrderResponseDeleteAt, len(orders))
	for i, order := range orders {
		pbOrders[i] = mapToProtoOrderResponseDeleteAtFromResult(order)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pborder.ApiResponsePaginationOrderDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active order",
		Data:       pbOrders,
		Pagination: paginationMeta,
	}, nil
}

func (s *orderQueryHandler) FindByTrashed(ctx context.Context, request *pborder.FindAllOrderRequest) (*pborder.ApiResponsePaginationOrderDeleteAt, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllOrder{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orders, totalRecords, err := s.orderQuery.FindTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbOrders := make([]*pborder.OrderResponseDeleteAt, len(orders))
	for i, order := range orders {
		pbOrders[i] = mapToProtoOrderResponseDeleteAtFromResult(order)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pborder.ApiResponsePaginationOrderDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed order",
		Data:       pbOrders,
		Pagination: paginationMeta,
	}, nil
}
