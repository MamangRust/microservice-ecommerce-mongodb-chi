package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/microservice-ecommerce-grpc-order-item/repository"
	sharedcache "github.com/MamangRust/microservice-ecommerce-shared/cache"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

const (
	orderItemAllCacheKey     = "order_item:all:page:%d:pageSize:%d:search:%s"
	orderItemActiveCacheKey  = "order_item:active:page:%d:pageSize:%d:search:%s"
	orderItemTrashedCacheKey = "order_item:trashed:page:%d:pageSize:%d:search:%s"
	orderItemByOrderCacheKey = "order_item:order:%d"

	ttlDefault = 5 * time.Minute
)

type orderItemCacheResponseDB struct {
	Data         []*repository.OrderItemResult `json:"data"`
	TotalRecords *int                          `json:"total_records"`
}

type orderItemActiveCacheResponseDB struct {
	Data         []*repository.OrderItemResult `json:"data"`
	TotalRecords *int                          `json:"total_records"`
}

type orderItemTrashedCacheResponseDB struct {
	Data         []*repository.OrderItemResult `json:"data"`
	TotalRecords *int                          `json:"total_records"`
}

type orderItemQueryCache struct {
	store *sharedcache.CacheStore
}

func NewOrderItemQueryCache(store *sharedcache.CacheStore) *orderItemQueryCache {
	return &orderItemQueryCache{store: store}
}

func (o *orderItemQueryCache) GetCachedOrderItemsAll(ctx context.Context, req *requests.FindAllOrderItems) ([]*repository.OrderItemResult, *int, bool) {
	key := fmt.Sprintf(orderItemAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := sharedcache.GetFromCache[orderItemCacheResponseDB](ctx, o.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (o *orderItemQueryCache) SetCachedOrderItemsAll(ctx context.Context, req *requests.FindAllOrderItems, data []*repository.OrderItemResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.OrderItemResult{}
	}
	key := fmt.Sprintf(orderItemAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &orderItemCacheResponseDB{Data: data, TotalRecords: total}
	sharedcache.SetToCache(ctx, o.store, key, payload, ttlDefault)
}

func (o *orderItemQueryCache) GetCachedOrderItemActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*repository.OrderItemResult, *int, bool) {
	key := fmt.Sprintf(orderItemActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := sharedcache.GetFromCache[orderItemActiveCacheResponseDB](ctx, o.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (o *orderItemQueryCache) SetCachedOrderItemActive(ctx context.Context, req *requests.FindAllOrderItems, data []*repository.OrderItemResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.OrderItemResult{}
	}
	key := fmt.Sprintf(orderItemActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &orderItemActiveCacheResponseDB{Data: data, TotalRecords: total}
	sharedcache.SetToCache(ctx, o.store, key, payload, ttlDefault)
}

func (o *orderItemQueryCache) GetCachedOrderItemTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*repository.OrderItemResult, *int, bool) {
	key := fmt.Sprintf(orderItemTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := sharedcache.GetFromCache[orderItemTrashedCacheResponseDB](ctx, o.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (o *orderItemQueryCache) SetCachedOrderItemTrashed(ctx context.Context, req *requests.FindAllOrderItems, data []*repository.OrderItemResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.OrderItemResult{}
	}
	key := fmt.Sprintf(orderItemTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &orderItemTrashedCacheResponseDB{Data: data, TotalRecords: total}
	sharedcache.SetToCache(ctx, o.store, key, payload, ttlDefault)
}

func (o *orderItemQueryCache) GetCachedOrderItems(ctx context.Context, orderID int) ([]*repository.OrderItemResult, bool) {
	key := fmt.Sprintf(orderItemByOrderCacheKey, orderID)
	result, found := sharedcache.GetFromCache[[]*repository.OrderItemResult](ctx, o.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (o *orderItemQueryCache) SetCachedOrderItems(ctx context.Context, data []*repository.OrderItemResult) {
	if len(data) == 0 {
		return
	}
	key := fmt.Sprintf(orderItemByOrderCacheKey, data[0].OrderID)
	sharedcache.SetToCache(ctx, o.store, key, &data, ttlDefault)
}
