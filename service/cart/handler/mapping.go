package handler

import (
	pbcart "github.com/MamangRust/microservice-ecommerce-grpc-pb/cart"
	pbcommon "github.com/MamangRust/microservice-ecommerce-grpc-pb/common"
	"math"

	"github.com/MamangRust/microservice-ecommerce-grpc-cart/repository"
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
	return &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}
}

func mapToProtoCartResponseFromResult(v *repository.CartResult) *pbcart.CartResponse {
	return &pbcart.CartResponse{
		Id:        v.CartID,
		UserId:    v.UserID,
		ProductId: v.ProductID,
		Name:      v.Name,
		Price:     v.Price,
		Image:     v.Image,
		Quantity:  v.Quantity,
		Weight:    v.Weight,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}
}

func mapToProtoCartResponseFromCreate(v *repository.CartCreateResult) *pbcart.CartResponse {
	return &pbcart.CartResponse{
		Id:        v.CartID,
		UserId:    v.UserID,
		ProductId: v.ProductID,
		Name:      v.Name,
		Price:     v.Price,
		Image:     v.Image,
		Quantity:  v.Quantity,
		Weight:    v.Weight,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}
}
