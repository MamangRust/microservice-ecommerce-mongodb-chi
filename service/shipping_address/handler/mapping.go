package handler

import (
	pbcommon "github.com/MamangRust/microservice-ecommerce-grpc-pb/common"
	pbshipping_address "github.com/MamangRust/microservice-ecommerce-grpc-pb/shipping_address"

	"github.com/MamangRust/microservice-ecommerce-grpc-shipping-address/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/convert"
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

func mapToProtoShippingResponse(shipping interface{}) *pbshipping_address.ShippingResponse {
	switch s := shipping.(type) {
	case *models.ShippingAddress:
		return &pbshipping_address.ShippingResponse{
			Id:             s.ShippingAddressID,
			OrderId:        s.OrderID,
			Alamat:         s.Alamat,
			Provinsi:       s.Provinsi,
			Negara:         s.Negara,
			Kota:           s.Kota,
			ShippingMethod: s.ShippingMethod,
			ShippingCost:   int32(s.ShippingCost),
			CreatedAt:      convert.FormatTimePtr(s.CreatedAt),
			UpdatedAt:      convert.FormatTimePtr(s.UpdatedAt),
			Courier:        s.Courier,
		}
	case *repository.ShippingAddressResult:
		return &pbshipping_address.ShippingResponse{
			Id:             s.ShippingAddressID,
			OrderId:        s.OrderID,
			Alamat:         s.Alamat,
			Provinsi:       s.Provinsi,
			Negara:         s.Negara,
			Kota:           s.Kota,
			ShippingMethod: s.ShippingMethod,
			ShippingCost:   int32(s.ShippingCost),
			CreatedAt:      convert.StrVal(s.CreatedAt),
			UpdatedAt:      convert.StrVal(s.UpdatedAt),
			Courier:        s.Courier,
		}
	default:
		return nil
	}
}

func mapToProtoShippingResponseDeleteAt(shipping interface{}) *pbshipping_address.ShippingResponseDeleteAt {
	switch s := shipping.(type) {
	case *models.ShippingAddress:
		return &pbshipping_address.ShippingResponseDeleteAt{
			Id:             s.ShippingAddressID,
			OrderId:        s.OrderID,
			Alamat:         s.Alamat,
			Provinsi:       s.Provinsi,
			Negara:         s.Negara,
			Kota:           s.Kota,
			ShippingMethod: s.ShippingMethod,
			ShippingCost:   int32(s.ShippingCost),
			CreatedAt:      convert.FormatTimePtr(s.CreatedAt),
			UpdatedAt:      convert.FormatTimePtr(s.UpdatedAt),
			DeletedAt:      convert.TimeToWrappers(s.DeletedAt),
			Courier:        s.Courier,
		}
	case *repository.ShippingAddressResult:
		return &pbshipping_address.ShippingResponseDeleteAt{
			Id:             s.ShippingAddressID,
			OrderId:        s.OrderID,
			Alamat:         s.Alamat,
			Provinsi:       s.Provinsi,
			Negara:         s.Negara,
			Kota:           s.Kota,
			ShippingMethod: s.ShippingMethod,
			ShippingCost:   int32(s.ShippingCost),
			CreatedAt:      convert.StrVal(s.CreatedAt),
			UpdatedAt:      convert.StrVal(s.UpdatedAt),
			DeletedAt:      convert.StrValToWrappers(s.DeletedAt),
			Courier:        s.Courier,
		}
	default:
		return nil
	}
}
