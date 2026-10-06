package handler

import (
	pbcommon "github.com/MamangRust/microservice-ecommerce-grpc-pb/common"
	pbmerchant_award "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant_award"
	"math"

	"github.com/MamangRust/microservice-ecommerce-grpc-merchant_award/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/convert"
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

func mapToProtoMerchantAwardResponse(v *repository.MerchantCertResult) *pbmerchant_award.MerchantAwardResponse {
	return mapToProtoMerchantAwardResponseFromResult(v)
}

func mapToProtoMerchantAwardResponseDeleteAt(v *repository.MerchantCertResult) *pbmerchant_award.MerchantAwardResponseDeleteAt {
	return mapToProtoMerchantAwardResponseDeleteAtFromResult(v)
}

func mapToProtoMerchantAwardResponseFromResult(v *repository.MerchantCertResult) *pbmerchant_award.MerchantAwardResponse {
	if v == nil {
		return nil
	}
	return &pbmerchant_award.MerchantAwardResponse{
		Id:             v.MerchantCertificationID,
		MerchantId:     v.MerchantID,
		Title:          v.Title,
		Description:    convert.StrVal(v.Description),
		IssuedBy:       convert.StrVal(v.IssuedBy),
		CertificateUrl: convert.StrVal(v.CertificateUrl),
		IssueDate:      convert.FormatTimePtr(v.IssueDate),
		ExpiryDate:     convert.FormatTimePtr(v.ExpiryDate),
		CreatedAt:      convert.FormatTimePtr(v.CreatedAt),
		UpdatedAt:      convert.FormatTimePtr(v.UpdatedAt),
	}
}

func mapToProtoMerchantAwardResponseFromModel(v *models.MerchantCertificationsAndAward) *pbmerchant_award.MerchantAwardResponse {
	if v == nil {
		return nil
	}
	return &pbmerchant_award.MerchantAwardResponse{
		Id:             v.MerchantCertificationID,
		MerchantId:     v.MerchantID,
		Title:          v.Title,
		Description:    convert.StrVal(v.Description),
		IssuedBy:       convert.StrVal(v.IssuedBy),
		CertificateUrl: convert.StrVal(v.CertificateUrl),
		IssueDate:      convert.FormatTimePtr(v.IssueDate),
		ExpiryDate:     convert.FormatTimePtr(v.ExpiryDate),
		CreatedAt:      convert.FormatTimePtr(v.CreatedAt),
		UpdatedAt:      convert.FormatTimePtr(v.UpdatedAt),
	}
}

func mapToProtoMerchantAwardResponseDeleteAtFromResult(v *repository.MerchantCertResult) *pbmerchant_award.MerchantAwardResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pbmerchant_award.MerchantAwardResponseDeleteAt{
		Id:             v.MerchantCertificationID,
		MerchantId:     v.MerchantID,
		Title:          v.Title,
		Description:    convert.StrVal(v.Description),
		IssuedBy:       convert.StrVal(v.IssuedBy),
		CertificateUrl: convert.StrVal(v.CertificateUrl),
		IssueDate:      convert.FormatTimePtr(v.IssueDate),
		ExpiryDate:     convert.FormatTimePtr(v.ExpiryDate),
		CreatedAt:      convert.FormatTimePtr(v.CreatedAt),
		UpdatedAt:      convert.FormatTimePtr(v.UpdatedAt),
	}
	if v.DeletedAt != nil {
		res.DeletedAt = &wrapperspb.StringValue{Value: convert.FormatTimePtr(v.DeletedAt)}
	}
	return res
}

func mapToProtoMerchantAwardResponseDeleteAtFromModel(v *models.MerchantCertificationsAndAward) *pbmerchant_award.MerchantAwardResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pbmerchant_award.MerchantAwardResponseDeleteAt{
		Id:             v.MerchantCertificationID,
		MerchantId:     v.MerchantID,
		Title:          v.Title,
		Description:    convert.StrVal(v.Description),
		IssuedBy:       convert.StrVal(v.IssuedBy),
		CertificateUrl: convert.StrVal(v.CertificateUrl),
		IssueDate:      convert.FormatTimePtr(v.IssueDate),
		ExpiryDate:     convert.FormatTimePtr(v.ExpiryDate),
		CreatedAt:      convert.FormatTimePtr(v.CreatedAt),
		UpdatedAt:      convert.FormatTimePtr(v.UpdatedAt),
	}
	if v.DeletedAt != nil {
		res.DeletedAt = &wrapperspb.StringValue{Value: convert.FormatTimePtr(v.DeletedAt)}
	}
	return res
}
