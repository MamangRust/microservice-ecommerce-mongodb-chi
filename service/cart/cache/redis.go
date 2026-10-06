package cache

import "github.com/MamangRust/microservice-ecommerce-shared/cache"

type cartMencache struct {
	CartQueryCache
}

func NewMencache(cacheStore *cache.CacheStore) CartMencache {
	return &cartMencache{
		CartQueryCache: NewCartQueryCache(cacheStore),
	}
}
