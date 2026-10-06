package repository

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type SliderResult struct {
	SliderID   int32
	Name       string
	Image      string
	CreatedAt  *string
	UpdatedAt  *string
	DeletedAt  *string
	TotalCount int64
}

type SliderQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllSlider) ([]*SliderResult, error)
	FindActive(ctx context.Context, req *requests.FindAllSlider) ([]*SliderResult, error)
	FindTrashed(ctx context.Context, req *requests.FindAllSlider) ([]*SliderResult, error)
	FindByID(ctx context.Context, slider_id int) (*SliderResult, error)
}

type SliderCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateSliderRequest) (*models.Slider, error)
	Update(ctx context.Context, request *requests.UpdateSliderRequest) (*models.Slider, error)
	Trash(ctx context.Context, slider_id int) (*models.Slider, error)
	Restore(ctx context.Context, slider_id int) (*models.Slider, error)
	DeletePermanent(ctx context.Context, slider_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
