package handler

import (
	pbslider "github.com/MamangRust/microservice-ecommerce-grpc-pb/slider"

	"github.com/MamangRust/microservice-ecommerce-grpc-slider/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/convert"
)

func mapToSliderResponse(slider interface{}) *pbslider.SliderResponse {
	switch v := slider.(type) {
	case *models.Slider:
		return &pbslider.SliderResponse{Id: v.SliderID, Name: v.Name, Image: v.Image, CreatedAt: convert.FormatTimePtr(v.CreatedAt), UpdatedAt: convert.FormatTimePtr(v.UpdatedAt)}
	case *repository.SliderResult:
		return &pbslider.SliderResponse{Id: v.SliderID, Name: v.Name, Image: v.Image, CreatedAt: convert.StrVal(v.CreatedAt), UpdatedAt: convert.StrVal(v.UpdatedAt)}
	default:
		return nil
	}
}

func mapToSliderResponseDeleteAt(slider interface{}) *pbslider.SliderResponseDeleteAt {
	switch v := slider.(type) {
	case *models.Slider:
		res := &pbslider.SliderResponseDeleteAt{Id: v.SliderID, Name: v.Name, Image: v.Image, CreatedAt: convert.FormatTimePtr(v.CreatedAt), UpdatedAt: convert.FormatTimePtr(v.UpdatedAt)}
		res.DeletedAt = convert.TimeToWrappers(v.DeletedAt)
		return res
	case *repository.SliderResult:
		res := &pbslider.SliderResponseDeleteAt{Id: v.SliderID, Name: v.Name, Image: v.Image, CreatedAt: convert.StrVal(v.CreatedAt), UpdatedAt: convert.StrVal(v.UpdatedAt)}
		res.DeletedAt = convert.StrValToWrappers(v.DeletedAt)
		return res
	default:
		return nil
	}
}
