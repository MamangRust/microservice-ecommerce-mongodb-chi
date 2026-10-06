package cache

import (
	"context"
	"fmt"

	"github.com/MamangRust/microservice-ecommerce-grpc-banner/repository"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type bannerCacheResponse struct {
	Data  []*repository.BannerResult `json:"data"`
	Total *int                       `json:"total"`
}

type bannerQueryCache struct {
	store *cache.CacheStore
}

func NewBannerQueryCache(store *cache.CacheStore) *bannerQueryCache {
	return &bannerQueryCache{store: store}
}

func (b *bannerQueryCache) GetCachedBannersCache(ctx context.Context, req *requests.FindAllBanner) ([]*repository.BannerResult, *int, bool) {
	key := fmt.Sprintf(bannerAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[bannerCacheResponse](ctx, b.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.Total, true
}

func (b *bannerQueryCache) SetCachedBannersCache(ctx context.Context, req *requests.FindAllBanner, data []*repository.BannerResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	key := fmt.Sprintf(bannerAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &bannerCacheResponse{Data: data, Total: total}
	cache.SetToCache(ctx, b.store, key, payload, ttlDefault)
}

func (b *bannerQueryCache) GetCachedBannerActiveCache(ctx context.Context, req *requests.FindAllBanner) ([]*repository.BannerResult, *int, bool) {
	key := fmt.Sprintf(bannerActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[bannerCacheResponse](ctx, b.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.Total, true
}

func (b *bannerQueryCache) SetCachedBannerActiveCache(ctx context.Context, req *requests.FindAllBanner, data []*repository.BannerResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.BannerResult{}
	}
	key := fmt.Sprintf(bannerActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &bannerCacheResponse{Data: data, Total: total}
	cache.SetToCache(ctx, b.store, key, payload, ttlDefault)
}

func (b *bannerQueryCache) GetCachedBannerTrashedCache(ctx context.Context, req *requests.FindAllBanner) ([]*repository.BannerResult, *int, bool) {
	key := fmt.Sprintf(bannerTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[bannerCacheResponse](ctx, b.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.Total, true
}

func (b *bannerQueryCache) SetCachedBannerTrashedCache(ctx context.Context, req *requests.FindAllBanner, data []*repository.BannerResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.BannerResult{}
	}
	key := fmt.Sprintf(bannerTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &bannerCacheResponse{Data: data, Total: total}
	cache.SetToCache(ctx, b.store, key, payload, ttlDefault)
}

func (b *bannerQueryCache) GetCachedBannerCache(ctx context.Context, id int) (*repository.BannerResult, bool) {
	key := fmt.Sprintf(bannerByIdCacheKey, id)
	result, found := cache.GetFromCache[repository.BannerResult](ctx, b.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (b *bannerQueryCache) SetCachedBannerCache(ctx context.Context, data *repository.BannerResult) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(bannerByIdCacheKey, data.BannerID)
	cache.SetToCache(ctx, b.store, key, data, ttlDefault)
}
