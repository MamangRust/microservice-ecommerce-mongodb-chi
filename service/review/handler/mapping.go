package handler

import (
	"fmt"
	pbreview "github.com/MamangRust/microservice-ecommerce-grpc-pb/review"
	"time"

	"github.com/MamangRust/microservice-ecommerce-grpc-review/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/convert"
)

// formatTimeShort formats *time.Time to "2006-01-02" (short date format)
func formatTimeShort(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}

func (h *reviewHandleGrpc) mapResponse(data interface{}) interface{} {
	switch v := data.(type) {
	case *models.Review:
		return h.mapReview(v)
	case *repository.ReviewResult:
		return h.mapReviewResult(v)
	case []*repository.ReviewResult:
		if len(v) > 0 && v[0].ReviewDetails != nil {
			res := make([]*pbreview.ReviewsDetailResponse, len(v))
			for i, r := range v {
				res[i] = h.mapReviewResultDetail(r)
			}
			return res
		}
		if len(v) > 0 && v[0].DeletedAt != nil && *v[0].DeletedAt != "" {
			res := make([]*pbreview.ReviewResponseDeleteAt, len(v))
			for i, r := range v {
				res[i] = h.mapReviewResultDeleteAt(r)
			}
			return res
		}
		res := make([]*pbreview.ReviewResponse, len(v))
		for i, r := range v {
			res[i] = h.mapReviewResultResponse(r)
		}
		return res
	default:
		_ = fmt.Sprintf("%T", data)
		return nil
	}
}

func (h *reviewHandleGrpc) mapReview(v *models.Review) *pbreview.ReviewResponseDeleteAt {
	res := &pbreview.ReviewResponseDeleteAt{
		Id:        v.ReviewID,
		UserId:    v.UserID,
		ProductId: v.ProductID,
		Name:      v.Name,
		Comment:   v.Comment,
		Rating:    v.Rating,
		CreatedAt: convert.FormatTimePtr(v.CreatedAt),
		UpdatedAt: convert.FormatTimePtr(v.UpdatedAt),
	}
	res.DeletedAt = convert.TimeToWrappers(v.DeletedAt)
	return res
}

func (h *reviewHandleGrpc) mapReviewResult(v *repository.ReviewResult) *pbreview.ReviewResponse {
	return &pbreview.ReviewResponse{
		Id:        v.ReviewID,
		UserId:    v.UserID,
		ProductId: v.ProductID,
		Name:      v.Name,
		Comment:   v.Comment,
		Rating:    v.Rating,
		CreatedAt: convert.StrVal(v.CreatedAt),
		UpdatedAt: convert.StrVal(v.UpdatedAt),
	}
}

func (h *reviewHandleGrpc) mapReviewResultResponse(v *repository.ReviewResult) *pbreview.ReviewResponse {
	return &pbreview.ReviewResponse{
		Id:        v.ReviewID,
		UserId:    v.UserID,
		ProductId: v.ProductID,
		Name:      v.Name,
		Comment:   v.Comment,
		Rating:    v.Rating,
		CreatedAt: convert.StrVal(v.CreatedAt),
		UpdatedAt: convert.StrVal(v.UpdatedAt),
	}
}

func (h *reviewHandleGrpc) mapReviewResultDeleteAt(v *repository.ReviewResult) *pbreview.ReviewResponseDeleteAt {
	res := &pbreview.ReviewResponseDeleteAt{
		Id:        v.ReviewID,
		UserId:    v.UserID,
		ProductId: v.ProductID,
		Name:      v.Name,
		Comment:   v.Comment,
		Rating:    v.Rating,
		CreatedAt: convert.StrVal(v.CreatedAt),
		UpdatedAt: convert.StrVal(v.UpdatedAt),
	}
	res.DeletedAt = convert.StrValToWrappers(v.DeletedAt)
	return res
}

func (h *reviewHandleGrpc) mapReviewResultDetail(v *repository.ReviewResult) *pbreview.ReviewsDetailResponse {
	res := &pbreview.ReviewsDetailResponse{
		Id:        v.ReviewID,
		UserId:    v.UserID,
		ProductId: v.ProductID,
		Name:      v.Name,
		Comment:   v.Comment,
		Rating:    v.Rating,
		CreatedAt: convert.StrVal(v.CreatedAt),
		UpdatedAt: convert.StrVal(v.UpdatedAt),
	}
	if v.DeletedAt != nil && *v.DeletedAt != "" {
		res.DeletedAt = *v.DeletedAt
	}
	return res
}
