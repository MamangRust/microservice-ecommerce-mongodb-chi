package handler

import (
	pbcommon "github.com/MamangRust/microservice-ecommerce-grpc-pb/common"
	pbmerchant "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant"
	pbmerchant_document "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant_document"
	"math"

	"github.com/MamangRust/microservice-ecommerce-grpc-merchant/repository"
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

func mapToProtoMerchantResponseFromResult(v *repository.MerchantResult) *pbmerchant.MerchantResponse {
	if v == nil {
		return nil
	}
	return &pbmerchant.MerchantResponse{
		Id:           v.MerchantID,
		UserId:       v.UserID,
		Name:         v.Name,
		Description:  convert.StrVal(v.Description),
		Address:      convert.StrVal(v.Address),
		ContactEmail: convert.StrVal(v.ContactEmail),
		ContactPhone: convert.StrVal(v.ContactPhone),
		Status:       v.Status,
		CreatedAt:    convert.StrVal(v.CreatedAt),
		UpdatedAt:    convert.StrVal(v.UpdatedAt),
	}
}

func mapToProtoMerchantResponseFromModel(v *models.Merchant) *pbmerchant.MerchantResponse {
	if v == nil {
		return nil
	}
	return &pbmerchant.MerchantResponse{
		Id:           v.MerchantID,
		UserId:       v.UserID,
		Name:         v.Name,
		Description:  convert.StrVal(v.Description),
		Address:      convert.StrVal(v.Address),
		ContactEmail: convert.StrVal(v.ContactEmail),
		ContactPhone: convert.StrVal(v.ContactPhone),
		Status:       v.Status,
		CreatedAt:    convert.FormatTimePtr(v.CreatedAt),
		UpdatedAt:    convert.FormatTimePtr(v.UpdatedAt),
	}
}

func mapToProtoMerchantResponse(m interface{}) *pbmerchant.MerchantResponse {
	switch v := m.(type) {
	case *repository.MerchantResult:
		return mapToProtoMerchantResponseFromResult(v)
	case *models.Merchant:
		return mapToProtoMerchantResponseFromModel(v)
	default:
		return nil
	}
}

func mapToProtoMerchantResponseDeleteAtFromResult(v *repository.MerchantResult) *pbmerchant.MerchantResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pbmerchant.MerchantResponseDeleteAt{
		Id:           v.MerchantID,
		UserId:       v.UserID,
		Name:         v.Name,
		Description:  convert.StrVal(v.Description),
		Address:      convert.StrVal(v.Address),
		ContactEmail: convert.StrVal(v.ContactEmail),
		ContactPhone: convert.StrVal(v.ContactPhone),
		Status:       v.Status,
		CreatedAt:    convert.StrVal(v.CreatedAt),
		UpdatedAt:    convert.StrVal(v.UpdatedAt),
	}
	res.DeletedAt = convert.StrValToWrappers(v.DeletedAt)
	return res
}

func mapToProtoMerchantResponseDeleteAtFromModel(v *models.Merchant) *pbmerchant.MerchantResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pbmerchant.MerchantResponseDeleteAt{
		Id:           v.MerchantID,
		UserId:       v.UserID,
		Name:         v.Name,
		Description:  convert.StrVal(v.Description),
		Address:      convert.StrVal(v.Address),
		ContactEmail: convert.StrVal(v.ContactEmail),
		ContactPhone: convert.StrVal(v.ContactPhone),
		Status:       v.Status,
		CreatedAt:    convert.FormatTimePtr(v.CreatedAt),
		UpdatedAt:    convert.FormatTimePtr(v.UpdatedAt),
	}
	res.DeletedAt = convert.TimeToWrappers(v.DeletedAt)
	return res
}

func mapToProtoMerchantResponseDeleteAt(m interface{}) *pbmerchant.MerchantResponseDeleteAt {
	switch v := m.(type) {
	case *repository.MerchantResult:
		return mapToProtoMerchantResponseDeleteAtFromResult(v)
	case *models.Merchant:
		return mapToProtoMerchantResponseDeleteAtFromModel(v)
	default:
		return nil
	}
}

func mapToProtoMerchantResponseTrashed(v *repository.MerchantResult) *pbmerchant.MerchantResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pbmerchant.MerchantResponseDeleteAt{
		Id:           v.MerchantID,
		UserId:       v.UserID,
		Name:         v.Name,
		Description:  convert.StrVal(v.Description),
		Address:      convert.StrVal(v.Address),
		ContactEmail: convert.StrVal(v.ContactEmail),
		ContactPhone: convert.StrVal(v.ContactPhone),
		Status:       v.Status,
		CreatedAt:    convert.StrVal(v.CreatedAt),
		UpdatedAt:    convert.StrVal(v.UpdatedAt),
	}
	if v.DeletedAt != nil && *v.DeletedAt != "" {
		res.DeletedAt = convert.StrValToWrappers(v.DeletedAt)
	}
	return res
}

func mapToProtoMerchantDocumentResponseFromResult(v *repository.MerchantDocumentResult) *pbmerchant_document.MerchantDocument {
	if v == nil {
		return nil
	}
	return &pbmerchant_document.MerchantDocument{
		DocumentId:   v.DocumentID,
		MerchantId:   v.MerchantID,
		DocumentType: v.DocumentType,
		DocumentUrl:  v.DocumentUrl,
		Status:       v.Status,
		Note:         convert.StrVal(v.Note),
		UploadedAt:   convert.StrVal(v.UploadedAt),
		UpdatedAt:    convert.StrVal(v.UpdatedAt),
	}
}

func mapToProtoMerchantDocumentResponseFromModel(v *models.MerchantDocument) *pbmerchant_document.MerchantDocument {
	if v == nil {
		return nil
	}
	return &pbmerchant_document.MerchantDocument{
		DocumentId:   v.DocumentID,
		MerchantId:   v.MerchantID,
		DocumentType: v.DocumentType,
		DocumentUrl:  v.DocumentUrl,
		Status:       v.Status,
		Note:         convert.StrVal(v.Note),
		UploadedAt:   convert.FormatTimePtr(v.UploadedAt),
		UpdatedAt:    convert.FormatTimePtr(v.UpdatedAt),
	}
}

func mapToProtoMerchantDocumentResponse(m interface{}) *pbmerchant_document.MerchantDocument {
	switch v := m.(type) {
	case *repository.MerchantDocumentResult:
		return mapToProtoMerchantDocumentResponseFromResult(v)
	case *models.MerchantDocument:
		return mapToProtoMerchantDocumentResponseFromModel(v)
	default:
		return nil
	}
}

func mapToProtoMerchantDocumentResponseAt(v *repository.MerchantDocumentResult) *pbmerchant_document.MerchantDocumentDeleteAt {
	if v == nil {
		return nil
	}
	res := &pbmerchant_document.MerchantDocumentDeleteAt{
		DocumentId:   v.DocumentID,
		MerchantId:   v.MerchantID,
		DocumentType: v.DocumentType,
		DocumentUrl:  v.DocumentUrl,
		Status:       v.Status,
		Note:         convert.StrVal(v.Note),
		UploadedAt:   convert.StrVal(v.UploadedAt),
		UpdatedAt:    convert.StrVal(v.UpdatedAt),
	}
	if v.DeletedAt != nil && *v.DeletedAt != "" {
		res.DeletedAt = convert.StrValToWrappers(v.DeletedAt)
	}
	return res
}
