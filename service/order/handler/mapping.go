package handler

import (
	pbcommon "github.com/MamangRust/microservice-ecommerce-grpc-pb/common"
	pborder "github.com/MamangRust/microservice-ecommerce-grpc-pb/order"
	"math"
	"time"

	"github.com/MamangRust/microservice-ecommerce-grpc-order/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"google.golang.org/protobuf/types/known/wrapperspb"
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
	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))
	return &pbcommon.PaginationMeta{CurrentPage: int32(page), PageSize: int32(pageSize), TotalPages: int32(totalPages), TotalRecords: int32(totalRecords)}
}

func formatTimePtr(t *time.Time) string {
	if t != nil {
		return t.Format("2006-01-02 15:04:05.000")
	}
	return ""
}

func mapToProtoOrderResponseFromModel(v *models.Order) *pborder.OrderResponse {
	if v == nil {
		return nil
	}
	return &pborder.OrderResponse{Id: v.OrderID, MerchantId: v.MerchantID, UserId: v.UserID, TotalPrice: int32(v.TotalPrice), CreatedAt: formatTimePtr(v.CreatedAt), UpdatedAt: formatTimePtr(v.UpdatedAt)}
}

func mapToProtoOrderResponseFromResult(v *repository.OrderResult) *pborder.OrderResponse {
	if v == nil {
		return nil
	}
	return &pborder.OrderResponse{Id: v.OrderID, MerchantId: v.MerchantID, UserId: v.UserID, TotalPrice: int32(v.TotalPrice), CreatedAt: formatTimePtr(v.CreatedAt), UpdatedAt: formatTimePtr(v.UpdatedAt)}
}

func mapToProtoOrderResponseDeleteAtFromModel(v *models.Order) *pborder.OrderResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pborder.OrderResponseDeleteAt{Id: v.OrderID, MerchantId: v.MerchantID, UserId: v.UserID, TotalPrice: int32(v.TotalPrice), CreatedAt: formatTimePtr(v.CreatedAt), UpdatedAt: formatTimePtr(v.UpdatedAt)}
	if v.DeletedAt != nil {
		res.DeletedAt = &wrapperspb.StringValue{Value: formatTimePtr(v.DeletedAt)}
	}
	return res
}

func mapToProtoOrderResponseDeleteAtFromResult(v *repository.OrderResult) *pborder.OrderResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pborder.OrderResponseDeleteAt{Id: v.OrderID, MerchantId: v.MerchantID, UserId: v.UserID, TotalPrice: int32(v.TotalPrice), CreatedAt: formatTimePtr(v.CreatedAt), UpdatedAt: formatTimePtr(v.UpdatedAt)}
	if v.DeletedAt != nil {
		res.DeletedAt = &wrapperspb.StringValue{Value: formatTimePtr(v.DeletedAt)}
	}
	return res
}
