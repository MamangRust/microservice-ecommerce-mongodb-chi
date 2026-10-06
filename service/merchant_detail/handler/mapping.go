package handler

import (
	pbcommon "github.com/MamangRust/microservice-ecommerce-grpc-pb/common"
	pbmerchant_detail "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant_detail"
	"math"

	"github.com/MamangRust/microservice-ecommerce-grpc-merchant_detail/repository"
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
	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))
	return &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}
}

func mapToProtoMerchantDetailResponseFromResult(v *repository.MerchantDetailResult) *pbmerchant_detail.MerchantDetailResponse {
	if v == nil {
		return nil
	}
	return &pbmerchant_detail.MerchantDetailResponse{
		Id:               v.MerchantDetailID,
		MerchantId:       v.MerchantID,
		DisplayName:      convert.StrVal(v.DisplayName),
		CoverImageUrl:    convert.StrVal(v.CoverImageUrl),
		LogoUrl:          convert.StrVal(v.LogoUrl),
		ShortDescription: convert.StrVal(v.ShortDescription),
		WebsiteUrl:       convert.StrVal(v.WebsiteUrl),
		CreatedAt:        convert.StrVal(v.CreatedAt),
		UpdatedAt:        convert.StrVal(v.UpdatedAt),
	}
}

func mapToProtoMerchantDetailResponseFromModel(v *models.MerchantDetail) *pbmerchant_detail.MerchantDetailResponse {
	if v == nil {
		return nil
	}
	return &pbmerchant_detail.MerchantDetailResponse{
		Id:               v.MerchantDetailID,
		MerchantId:       v.MerchantID,
		DisplayName:      convert.StrVal(v.DisplayName),
		CoverImageUrl:    convert.StrVal(v.CoverImageUrl),
		LogoUrl:          convert.StrVal(v.LogoUrl),
		ShortDescription: convert.StrVal(v.ShortDescription),
		WebsiteUrl:       convert.StrVal(v.WebsiteUrl),
		CreatedAt:        convert.FormatTimePtr(v.CreatedAt),
		UpdatedAt:        convert.FormatTimePtr(v.UpdatedAt),
	}
}

func mapToProtoMerchantDetailResponseDeleteAtFromResult(v *repository.MerchantDetailResult) *pbmerchant_detail.MerchantDetailResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pbmerchant_detail.MerchantDetailResponseDeleteAt{
		Id:               v.MerchantDetailID,
		MerchantId:       v.MerchantID,
		DisplayName:      convert.StrVal(v.DisplayName),
		CoverImageUrl:    convert.StrVal(v.CoverImageUrl),
		LogoUrl:          convert.StrVal(v.LogoUrl),
		ShortDescription: convert.StrVal(v.ShortDescription),
		WebsiteUrl:       convert.StrVal(v.WebsiteUrl),
		CreatedAt:        convert.StrVal(v.CreatedAt),
		UpdatedAt:        convert.StrVal(v.UpdatedAt),
	}
	if v.DeletedAt != nil && *v.DeletedAt != "" {
		res.DeletedAt = convert.StrValToWrappers(v.DeletedAt)
	}
	return res
}

func mapToProtoMerchantDetailResponseDeleteAtFromModel(v *models.MerchantDetail) *pbmerchant_detail.MerchantDetailResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pbmerchant_detail.MerchantDetailResponseDeleteAt{
		Id:               v.MerchantDetailID,
		MerchantId:       v.MerchantID,
		DisplayName:      convert.StrVal(v.DisplayName),
		CoverImageUrl:    convert.StrVal(v.CoverImageUrl),
		LogoUrl:          convert.StrVal(v.LogoUrl),
		ShortDescription: convert.StrVal(v.ShortDescription),
		WebsiteUrl:       convert.StrVal(v.WebsiteUrl),
		CreatedAt:        convert.FormatTimePtr(v.CreatedAt),
		UpdatedAt:        convert.FormatTimePtr(v.UpdatedAt),
	}
	if v.DeletedAt != nil {
		res.DeletedAt = convert.TimeToWrappers(v.DeletedAt)
	}
	return res
}

func mapToProtoMerchantDetailResponse(m interface{}) *pbmerchant_detail.MerchantDetailResponse {
	switch v := m.(type) {
	case *repository.MerchantDetailResult:
		return mapToProtoMerchantDetailResponseFromResult(v)
	case *models.MerchantDetail:
		return mapToProtoMerchantDetailResponseFromModel(v)
	default:
		return nil
	}
}

func mapToProtoMerchantDetailResponseDeleteAt(m interface{}) *pbmerchant_detail.MerchantDetailResponseDeleteAt {
	switch v := m.(type) {
	case *repository.MerchantDetailResult:
		return mapToProtoMerchantDetailResponseDeleteAtFromResult(v)
	case *models.MerchantDetail:
		return mapToProtoMerchantDetailResponseDeleteAtFromModel(v)
	default:
		return nil
	}
}

func mapToProtoMerchantSocialLinkResponse(m interface{}) *pbmerchant_detail.MerchantSocialMediaLinkResponse {
	switch v := m.(type) {
	case *models.MerchantSocialMediaLink:
		return &pbmerchant_detail.MerchantSocialMediaLinkResponse{
			Id:               v.MerchantSocialID,
			MerchantDetailId: v.MerchantDetailID,
			Platform:         v.Platform,
			Url:              v.Url,
			CreatedAt:        convert.FormatTimePtr(v.CreatedAt),
			UpdatedAt:        convert.FormatTimePtr(v.UpdatedAt),
		}
	default:
		return nil
	}
}
