package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/microservice-ecommerce-grpc-slider/repository"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

const (
	sliderAllCacheKey     = "slider:all:page:%d:pageSize:%d:search:%s"
	sliderActiveCacheKey  = "slider:active:page:%d:pageSize:%d:search:%s"
	sliderTrashedCacheKey = "slider:trashed:page:%d:pageSize:%d:search:%s"
	sliderIdKey           = "slider:id:%d"
	ttlDefault            = 5 * time.Minute
)

type sliderCacheResponse struct {
	Data  []*repository.SliderResult `json:"data"`
	Total *int                       `json:"total_records"`
}

type sliderQueryCache struct {
	store *cache.CacheStore
}

func NewSliderQueryCache(store *cache.CacheStore) *sliderQueryCache {
	return &sliderQueryCache{store: store}
}

func (s *sliderQueryCache) GetSliderAllCache(ctx context.Context, req *requests.FindAllSlider) ([]*repository.SliderResult, *int, bool) {
	key := fmt.Sprintf(sliderAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[sliderCacheResponse](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.Total, true
}

func (s *sliderQueryCache) SetSliderAllCache(ctx context.Context, req *requests.FindAllSlider, data []*repository.SliderResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.SliderResult{}
	}
	key := fmt.Sprintf(sliderAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &sliderCacheResponse{Data: data, Total: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *sliderQueryCache) GetSliderActiveCache(ctx context.Context, req *requests.FindAllSlider) ([]*repository.SliderResult, *int, bool) {
	key := fmt.Sprintf(sliderActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[sliderCacheResponse](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.Total, true
}

func (s *sliderQueryCache) SetSliderActiveCache(ctx context.Context, req *requests.FindAllSlider, data []*repository.SliderResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.SliderResult{}
	}
	key := fmt.Sprintf(sliderActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &sliderCacheResponse{Data: data, Total: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *sliderQueryCache) GetSliderTrashedCache(ctx context.Context, req *requests.FindAllSlider) ([]*repository.SliderResult, *int, bool) {
	key := fmt.Sprintf(sliderTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[sliderCacheResponse](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.Total, true
}

func (s *sliderQueryCache) SetSliderTrashedCache(ctx context.Context, req *requests.FindAllSlider, data []*repository.SliderResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.SliderResult{}
	}
	key := fmt.Sprintf(sliderTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &sliderCacheResponse{Data: data, Total: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *sliderQueryCache) GetSliderCache(ctx context.Context, slider_id int) (*repository.SliderResult, bool) {
	key := fmt.Sprintf(sliderIdKey, slider_id)
	result, found := cache.GetFromCache[repository.SliderResult](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (s *sliderQueryCache) SetSliderCache(ctx context.Context, data *repository.SliderResult) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(sliderIdKey, data.SliderID)
	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}
