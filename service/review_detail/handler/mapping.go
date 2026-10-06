package handler

import (
	pbcommon "github.com/MamangRust/microservice-ecommerce-grpc-pb/common"
	pbreview_detail "github.com/MamangRust/microservice-ecommerce-grpc-pb/review_detail"
	"math"

	"github.com/MamangRust/microservice-ecommerce-grpc-review-detail/repository"
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

func mapToReviewDetailResponseFromResult(v *repository.ReviewDetailResult) *pbreview_detail.ReviewDetailsResponse {
	if v == nil {
		return nil
	}
	return &pbreview_detail.ReviewDetailsResponse{
		Id:        v.ReviewDetailID,
		ReviewId:  v.ReviewID,
		Type:      v.Type,
		Url:       v.Url,
		Caption:   convert.StrVal(v.Caption),
		CreatedAt: convert.StrVal(v.CreatedAt),
		UpdatedAt: convert.StrVal(v.UpdatedAt),
	}
}

func mapToReviewDetailResponseFromModel(v *models.ReviewDetail) *pbreview_detail.ReviewDetailsResponse {
	if v == nil {
		return nil
	}
	return &pbreview_detail.ReviewDetailsResponse{
		Id:        v.ReviewDetailID,
		ReviewId:  v.ReviewID,
		Type:      v.Type,
		Url:       v.Url,
		Caption:   convert.StrVal(v.Caption),
		CreatedAt: convert.FormatTimePtr(v.CreatedAt),
		UpdatedAt: convert.FormatTimePtr(v.UpdatedAt),
	}
}

func mapToReviewDetailResponseDeleteAtFromResult(v *repository.ReviewDetailResult) *pbreview_detail.ReviewDetailsResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pbreview_detail.ReviewDetailsResponseDeleteAt{
		Id:        v.ReviewDetailID,
		ReviewId:  v.ReviewID,
		Type:      v.Type,
		Url:       v.Url,
		Caption:   convert.StrVal(v.Caption),
		CreatedAt: convert.StrVal(v.CreatedAt),
		UpdatedAt: convert.StrVal(v.UpdatedAt),
	}
	if v.DeletedAt != nil && *v.DeletedAt != "" {
		res.DeletedAt = &wrapperspb.StringValue{Value: *v.DeletedAt}
	}
	return res
}

func mapToReviewDetailResponseDeleteAtFromModel(v *models.ReviewDetail) *pbreview_detail.ReviewDetailsResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pbreview_detail.ReviewDetailsResponseDeleteAt{
		Id:        v.ReviewDetailID,
		ReviewId:  v.ReviewID,
		Type:      v.Type,
		Url:       v.Url,
		Caption:   convert.StrVal(v.Caption),
		CreatedAt: convert.FormatTimePtr(v.CreatedAt),
		UpdatedAt: convert.FormatTimePtr(v.UpdatedAt),
	}
	if v.DeletedAt != nil {
		res.DeletedAt = &wrapperspb.StringValue{Value: convert.FormatTimePtr(v.DeletedAt)}
	}
	return res
}
