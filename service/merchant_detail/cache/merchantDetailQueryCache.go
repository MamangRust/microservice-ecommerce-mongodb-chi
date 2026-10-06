package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/microservice-ecommerce-grpc-merchant_detail/repository"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

const (
	merchantDetailAllCacheKey     = "merchant_detail:all:page:%d:pageSize:%d:search:%s"
	merchantDetailByIdCacheKey    = "merchant_detail:id:%d"
	merchantDetailActiveCacheKey  = "merchant_detail:active:page:%d:pageSize:%d:search:%s"
	merchantDetailTrashedCacheKey = "merchant_detail:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type merchantDetailCacheResponseDB struct {
	Data         []*repository.MerchantDetailResult `json:"data"`
	TotalRecords *int                               `json:"total_records"`
}

type merchantDetailQueryCache struct {
	store *cache.CacheStore
}

func NewMerchantDetailQueryCache(store *cache.CacheStore) *merchantDetailQueryCache {
	return &merchantDetailQueryCache{store: store}
}

func (m *merchantDetailQueryCache) GetCachedMerchantDetailsAll(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantDetailResult, *int, bool) {
	key := fmt.Sprintf(merchantDetailAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[merchantDetailCacheResponseDB](ctx, m.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (m *merchantDetailQueryCache) SetCachedMerchantDetailsAll(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantDetailResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.MerchantDetailResult{}
	}
	key := fmt.Sprintf(merchantDetailAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &merchantDetailCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantDetailQueryCache) GetCachedMerchantDetailsActive(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantDetailResult, *int, bool) {
	key := fmt.Sprintf(merchantDetailActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[merchantDetailCacheResponseDB](ctx, m.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (m *merchantDetailQueryCache) SetCachedMerchantDetailsActive(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantDetailResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.MerchantDetailResult{}
	}
	key := fmt.Sprintf(merchantDetailActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &merchantDetailCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantDetailQueryCache) GetCachedMerchantDetailsTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantDetailResult, *int, bool) {
	key := fmt.Sprintf(merchantDetailTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[merchantDetailCacheResponseDB](ctx, m.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (m *merchantDetailQueryCache) SetCachedMerchantDetailsTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantDetailResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.MerchantDetailResult{}
	}
	key := fmt.Sprintf(merchantDetailTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &merchantDetailCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantDetailQueryCache) GetCachedMerchantDetail(ctx context.Context, id int) (*repository.MerchantDetailResult, bool) {
	key := fmt.Sprintf(merchantDetailByIdCacheKey, id)
	result, found := cache.GetFromCache[repository.MerchantDetailResult](ctx, m.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (m *merchantDetailQueryCache) SetCachedMerchantDetail(ctx context.Context, data *repository.MerchantDetailResult) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(merchantDetailByIdCacheKey, data.MerchantDetailID)
	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}
