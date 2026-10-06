package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/microservice-ecommerce-grpc-category/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

const (
	categoryAllCacheKey     = "category:all:page:%d:pageSize:%d:search:%s"
	categoryByIdCacheKey    = "category:id:%d"
	categoryActiveCacheKey  = "category:active:page:%d:pageSize:%d:search:%s"
	categoryTrashedCacheKey = "category:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type categoryCacheResponseDB struct {
	Data         []*repository.CategoryResult `json:"data"`
	TotalRecords *int                         `json:"totalRecords"`
}

type categoryQueryCache struct {
	store *cache.CacheStore
}

func NewCategoryQueryCache(store *cache.CacheStore) *categoryQueryCache {
	return &categoryQueryCache{store: store}
}

func (s *categoryQueryCache) GetCachedCategoriesCache(ctx context.Context, req *requests.FindAllCategory) ([]*repository.CategoryResult, *int, bool) {
	key := fmt.Sprintf(categoryAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[categoryCacheResponseDB](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (s *categoryQueryCache) SetCachedCategoriesCache(ctx context.Context, req *requests.FindAllCategory, data []*repository.CategoryResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.CategoryResult{}
	}
	key := fmt.Sprintf(categoryAllCacheKey, req.Page, req.PageSize, req.Search)
	cache.SetToCache(ctx, s.store, key, &categoryCacheResponseDB{Data: data, TotalRecords: total}, ttlDefault)
}

func (s *categoryQueryCache) GetCachedCategoryActiveCache(ctx context.Context, req *requests.FindAllCategory) ([]*repository.CategoryResult, *int, bool) {
	key := fmt.Sprintf(categoryActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[categoryCacheResponseDB](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (s *categoryQueryCache) SetCachedCategoryActiveCache(ctx context.Context, req *requests.FindAllCategory, data []*repository.CategoryResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.CategoryResult{}
	}
	key := fmt.Sprintf(categoryActiveCacheKey, req.Page, req.PageSize, req.Search)
	cache.SetToCache(ctx, s.store, key, &categoryCacheResponseDB{Data: data, TotalRecords: total}, ttlDefault)
}

func (s *categoryQueryCache) GetCachedCategoryTrashedCache(ctx context.Context, req *requests.FindAllCategory) ([]*repository.CategoryResult, *int, bool) {
	key := fmt.Sprintf(categoryTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[categoryCacheResponseDB](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (s *categoryQueryCache) SetCachedCategoryTrashedCache(ctx context.Context, req *requests.FindAllCategory, data []*repository.CategoryResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.CategoryResult{}
	}
	key := fmt.Sprintf(categoryTrashedCacheKey, req.Page, req.PageSize, req.Search)
	cache.SetToCache(ctx, s.store, key, &categoryCacheResponseDB{Data: data, TotalRecords: total}, ttlDefault)
}

func (s *categoryQueryCache) GetCachedCategoryCache(ctx context.Context, id int) (*models.Category, bool) {
	key := fmt.Sprintf(categoryByIdCacheKey, id)
	result, found := cache.GetFromCache[*models.Category](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *categoryQueryCache) SetCachedCategoryCache(ctx context.Context, data *models.Category) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryByIdCacheKey, data.CategoryID)
	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}
