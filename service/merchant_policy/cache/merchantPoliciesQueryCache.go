package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/microservice-ecommerce-grpc-merchant_policy/repository"
	sharedcache "github.com/MamangRust/microservice-ecommerce-shared/cache"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

const (
	merchantPolicyAllCacheKey     = "merchant_policy:all:page:%d:pageSize:%d:search:%s"
	merchantPolicyByIdCacheKey    = "merchant_policy:id:%d"
	merchantPolicyActiveCacheKey  = "merchant_policy:active:page:%d:pageSize:%d:search:%s"
	merchantPolicyTrashedCacheKey = "merchant_policy:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type merchantPolicyCacheResponseDB struct {
	Data         []*repository.MerchantPolicyResult `json:"data"`
	TotalRecords *int                               `json:"total_records"`
}

type merchantPolicyActiveCacheResponseDB struct {
	Data         []*repository.MerchantPolicyResult `json:"data"`
	TotalRecords *int                               `json:"total_records"`
}

type merchantPolicyTrashedCacheResponseDB struct {
	Data         []*repository.MerchantPolicyResult `json:"data"`
	TotalRecords *int                               `json:"total_records"`
}

type merchantPoliciesQueryCache struct {
	store *sharedcache.CacheStore
}

func NewMerchantPoliciesQueryCache(store *sharedcache.CacheStore) MerchantPoliciesQueryCache {
	return &merchantPoliciesQueryCache{store: store}
}

func (m *merchantPoliciesQueryCache) GetCachedMerchantPolicyAll(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantPolicyResult, *int, bool) {
	key := fmt.Sprintf(merchantPolicyAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := sharedcache.GetFromCache[merchantPolicyCacheResponseDB](ctx, m.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (m *merchantPoliciesQueryCache) SetCachedMerchantPolicyAll(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantPolicyResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.MerchantPolicyResult{}
	}
	key := fmt.Sprintf(merchantPolicyAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &merchantPolicyCacheResponseDB{Data: data, TotalRecords: total}
	sharedcache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantPoliciesQueryCache) GetCachedMerchantPolicyActive(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantPolicyResult, *int, bool) {
	key := fmt.Sprintf(merchantPolicyActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := sharedcache.GetFromCache[merchantPolicyActiveCacheResponseDB](ctx, m.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (m *merchantPoliciesQueryCache) SetCachedMerchantPolicyActive(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantPolicyResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.MerchantPolicyResult{}
	}
	key := fmt.Sprintf(merchantPolicyActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &merchantPolicyActiveCacheResponseDB{Data: data, TotalRecords: total}
	sharedcache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantPoliciesQueryCache) GetCachedMerchantPolicyTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantPolicyResult, *int, bool) {
	key := fmt.Sprintf(merchantPolicyTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := sharedcache.GetFromCache[merchantPolicyTrashedCacheResponseDB](ctx, m.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (m *merchantPoliciesQueryCache) SetCachedMerchantPolicyTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantPolicyResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.MerchantPolicyResult{}
	}
	key := fmt.Sprintf(merchantPolicyTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &merchantPolicyTrashedCacheResponseDB{Data: data, TotalRecords: total}
	sharedcache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantPoliciesQueryCache) GetCachedMerchantPolicy(ctx context.Context, id int) (*repository.MerchantPolicyResult, bool) {
	key := fmt.Sprintf(merchantPolicyByIdCacheKey, id)
	result, found := sharedcache.GetFromCache[repository.MerchantPolicyResult](ctx, m.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (m *merchantPoliciesQueryCache) SetCachedMerchantPolicy(ctx context.Context, data *repository.MerchantPolicyResult) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(merchantPolicyByIdCacheKey, data.MerchantPolicyID)
	sharedcache.SetToCache(ctx, m.store, key, data, ttlDefault)
}
