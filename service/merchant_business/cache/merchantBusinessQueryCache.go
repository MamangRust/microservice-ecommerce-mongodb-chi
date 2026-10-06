package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/microservice-ecommerce-grpc-merchant_business/repository"
	sharedcache "github.com/MamangRust/microservice-ecommerce-shared/cache"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

const (
	merchantBusinessAllCacheKey     = "merchant_business:all:page:%d:pageSize:%d:search:%s"
	merchantBusinessByIdCacheKey    = "merchant_business:id:%d"
	merchantBusinessActiveCacheKey  = "merchant_business:active:page:%d:pageSize:%d:search:%s"
	merchantBusinessTrashedCacheKey = "merchant_business:trashed:page:%d:pageSize:%d:search:%s"
	ttlDefault                      = 5 * time.Minute
)

type merchantBusinessCacheResponseDB struct {
	Data         []*repository.MerchantBusinessResult `json:"data"`
	TotalRecords *int                                 `json:"total_records"`
}

type merchantBusinessQueryCache struct {
	store *sharedcache.CacheStore
}

func NewMerchantBusinessQueryCache(store *sharedcache.CacheStore) *merchantBusinessQueryCache {
	return &merchantBusinessQueryCache{store: store}
}

func (m *merchantBusinessQueryCache) GetCachedMerchantBusinessAll(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantBusinessResult, *int, bool) {
	key := fmt.Sprintf(merchantBusinessAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := sharedcache.GetFromCache[merchantBusinessCacheResponseDB](ctx, m.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (m *merchantBusinessQueryCache) SetCachedMerchantBusinessAll(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantBusinessResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.MerchantBusinessResult{}
	}
	key := fmt.Sprintf(merchantBusinessAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &merchantBusinessCacheResponseDB{Data: data, TotalRecords: total}
	sharedcache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantBusinessQueryCache) GetCachedMerchantBusinessActive(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantBusinessResult, *int, bool) {
	key := fmt.Sprintf(merchantBusinessActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := sharedcache.GetFromCache[merchantBusinessCacheResponseDB](ctx, m.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (m *merchantBusinessQueryCache) SetCachedMerchantBusinessActive(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantBusinessResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.MerchantBusinessResult{}
	}
	key := fmt.Sprintf(merchantBusinessActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &merchantBusinessCacheResponseDB{Data: data, TotalRecords: total}
	sharedcache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantBusinessQueryCache) GetCachedMerchantBusinessTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantBusinessResult, *int, bool) {
	key := fmt.Sprintf(merchantBusinessTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := sharedcache.GetFromCache[merchantBusinessCacheResponseDB](ctx, m.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (m *merchantBusinessQueryCache) SetCachedMerchantBusinessTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantBusinessResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.MerchantBusinessResult{}
	}
	key := fmt.Sprintf(merchantBusinessTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &merchantBusinessCacheResponseDB{Data: data, TotalRecords: total}
	sharedcache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantBusinessQueryCache) GetCachedMerchantBusiness(ctx context.Context, id int) (*repository.MerchantBusinessResult, bool) {
	key := fmt.Sprintf(merchantBusinessByIdCacheKey, id)
	result, found := sharedcache.GetFromCache[repository.MerchantBusinessResult](ctx, m.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (m *merchantBusinessQueryCache) SetCachedMerchantBusiness(ctx context.Context, data *repository.MerchantBusinessResult) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(merchantBusinessByIdCacheKey, data.MerchantBusinessInfoID)
	sharedcache.SetToCache(ctx, m.store, key, data, ttlDefault)
}
