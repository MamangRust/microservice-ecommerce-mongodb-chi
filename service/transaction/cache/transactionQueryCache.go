package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/microservice-ecommerce-grpc-transaction/repository"
	sharedcache "github.com/MamangRust/microservice-ecommerce-shared/cache"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

const (
	transactionAllCacheKey        = "transaction:all:page:%d:pageSize:%d:search:%s"
	transactionByIdCacheKey       = "transaction:id:%d"
	transactionByMerchantCacheKey = "transaction:merchant:%d:page:%d:pageSize:%d:search:%s"
	transactionActiveCacheKey     = "transaction:active:page:%d:pageSize:%d:search:%s"
	transactionTrashedCacheKey    = "transaction:trashed:page:%d:pageSize:%d:search:%s"
	transactionByOrderCacheKey    = "transaction:order:%d"
	ttlDefault                    = 5 * time.Minute
)

type transactionCacheResponseDB struct {
	Data         []*repository.TransactionResult `json:"data"`
	TotalRecords *int                            `json:"totalRecords"`
}

type transactionMechantCacheResponseDB struct {
	Data         []*repository.TransactionResult `json:"data"`
	TotalRecords *int                            `json:"totalRecords"`
}

type transactionActiveCacheResponseDB struct {
	Data         []*repository.TransactionResult `json:"data"`
	TotalRecords *int                            `json:"totalRecords"`
}

type transactionTrashedCacheResponseDB struct {
	Data         []*repository.TransactionResult `json:"data"`
	TotalRecords *int                            `json:"totalRecords"`
}

type transactionQueryCache struct {
	store *sharedcache.CacheStore
}

func NewTransactionQueryCache(store *sharedcache.CacheStore) *transactionQueryCache {
	return &transactionQueryCache{store: store}
}

func (t *transactionQueryCache) GetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransaction) ([]*repository.TransactionResult, *int, bool) {
	key := fmt.Sprintf(transactionAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := sharedcache.GetFromCache[transactionCacheResponseDB](ctx, t.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (t *transactionQueryCache) SetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransaction, data []*repository.TransactionResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.TransactionResult{}
	}

	key := fmt.Sprintf(transactionAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &transactionCacheResponseDB{Data: data, TotalRecords: total}
	sharedcache.SetToCache(ctx, t.store, key, payload, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant) ([]*repository.TransactionResult, *int, bool) {
	key := fmt.Sprintf(transactionByMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)
	result, found := sharedcache.GetFromCache[transactionMechantCacheResponseDB](ctx, t.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (t *transactionQueryCache) SetCachedTransactionByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant, data []*repository.TransactionResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.TransactionResult{}
	}

	key := fmt.Sprintf(transactionByMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)
	payload := &transactionMechantCacheResponseDB{Data: data, TotalRecords: total}
	sharedcache.SetToCache(ctx, t.store, key, payload, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransaction) ([]*repository.TransactionResult, *int, bool) {
	key := fmt.Sprintf(transactionActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := sharedcache.GetFromCache[transactionActiveCacheResponseDB](ctx, t.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (t *transactionQueryCache) SetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransaction, data []*repository.TransactionResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.TransactionResult{}
	}

	key := fmt.Sprintf(transactionActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &transactionActiveCacheResponseDB{Data: data, TotalRecords: total}
	sharedcache.SetToCache(ctx, t.store, key, payload, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransaction) ([]*repository.TransactionResult, *int, bool) {
	key := fmt.Sprintf(transactionTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := sharedcache.GetFromCache[transactionTrashedCacheResponseDB](ctx, t.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (t *transactionQueryCache) SetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransaction, data []*repository.TransactionResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.TransactionResult{}
	}

	key := fmt.Sprintf(transactionTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &transactionTrashedCacheResponseDB{Data: data, TotalRecords: total}
	sharedcache.SetToCache(ctx, t.store, key, payload, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionCache(ctx context.Context, id int) (*repository.TransactionResult, bool) {
	key := fmt.Sprintf(transactionByIdCacheKey, id)
	result, found := sharedcache.GetFromCache[repository.TransactionResult](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionQueryCache) SetCachedTransactionCache(ctx context.Context, data *repository.TransactionResult) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(transactionByIdCacheKey, data.TransactionID)
	sharedcache.SetToCache(ctx, t.store, key, data, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionByOrderId(ctx context.Context, orderID int) (*repository.TransactionResult, bool) {
	key := fmt.Sprintf(transactionByOrderCacheKey, orderID)
	result, found := sharedcache.GetFromCache[repository.TransactionResult](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionQueryCache) SetCachedTransactionByOrderId(ctx context.Context, orderID int, data *repository.TransactionResult) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(transactionByOrderCacheKey, orderID)
	sharedcache.SetToCache(ctx, t.store, key, data, ttlDefault)
}
