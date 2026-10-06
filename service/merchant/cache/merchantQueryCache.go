package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/microservice-ecommerce-grpc-merchant/repository"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

const (
	merchantAllCacheKey     = "merchant:all:page:%d:pageSize:%d:search:%s"
	merchantByIdCacheKey    = "merchant:id:%d"
	merchantActiveCacheKey  = "merchant:active:page:%d:pageSize:%d:search:%s"
	merchantTrashedCacheKey = "merchant:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type merchantCachedResponseDB struct {
	Data         []*repository.MerchantResult `json:"data"`
	TotalRecords *int                         `json:"total_records"`
}

type merchantQueryCache struct {
	store *cache.CacheStore
}

func NewMerchantQueryCache(store *cache.CacheStore) *merchantQueryCache {
	return &merchantQueryCache{store: store}
}

func (m *merchantQueryCache) GetMerchantAllCache(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantResult, *int, bool) {
	key := fmt.Sprintf(merchantAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[merchantCachedResponseDB](ctx, m.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (m *merchantQueryCache) SetMerchantAllCache(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.MerchantResult{}
	}
	key := fmt.Sprintf(merchantAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &merchantCachedResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantQueryCache) GetMerchantActiveCache(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantResult, *int, bool) {
	key := fmt.Sprintf(merchantActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[merchantCachedResponseDB](ctx, m.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (m *merchantQueryCache) SetMerchantActiveCache(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.MerchantResult{}
	}
	key := fmt.Sprintf(merchantActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &merchantCachedResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantQueryCache) GetMerchantTrashedCache(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantResult, *int, bool) {
	key := fmt.Sprintf(merchantTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[merchantCachedResponseDB](ctx, m.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (m *merchantQueryCache) SetMerchantTrashedCache(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.MerchantResult{}
	}
	key := fmt.Sprintf(merchantTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &merchantCachedResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantQueryCache) GetCachedMerchantCache(ctx context.Context, id int) (*repository.MerchantResult, bool) {
	key := fmt.Sprintf(merchantByIdCacheKey, id)
	result, found := cache.GetFromCache[repository.MerchantResult](ctx, m.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (m *merchantQueryCache) SetCachedMerchantCache(ctx context.Context, data *repository.MerchantResult) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(merchantByIdCacheKey, data.MerchantID)
	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}

func (m *merchantQueryCache) InvalidateMerchantCache(ctx context.Context) {
	// no-op for now
}
