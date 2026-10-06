package handler

import (
	"context"
	pbtransaction "github.com/MamangRust/microservice-ecommerce-grpc-pb/transaction"

	"github.com/MamangRust/microservice-ecommerce-grpc-transaction/service"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

type transactionCommandHandler struct {
	pbtransaction.UnimplementedTransactionCommandServiceServer
	service service.TransactionCommandService
	logger  logger.LoggerInterface
}

func NewTransactionCommandHandler(service service.TransactionCommandService, logger logger.LoggerInterface) *transactionCommandHandler {
	return &transactionCommandHandler{
		service: service,
		logger:  logger,
	}
}

func (h *transactionCommandHandler) Create(ctx context.Context, req *pbtransaction.CreateTransactionRequest) (*pbtransaction.ApiResponseTransaction, error) {
	request := &requests.CreateTransactionRequest{
		OrderID:       int(req.GetOrderId()),
		MerchantID:    int(req.GetMerchantId()),
		UserID:        int(req.GetUserId()),
		PaymentMethod: req.GetPaymentMethod(),
		Amount:        int(req.GetAmount()),
		PaymentStatus: &[]string{req.GetPaymentStatus()}[0],
	}

	data, err := h.service.Create(ctx, request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbtransaction.ApiResponseTransaction{
		Status:  "success",
		Message: "Successfully created transaction",
		Data:    mapToProtoTransactionResponse(data),
	}, nil
}

func (h *transactionCommandHandler) Update(ctx context.Context, req *pbtransaction.UpdateTransactionRequest) (*pbtransaction.ApiResponseTransaction, error) {
	transactionID := int(req.GetTransactionId())
	request := &requests.UpdateTransactionRequest{
		TransactionID: &transactionID,
		MerchantID:    int(req.GetMerchantId()),
		OrderID:       int(req.GetOrderId()),
		PaymentMethod: req.GetPaymentMethod(),
		Amount:        int(req.GetAmount()),
		PaymentStatus: &[]string{req.GetPaymentStatus()}[0],
	}

	data, err := h.service.Update(ctx, request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbtransaction.ApiResponseTransaction{
		Status:  "success",
		Message: "Successfully updated transaction",
		Data:    mapToProtoTransactionResponse(data),
	}, nil
}

func (h *transactionCommandHandler) TrashedTransaction(ctx context.Context, req *pbtransaction.FindByIdTransactionRequest) (*pbtransaction.ApiResponseTransactionDeleteAt, error) {
	data, err := h.service.Trash(ctx, int(req.GetId()))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbtransaction.ApiResponseTransactionDeleteAt{
		Status:  "success",
		Message: "Successfully trashed transaction",
		Data:    mapToProtoTransactionResponseDeleteAt(data),
	}, nil
}

func (h *transactionCommandHandler) RestoreTransaction(ctx context.Context, req *pbtransaction.FindByIdTransactionRequest) (*pbtransaction.ApiResponseTransactionDeleteAt, error) {
	data, err := h.service.Restore(ctx, int(req.GetId()))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbtransaction.ApiResponseTransactionDeleteAt{
		Status:  "success",
		Message: "Successfully restored transaction",
		Data:    mapToProtoTransactionResponseDeleteAt(data),
	}, nil
}

func (h *transactionCommandHandler) DeleteTransactionPermanent(ctx context.Context, req *pbtransaction.FindByIdTransactionRequest) (*pbtransaction.ApiResponseTransactionDelete, error) {
	_, err := h.service.DeletePermanent(ctx, int(req.GetId()))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbtransaction.ApiResponseTransactionDelete{
		Status:  "success",
		Message: "Successfully deleted transaction permanently",
	}, nil
}

func (h *transactionCommandHandler) RestoreAllTransaction(ctx context.Context, req *emptypb.Empty) (*pbtransaction.ApiResponseTransactionAll, error) {
	_, err := h.service.RestoreAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbtransaction.ApiResponseTransactionAll{
		Status:  "success",
		Message: "Successfully restored all transactions",
	}, nil
}

func (h *transactionCommandHandler) DeleteTransactionByOrderPermanent(ctx context.Context, req *pbtransaction.FindByIdTransactionRequest) (*pbtransaction.ApiResponseTransactionDelete, error) {
	_, err := h.service.DeleteByOrderIDPermanent(ctx, int(req.GetId()))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbtransaction.ApiResponseTransactionDelete{
		Status:  "success",
		Message: "Successfully deleted transactions by order permanently",
	}, nil
}

func (h *transactionCommandHandler) DeleteAllTransactionPermanent(ctx context.Context, req *emptypb.Empty) (*pbtransaction.ApiResponseTransactionAll, error) {
	_, err := h.service.DeleteAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbtransaction.ApiResponseTransactionAll{
		Status:  "success",
		Message: "Successfully deleted all transactions permanently",
	}, nil
}
