package handler

import (
	"fmt"
	pbcommon "github.com/MamangRust/microservice-ecommerce-grpc-pb/common"
	pborder_item "github.com/MamangRust/microservice-ecommerce-grpc-pb/order_item"
	"time"

	"github.com/MamangRust/microservice-ecommerce-grpc-order-item/repository"
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
	totalPages := (totalRecords + pageSize - 1) / pageSize
	return &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}
}

func formatTimePtr(t *time.Time) string {
	if t != nil {
		return t.Format("2006-01-02 15:04:05.000")
	}
	return ""
}

func getStringValue(t *time.Time) *wrapperspb.StringValue {
	if t != nil {
		return wrapperspb.String(t.Format("2006-01-02 15:04:05.000"))
	}
	return nil
}

func mapToProtoOrderItemResponseFromModel(v *models.OrderItem) *pborder_item.OrderItemResponse {
	if v == nil {
		return nil
	}
	return &pborder_item.OrderItemResponse{
		Id:        v.OrderItemID,
		OrderId:   v.OrderID,
		ProductId: v.ProductID,
		Quantity:  v.Quantity,
		Price:     v.Price,
		CreatedAt: formatTimePtr(v.CreatedAt),
		UpdatedAt: formatTimePtr(v.UpdatedAt),
	}
}

func mapToProtoOrderItemResponseFromResult(v *repository.OrderItemResult) *pborder_item.OrderItemResponse {
	if v == nil {
		return nil
	}
	return &pborder_item.OrderItemResponse{
		Id:        v.OrderItemID,
		OrderId:   v.OrderID,
		ProductId: v.ProductID,
		Quantity:  v.Quantity,
		Price:     v.Price,
		CreatedAt: formatTimePtr(v.CreatedAt),
		UpdatedAt: formatTimePtr(v.UpdatedAt),
	}
}

func mapToProtoOrderItemResponseDeleteAtFromModel(v *models.OrderItem) *pborder_item.OrderItemResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pborder_item.OrderItemResponseDeleteAt{
		Id:        v.OrderItemID,
		OrderId:   v.OrderID,
		ProductId: v.ProductID,
		Quantity:  v.Quantity,
		Price:     v.Price,
		CreatedAt: formatTimePtr(v.CreatedAt),
		UpdatedAt: formatTimePtr(v.UpdatedAt),
	}
	if v.DeletedAt != nil {
		res.DeletedAt = &wrapperspb.StringValue{Value: formatTimePtr(v.DeletedAt)}
	}
	return res
}

func mapToProtoOrderItemResponseDeleteAtFromResult(v *repository.OrderItemResult) *pborder_item.OrderItemResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pborder_item.OrderItemResponseDeleteAt{
		Id:        v.OrderItemID,
		OrderId:   v.OrderID,
		ProductId: v.ProductID,
		Quantity:  v.Quantity,
		Price:     v.Price,
		CreatedAt: formatTimePtr(v.CreatedAt),
		UpdatedAt: formatTimePtr(v.UpdatedAt),
	}
	if v.DeletedAt != nil {
		res.DeletedAt = &wrapperspb.StringValue{Value: formatTimePtr(v.DeletedAt)}
	}
	return res
}

// Suppress unused import warning
var _ = fmt.Sprintf
