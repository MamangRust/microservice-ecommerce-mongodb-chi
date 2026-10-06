package handler

import (
	"context"
	pbtransaction "github.com/MamangRust/microservice-ecommerce-grpc-pb/transaction"

	"github.com/MamangRust/microservice-ecommerce-grpc-transaction/service"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors"
)

type transactionQueryHandler struct {
	pbtransaction.UnimplementedTransactionQueryServiceServer
	service service.TransactionQueryService
	logger  logger.LoggerInterface
}

func NewTransactionQueryHandler(service service.TransactionQueryService, logger logger.LoggerInterface) *transactionQueryHandler {
	return &transactionQueryHandler{
		service: service,
		logger:  logger,
	}
}

func (h *transactionQueryHandler) FindAllTransactions(ctx context.Context, req *pbtransaction.FindAllTransactionRequest) (*pbtransaction.ApiResponsePaginationTransaction, error) {
	request := &requests.FindAllTransaction{
		Page:     int(req.GetPage()),
		PageSize: int(req.GetPageSize()),
		Search:   req.GetSearch(),
	}

	data, total, err := h.service.FindAll(ctx, request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var transactions []*pbtransaction.TransactionResponse
	for _, v := range data {
		transactions = append(transactions, mapToProtoTransactionResponse(v))
	}

	return &pbtransaction.ApiResponsePaginationTransaction{
		Status:     "success",
		Message:    "Successfully fetched transactions",
		Data:       transactions,
		Pagination: createPaginationMeta(request.Page, request.PageSize, *total),
	}, nil
}

func (h *transactionQueryHandler) FindByActive(ctx context.Context, req *pbtransaction.FindAllTransactionRequest) (*pbtransaction.ApiResponsePaginationTransaction, error) {
	request := &requests.FindAllTransaction{
		Page:     int(req.GetPage()),
		PageSize: int(req.GetPageSize()),
		Search:   req.GetSearch(),
	}

	data, total, err := h.service.FindActive(ctx, request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var transactions []*pbtransaction.TransactionResponse
	for _, v := range data {
		transactions = append(transactions, mapToProtoTransactionResponse(v))
	}

	return &pbtransaction.ApiResponsePaginationTransaction{
		Status:     "success",
		Message:    "Successfully fetched active transactions",
		Data:       transactions,
		Pagination: createPaginationMeta(request.Page, request.PageSize, *total),
	}, nil
}

func (h *transactionQueryHandler) FindByTrashed(ctx context.Context, req *pbtransaction.FindAllTransactionRequest) (*pbtransaction.ApiResponsePaginationTransactionDeleteAt, error) {
	request := &requests.FindAllTransaction{
		Page:     int(req.GetPage()),
		PageSize: int(req.GetPageSize()),
		Search:   req.GetSearch(),
	}

	data, total, err := h.service.FindTrashed(ctx, request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var transactions []*pbtransaction.TransactionResponseDeleteAt
	for _, v := range data {
		transactions = append(transactions, mapToProtoTransactionResponseDeleteAt(v))
	}

	return &pbtransaction.ApiResponsePaginationTransactionDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed transactions",
		Data:       transactions,
		Pagination: createPaginationMeta(request.Page, request.PageSize, *total),
	}, nil
}

func (h *transactionQueryHandler) FindById(ctx context.Context, req *pbtransaction.FindByIdTransactionRequest) (*pbtransaction.ApiResponseTransaction, error) {
	data, err := h.service.FindByID(ctx, int(req.GetId()))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbtransaction.ApiResponseTransaction{
		Status:  "success",
		Message: "Successfully fetched transaction",
		Data:    mapToProtoTransactionResponse(data),
	}, nil
}

func (h *transactionQueryHandler) FindByOrderId(ctx context.Context, req *pbtransaction.FindByOrderIdTransactionRequest) (*pbtransaction.ApiResponseTransaction, error) {
	data, err := h.service.FindByOrderID(ctx, int(req.GetOrderId()))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbtransaction.ApiResponseTransaction{
		Status:  "success",
		Message: "Successfully fetched transaction by order id",
		Data:    mapToProtoTransactionResponse(data),
	}, nil
}

func (h *transactionQueryHandler) FindByMerchant(ctx context.Context, req *pbtransaction.FindAllTransactionByMerchantRequest) (*pbtransaction.ApiResponsePaginationTransaction, error) {
	request := &requests.FindAllTransactionByMerchant{
		MerchantID: int(req.GetMerchantId()),
		Page:       int(req.GetPage()),
		PageSize:   int(req.GetPageSize()),
		Search:     req.GetSearch(),
	}

	data, total, err := h.service.FindByMerchant(ctx, request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var transactions []*pbtransaction.TransactionResponse
	for _, v := range data {
		transactions = append(transactions, mapToProtoTransactionResponse(v))
	}

	return &pbtransaction.ApiResponsePaginationTransaction{
		Status:     "success",
		Message:    "Successfully fetched transactions by merchant",
		Data:       transactions,
		Pagination: createPaginationMeta(request.Page, request.PageSize, *total),
	}, nil
}
