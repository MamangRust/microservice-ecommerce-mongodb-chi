package handler

import (
	pbcommon "github.com/MamangRust/microservice-ecommerce-grpc-pb/common"
	pbmerchant_business "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant_business"

	"github.com/MamangRust/microservice-ecommerce-grpc-merchant_business/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/convert"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func int32Deref(v *int32) int32 {
	if v == nil {
		return 0
	}
	return *v
}

func stringDeref(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func mapToProtoMerchantBusinessResponseFromResult(v *repository.MerchantBusinessResult) *pbmerchant_business.MerchantBusinessResponse {
	if v == nil {
		return nil
	}
	return &pbmerchant_business.MerchantBusinessResponse{
		Id:                v.MerchantBusinessInfoID,
		MerchantId:        v.MerchantID,
		BusinessType:      stringDeref(v.BusinessType),
		TaxId:             stringDeref(v.TaxID),
		EstablishedYear:   int32Deref(v.EstablishedYear),
		NumberOfEmployees: int32Deref(v.NumberOfEmployees),
		WebsiteUrl:        stringDeref(v.WebsiteUrl),
		CreatedAt:         convert.FormatTimePtr(v.CreatedAt),
		UpdatedAt:         convert.FormatTimePtr(v.UpdatedAt),
	}
}

func mapToProtoMerchantBusinessResponseFromModel(v *models.MerchantBusinessInformation) *pbmerchant_business.MerchantBusinessResponse {
	if v == nil {
		return nil
	}
	return &pbmerchant_business.MerchantBusinessResponse{
		Id:                v.MerchantBusinessInfoID,
		MerchantId:        v.MerchantID,
		BusinessType:      stringDeref(v.BusinessType),
		TaxId:             stringDeref(v.TaxID),
		EstablishedYear:   int32Deref(v.EstablishedYear),
		NumberOfEmployees: int32Deref(v.NumberOfEmployees),
		WebsiteUrl:        stringDeref(v.WebsiteUrl),
		CreatedAt:         convert.FormatTimePtr(v.CreatedAt),
		UpdatedAt:         convert.FormatTimePtr(v.UpdatedAt),
	}
}

func mapToProtoMerchantBusinessResponseDeleteAtFromResult(v *repository.MerchantBusinessResult) *pbmerchant_business.MerchantBusinessResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pbmerchant_business.MerchantBusinessResponseDeleteAt{
		Id:                v.MerchantBusinessInfoID,
		MerchantId:        v.MerchantID,
		BusinessType:      stringDeref(v.BusinessType),
		TaxId:             stringDeref(v.TaxID),
		EstablishedYear:   int32Deref(v.EstablishedYear),
		NumberOfEmployees: int32Deref(v.NumberOfEmployees),
		WebsiteUrl:        stringDeref(v.WebsiteUrl),
		CreatedAt:         convert.FormatTimePtr(v.CreatedAt),
		UpdatedAt:         convert.FormatTimePtr(v.UpdatedAt),
	}
	if v.DeletedAt != nil {
		res.DeletedAt = &wrapperspb.StringValue{Value: convert.FormatTimePtr(v.DeletedAt)}
	}
	return res
}

func mapToProtoMerchantBusinessResponseDeleteAtFromModel(v *models.MerchantBusinessInformation) *pbmerchant_business.MerchantBusinessResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pbmerchant_business.MerchantBusinessResponseDeleteAt{
		Id:                v.MerchantBusinessInfoID,
		MerchantId:        v.MerchantID,
		BusinessType:      stringDeref(v.BusinessType),
		TaxId:             stringDeref(v.TaxID),
		EstablishedYear:   int32Deref(v.EstablishedYear),
		NumberOfEmployees: int32Deref(v.NumberOfEmployees),
		WebsiteUrl:        stringDeref(v.WebsiteUrl),
		CreatedAt:         convert.FormatTimePtr(v.CreatedAt),
		UpdatedAt:         convert.FormatTimePtr(v.UpdatedAt),
	}
	if v.DeletedAt != nil {
		res.DeletedAt = &wrapperspb.StringValue{Value: convert.FormatTimePtr(v.DeletedAt)}
	}
	return res
}

func normalizePage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func mapToProtoMerchantBusinessResponse(v *repository.MerchantBusinessResult) *pbmerchant_business.MerchantBusinessResponse {
	return mapToProtoMerchantBusinessResponseFromResult(v)
}

func mapToProtoMerchantBusinessResponseDeleteAt(v *repository.MerchantBusinessResult) *pbmerchant_business.MerchantBusinessResponseDeleteAt {
	return mapToProtoMerchantBusinessResponseDeleteAtFromResult(v)
}

func createPaginationMeta(page, pageSize, totalRecords int) *pbcommon.PaginationMeta {
	totalPages := int32(totalRecords / pageSize)
	if totalRecords%pageSize != 0 {
		totalPages++
	}
	return &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   totalPages,
		TotalRecords: int32(totalRecords),
	}
}
