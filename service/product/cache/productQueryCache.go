package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/microservice-ecommerce-grpc-product/repository"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

const (
	productAllCacheKey      = "product:all:page:%d:pageSize:%d:search:%s"
	productCategoryCacheKey = "product:category:%s:page:%d:pageSize:%d:search:%s"
	productMerchantCacheKey = "product:merchant:%d:page:%d:pageSize:%d:search:%s"

	productActiveCacheKey  = "product:active:page:%d:pageSize:%d:search:%s"
	productTrashedCacheKey = "product:trashed:page:%d:pageSize:%d:search:%s"
	productByIdCacheKey    = "product:id:%d"

	ttlDefault = 5 * time.Minute
)

type productCacheResponse struct {
	Data         []*repository.ProductResult `json:"data"`
	TotalRecords *int                        `json:"total_records"`
}

type productQueryCache struct {
	store *cache.CacheStore
}

func NewProductQueryCache(store *cache.CacheStore) *productQueryCache {
	return &productQueryCache{store: store}
}

func (p *productQueryCache) GetCachedProducts(ctx context.Context, req *requests.FindAllProduct) ([]*repository.ProductResult, *int, bool) {
	key := fmt.Sprintf(productAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[productCacheResponse](ctx, p.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (p *productQueryCache) SetCachedProducts(ctx context.Context, req *requests.FindAllProduct, data []*repository.ProductResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.ProductResult{}
	}
	key := fmt.Sprintf(productAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &productCacheResponse{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, p.store, key, payload, ttlDefault)
}

func (p *productQueryCache) GetCachedProductsByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant) ([]*repository.ProductResult, *int, bool) {
	key := fmt.Sprintf(productMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[productCacheResponse](ctx, p.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (p *productQueryCache) SetCachedProductsByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant, data []*repository.ProductResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.ProductResult{}
	}
	key := fmt.Sprintf(productMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)
	payload := &productCacheResponse{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, p.store, key, payload, ttlDefault)
}

func (p *productQueryCache) GetCachedProductsByCategory(ctx context.Context, req *requests.FindAllProductByCategory) ([]*repository.ProductResult, *int, bool) {
	key := fmt.Sprintf(productCategoryCacheKey, req.CategoryName, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[productCacheResponse](ctx, p.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (p *productQueryCache) SetCachedProductsByCategory(ctx context.Context, req *requests.FindAllProductByCategory, data []*repository.ProductResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.ProductResult{}
	}
	key := fmt.Sprintf(productCategoryCacheKey, req.CategoryName, req.Page, req.PageSize, req.Search)
	payload := &productCacheResponse{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, p.store, key, payload, ttlDefault)
}

func (p *productQueryCache) GetCachedProductActive(ctx context.Context, req *requests.FindAllProduct) ([]*repository.ProductResult, *int, bool) {
	key := fmt.Sprintf(productActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[productCacheResponse](ctx, p.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (p *productQueryCache) SetCachedProductActive(ctx context.Context, req *requests.FindAllProduct, data []*repository.ProductResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.ProductResult{}
	}
	key := fmt.Sprintf(productActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &productCacheResponse{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, p.store, key, payload, ttlDefault)
}

func (p *productQueryCache) GetCachedProductTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*repository.ProductResult, *int, bool) {
	key := fmt.Sprintf(productTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[productCacheResponse](ctx, p.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (p *productQueryCache) SetCachedProductTrashed(ctx context.Context, req *requests.FindAllProduct, data []*repository.ProductResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.ProductResult{}
	}
	key := fmt.Sprintf(productTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &productCacheResponse{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, p.store, key, payload, ttlDefault)
}

func (p *productQueryCache) GetCachedProduct(ctx context.Context, productID int) (*repository.ProductResult, bool) {
	key := fmt.Sprintf(productByIdCacheKey, productID)
	result, found := cache.GetFromCache[repository.ProductResult](ctx, p.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (p *productQueryCache) SetCachedProduct(ctx context.Context, data *repository.ProductResult) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(productByIdCacheKey, data.ProductID)
	cache.SetToCache(ctx, p.store, key, data, ttlDefault)
}
