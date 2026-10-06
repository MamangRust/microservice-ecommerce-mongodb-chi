package handler

import (
	pbcommon "github.com/MamangRust/microservice-ecommerce-grpc-pb/common"
	pbtransaction "github.com/MamangRust/microservice-ecommerce-grpc-pb/transaction"
	"github.com/MamangRust/microservice-ecommerce-shared/convert"

	"github.com/MamangRust/microservice-ecommerce-grpc-transaction/repository"
)

func normalizePage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func createPaginationMeta(page, pageSize, totalRecords int) *pbcommon.PaginationMeta {
	totalPages := (totalRecords + pageSize - 1) / pageSize
	return &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}
}

func mapToProtoTransactionResponse(v *repository.TransactionResult) *pbtransaction.TransactionResponse {
	if v == nil {
		return nil
	}
	return &pbtransaction.TransactionResponse{
		Id:            v.TransactionID,
		OrderId:       v.OrderID,
		MerchantId:    v.MerchantID,
		PaymentMethod: v.PaymentMethod,
		Amount:        v.Amount,
		PaymentStatus: v.PaymentStatus,
		CreatedAt:     convert.FormatTimePtr(v.CreatedAt),
		UpdatedAt:     convert.FormatTimePtr(v.UpdatedAt),
	}
}

func mapToProtoTransactionResponseDeleteAt(v *repository.TransactionResult) *pbtransaction.TransactionResponseDeleteAt {
	if v == nil {
		return nil
	}
	return &pbtransaction.TransactionResponseDeleteAt{
		Id:            v.TransactionID,
		OrderId:       v.OrderID,
		MerchantId:    v.MerchantID,
		PaymentMethod: v.PaymentMethod,
		Amount:        v.Amount,
		PaymentStatus: v.PaymentStatus,
		CreatedAt:     convert.FormatTimePtr(v.CreatedAt),
		UpdatedAt:     convert.FormatTimePtr(v.UpdatedAt),
		DeletedAt:     convert.TimeToWrappers(v.DeletedAt),
	}
}
