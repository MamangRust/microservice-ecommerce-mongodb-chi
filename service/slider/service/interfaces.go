package service

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-slider/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type SliderQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllSlider) ([]*repository.SliderResult, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllSlider) ([]*repository.SliderResult, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllSlider) ([]*repository.SliderResult, *int, error)
	FindByID(ctx context.Context, slider_id int) (*repository.SliderResult, error)
}

type SliderCommandService interface {
	Create(ctx context.Context, request *requests.CreateSliderRequest) (*models.Slider, error)
	Update(ctx context.Context, request *requests.UpdateSliderRequest) (*models.Slider, error)
	Trash(ctx context.Context, slider_id int) (*models.Slider, error)
	Restore(ctx context.Context, slider_id int) (*models.Slider, error)
	DeletePermanent(ctx context.Context, slider_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
