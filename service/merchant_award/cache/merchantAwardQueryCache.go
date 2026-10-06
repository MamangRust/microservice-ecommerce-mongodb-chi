package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/microservice-ecommerce-grpc-merchant_award/repository"
	sharedcache "github.com/MamangRust/microservice-ecommerce-shared/cache"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

const (
	merchantAwardAllCacheKey     = "merchant_award:all:page:%d:pageSize:%d:search:%s"
	merchantAwardByIdCacheKey    = "merchant_award:id:%d"
	merchantAwardActiveCacheKey  = "merchant_award:active:page:%d:pageSize:%d:search:%s"
	merchantAwardTrashedCacheKey = "merchant_award:trashed:page:%d:pageSize:%d:search:%s"
	ttlDefault                   = 5 * time.Minute
)

type merchantAwardCacheResponseDB struct {
	Data         []*repository.MerchantCertResult `json:"data"`
	TotalRecords *int                             `json:"total_records"`
}

type merchantAwardQueryCache struct {
	store *sharedcache.CacheStore
}

func NewMerchantAwardQueryCache(store *sharedcache.CacheStore) *merchantAwardQueryCache {
	return &merchantAwardQueryCache{store: store}
}

func (m *merchantAwardQueryCache) GetCachedMerchantAwardAll(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantCertResult, *int, bool) {
	key := fmt.Sprintf(merchantAwardAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := sharedcache.GetFromCache[merchantAwardCacheResponseDB](ctx, m.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (m *merchantAwardQueryCache) SetCachedMerchantAwardAll(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantCertResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.MerchantCertResult{}
	}
	key := fmt.Sprintf(merchantAwardAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &merchantAwardCacheResponseDB{Data: data, TotalRecords: total}
	sharedcache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantAwardQueryCache) GetCachedMerchantAwardActive(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantCertResult, *int, bool) {
	key := fmt.Sprintf(merchantAwardActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := sharedcache.GetFromCache[merchantAwardCacheResponseDB](ctx, m.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (m *merchantAwardQueryCache) SetCachedMerchantAwardActive(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantCertResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.MerchantCertResult{}
	}
	key := fmt.Sprintf(merchantAwardActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &merchantAwardCacheResponseDB{Data: data, TotalRecords: total}
	sharedcache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantAwardQueryCache) GetCachedMerchantAwardTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantCertResult, *int, bool) {
	key := fmt.Sprintf(merchantAwardTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := sharedcache.GetFromCache[merchantAwardCacheResponseDB](ctx, m.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (m *merchantAwardQueryCache) SetCachedMerchantAwardTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantCertResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.MerchantCertResult{}
	}
	key := fmt.Sprintf(merchantAwardTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &merchantAwardCacheResponseDB{Data: data, TotalRecords: total}
	sharedcache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantAwardQueryCache) GetCachedMerchantAward(ctx context.Context, id int) (*repository.MerchantCertResult, bool) {
	key := fmt.Sprintf(merchantAwardByIdCacheKey, id)
	result, found := sharedcache.GetFromCache[repository.MerchantCertResult](ctx, m.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (m *merchantAwardQueryCache) SetCachedMerchantAward(ctx context.Context, data *repository.MerchantCertResult) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(merchantAwardByIdCacheKey, data.MerchantCertificationID)
	sharedcache.SetToCache(ctx, m.store, key, data, ttlDefault)
}
