package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/microservice-ecommerce-grpc-cart/repository"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"go.uber.org/zap"
)

const (
	cartSvcAllCacheKey = "cart:svc:all:user:%d:page:%d:pageSize:%d:search:%s"
	ttlDefault         = 5 * time.Minute
)

type cartCacheResponse struct {
	Data  []*repository.CartResult `json:"data"`
	Total *int                     `json:"total_records"`
}

type cartQueryCache struct {
	store *cache.CacheStore
}

func NewCartQueryCache(store *cache.CacheStore) *cartQueryCache {
	return &cartQueryCache{store: store}
}

func (c *cartQueryCache) GetCachedCartsCache(ctx context.Context, request *requests.FindAllCarts) ([]*repository.CartResult, *int, bool) {
	key := fmt.Sprintf(cartSvcAllCacheKey, request.UserID, request.Page, request.PageSize, request.Search)

	result, found := cache.GetFromCache[cartCacheResponse](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (c *cartQueryCache) SetCartsCache(ctx context.Context, request *requests.FindAllCarts, response []*repository.CartResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	key := fmt.Sprintf(cartSvcAllCacheKey, request.UserID, request.Page, request.PageSize, request.Search)
	payload := &cartCacheResponse{Data: response, Total: total}
	cache.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

func (c *cartQueryCache) DeleteCartsCache(ctx context.Context, userID int) {
	pattern := fmt.Sprintf("cart:svc:all:user:%d:*", userID)
	if _, err := c.store.InvalidateCache(ctx, pattern); err != nil {
		c.store.Logger.Error("Failed to invalidate cart service cache", zap.Error(err))
	}
}
