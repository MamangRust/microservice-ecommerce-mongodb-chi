package cache

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-slider/repository"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type SliderQueryCache interface {
	GetSliderAllCache(ctx context.Context, req *requests.FindAllSlider) ([]*repository.SliderResult, *int, bool)
	SetSliderAllCache(ctx context.Context, req *requests.FindAllSlider, data []*repository.SliderResult, total *int)

	GetSliderActiveCache(ctx context.Context, req *requests.FindAllSlider) ([]*repository.SliderResult, *int, bool)
	SetSliderActiveCache(ctx context.Context, req *requests.FindAllSlider, data []*repository.SliderResult, total *int)

	GetSliderTrashedCache(ctx context.Context, req *requests.FindAllSlider) ([]*repository.SliderResult, *int, bool)
	SetSliderTrashedCache(ctx context.Context, req *requests.FindAllSlider, data []*repository.SliderResult, total *int)

	GetSliderCache(ctx context.Context, slider_id int) (*repository.SliderResult, bool)
	SetSliderCache(ctx context.Context, data *repository.SliderResult)
}

type SliderCommandCache interface {
	DeleteSliderCache(ctx context.Context, slider_id int)
	InvalidateSliderCache(ctx context.Context)
}
