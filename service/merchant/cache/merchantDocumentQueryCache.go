package cache

import (
	"context"
	"fmt"

	"github.com/MamangRust/microservice-ecommerce-grpc-merchant/repository"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

const (
	merchantDocumentAllCacheKey     = "merchant_document:all:page:%d:pageSize:%d:search:%s"
	merchantDocumentByIdCacheKey    = "merchant_document:id:%d"
	merchantDocumentActiveCacheKey  = "merchant_document:active:page:%d:pageSize:%d:search:%s"
	merchantDocumentTrashedCacheKey = "merchant_document:trashed:page:%d:pageSize:%d:search:%s"
)

type merchantDocumentCachedResponseDB struct {
	Data         []*repository.MerchantDocumentResult `json:"data"`
	TotalRecords *int                                 `json:"total_records"`
}

type merchantDocumentQueryCache struct {
	store *cache.CacheStore
}

func NewMerchantDocumentQueryCache(store *cache.CacheStore) *merchantDocumentQueryCache {
	return &merchantDocumentQueryCache{store: store}
}

func (s *merchantDocumentQueryCache) GetCachedMerchantDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*repository.MerchantDocumentResult, *int, bool) {
	key := fmt.Sprintf(merchantDocumentAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[merchantDocumentCachedResponseDB](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (s *merchantDocumentQueryCache) SetCachedMerchantDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*repository.MerchantDocumentResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.MerchantDocumentResult{}
	}
	key := fmt.Sprintf(merchantDocumentAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &merchantDocumentCachedResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *merchantDocumentQueryCache) GetCachedMerchantDocumentsActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*repository.MerchantDocumentResult, *int, bool) {
	key := fmt.Sprintf(merchantDocumentActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[merchantDocumentCachedResponseDB](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (s *merchantDocumentQueryCache) SetCachedMerchantDocumentsActive(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*repository.MerchantDocumentResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.MerchantDocumentResult{}
	}
	key := fmt.Sprintf(merchantDocumentActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &merchantDocumentCachedResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *merchantDocumentQueryCache) GetCachedMerchantDocumentsTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*repository.MerchantDocumentResult, *int, bool) {
	key := fmt.Sprintf(merchantDocumentTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[merchantDocumentCachedResponseDB](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (s *merchantDocumentQueryCache) SetCachedMerchantDocumentsTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*repository.MerchantDocumentResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.MerchantDocumentResult{}
	}
	key := fmt.Sprintf(merchantDocumentTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &merchantDocumentCachedResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *merchantDocumentQueryCache) GetCachedMerchantDocument(ctx context.Context, id int) (*repository.MerchantDocumentResult, bool) {
	key := fmt.Sprintf(merchantDocumentByIdCacheKey, id)
	result, found := cache.GetFromCache[repository.MerchantDocumentResult](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (s *merchantDocumentQueryCache) SetCachedMerchantDocument(ctx context.Context, data *repository.MerchantDocumentResult) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(merchantDocumentByIdCacheKey, data.DocumentID)
	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}
