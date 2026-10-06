package handler

import (
	"encoding/json"
	pbcategory "github.com/MamangRust/microservice-ecommerce-grpc-pb/category"
	"time"

	"github.com/MamangRust/microservice-ecommerce-grpc-category/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func formatTimePtr(t *time.Time) string {
	if t != nil {
		return t.Format("2006-01-02")
	}
	return ""
}

func (h *Handler) mapToCategoryResponseFromModel(v *models.Category) *pbcategory.CategoryResponse {
	if v == nil {
		return nil
	}
	return &pbcategory.CategoryResponse{
		Id:            v.CategoryID,
		Name:          v.Name,
		Description:   ptrStr(v.Description),
		SlugCategory:  ptrStr(v.SlugCategory),
		ImageCategory: ptrStr(v.ImageCategory),
		CreatedAt:     formatTimePtr(v.CreatedAt),
		UpdatedAt:     formatTimePtr(v.UpdatedAt),
	}
}

func (h *Handler) mapToCategoryResponseFromResult(v *repository.CategoryResult) *pbcategory.CategoryResponse {
	if v == nil {
		return nil
	}
	return &pbcategory.CategoryResponse{
		Id:            v.CategoryID,
		Name:          v.Name,
		Description:   ptrStr(v.Description),
		SlugCategory:  ptrStr(v.SlugCategory),
		ImageCategory: ptrStr(v.ImageCategory),
		CreatedAt:     formatTimePtr(v.CreatedAt),
		UpdatedAt:     formatTimePtr(v.UpdatedAt),
	}
}

func (h *Handler) mapToDeleteAtResponseFromModel(v *models.Category) *pbcategory.CategoryResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pbcategory.CategoryResponseDeleteAt{
		Id:            v.CategoryID,
		Name:          v.Name,
		Description:   ptrStr(v.Description),
		SlugCategory:  ptrStr(v.SlugCategory),
		ImageCategory: ptrStr(v.ImageCategory),
		CreatedAt:     formatTimePtr(v.CreatedAt),
		UpdatedAt:     formatTimePtr(v.UpdatedAt),
	}
	if v.DeletedAt != nil {
		res.DeletedAt = &wrapperspb.StringValue{Value: formatTimePtr(v.DeletedAt)}
	}
	return res
}

func (h *Handler) mapToDeleteAtResponseFromResult(v *repository.CategoryResult) *pbcategory.CategoryResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pbcategory.CategoryResponseDeleteAt{
		Id:            v.CategoryID,
		Name:          v.Name,
		Description:   ptrStr(v.Description),
		SlugCategory:  ptrStr(v.SlugCategory),
		ImageCategory: ptrStr(v.ImageCategory),
		CreatedAt:     formatTimePtr(v.CreatedAt),
		UpdatedAt:     formatTimePtr(v.UpdatedAt),
	}
	if v.DeletedAt != nil {
		res.DeletedAt = &wrapperspb.StringValue{Value: formatTimePtr(v.DeletedAt)}
	}
	return res
}

func (h *Handler) mapToPayload(data interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

func ptrStr(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
