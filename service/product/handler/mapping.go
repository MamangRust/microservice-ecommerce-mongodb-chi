package handler

import (
	pbcommon "github.com/MamangRust/microservice-ecommerce-grpc-pb/common"
	pbproduct "github.com/MamangRust/microservice-ecommerce-grpc-pb/product"
	"math"

	"github.com/MamangRust/microservice-ecommerce-grpc-product/repository"
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
	return &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}
}

func stringPtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func float64PtrToFloat32(f *float64) float32 {
	if f == nil {
		return 0
	}
	return float32(*f)
}

func int32PtrToInt32(i *int32) int32 {
	if i == nil {
		return 0
	}
	return *i
}

func formatTimePtr(t interface{}) string {
	switch v := t.(type) {
	case *string:
		if v != nil {
			return *v
		}
	case string:
		return v
	}
	return ""
}

func mapToProtoProductResponse(item interface{}) *pbproduct.ProductResponse {
	switch v := item.(type) {
	case *models.Product:
		return &pbproduct.ProductResponse{
			Id:           v.ProductID,
			MerchantId:   v.MerchantID,
			CategoryId:   v.CategoryID,
			Name:         v.Name,
			Description:  stringPtrToString(v.Description),
			Price:        v.Price,
			CountInStock: v.CountInStock,
			Brand:        stringPtrToString(v.Brand),
			Weight:       int32PtrToInt32(v.Weight),
			Rating:       float64PtrToFloat32(v.Rating),
			SlugProduct:  stringPtrToString(v.SlugProduct),
			ImageProduct: stringPtrToString(v.ImageProduct),
			CreatedAt:    formatTimePtr(v.CreatedAt),
			UpdatedAt:    formatTimePtr(v.UpdatedAt),
		}
	case *repository.ProductResult:
		return &pbproduct.ProductResponse{
			Id:           v.ProductID,
			MerchantId:   v.MerchantID,
			CategoryId:   v.CategoryID,
			Name:         v.Name,
			Description:  stringPtrToString(v.Description),
			Price:        v.Price,
			CountInStock: v.CountInStock,
			Brand:        stringPtrToString(v.Brand),
			Weight:       int32PtrToInt32(v.Weight),
			Rating:       float64PtrToFloat32(v.Rating),
			SlugProduct:  stringPtrToString(v.SlugProduct),
			ImageProduct: stringPtrToString(v.ImageProduct),
			CreatedAt:    formatTimePtr(v.CreatedAt),
			UpdatedAt:    formatTimePtr(v.UpdatedAt),
		}
	default:
		return nil
	}
}

func mapToProtoProductResponseDeleteAt(item interface{}) *pbproduct.ProductResponseDeleteAt {
	var res *pbproduct.ProductResponseDeleteAt
	var deletedAt interface{}

	switch v := item.(type) {
	case *models.Product:
		res = &pbproduct.ProductResponseDeleteAt{
			Id:           v.ProductID,
			MerchantId:   v.MerchantID,
			CategoryId:   v.CategoryID,
			Name:         v.Name,
			Description:  stringPtrToString(v.Description),
			Price:        v.Price,
			CountInStock: v.CountInStock,
			Brand:        stringPtrToString(v.Brand),
			Weight:       int32PtrToInt32(v.Weight),
			Rating:       float64PtrToFloat32(v.Rating),
			SlugProduct:  stringPtrToString(v.SlugProduct),
			ImageProduct: stringPtrToString(v.ImageProduct),
			CreatedAt:    formatTimePtr(v.CreatedAt),
			UpdatedAt:    formatTimePtr(v.UpdatedAt),
		}
		deletedAt = v.DeletedAt
	case *repository.ProductResult:
		res = &pbproduct.ProductResponseDeleteAt{
			Id:           v.ProductID,
			MerchantId:   v.MerchantID,
			CategoryId:   v.CategoryID,
			Name:         v.Name,
			Description:  stringPtrToString(v.Description),
			Price:        v.Price,
			CountInStock: v.CountInStock,
			Brand:        stringPtrToString(v.Brand),
			Weight:       int32PtrToInt32(v.Weight),
			Rating:       float64PtrToFloat32(v.Rating),
			SlugProduct:  stringPtrToString(v.SlugProduct),
			ImageProduct: stringPtrToString(v.ImageProduct),
			CreatedAt:    formatTimePtr(v.CreatedAt),
			UpdatedAt:    formatTimePtr(v.UpdatedAt),
		}
		deletedAt = v.DeletedAt
	default:
		return nil
	}

	if val := formatTimePtr(deletedAt); val != "" {
		res.DeletedAt = &wrapperspb.StringValue{Value: val}
	}

	return res
}
