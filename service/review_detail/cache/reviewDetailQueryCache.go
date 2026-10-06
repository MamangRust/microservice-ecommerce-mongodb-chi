package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/microservice-ecommerce-grpc-review-detail/repository"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

const (
	reviewDetailAllCacheKey         = "review_detail:all:page:%d:pageSize:%d:search:%s"
	reviewDetailByIdCacheKey        = "review_detail:id:%d"
	reviewDetailActiveCacheKey      = "review_detail:active:page:%d:pageSize:%d:search:%s"
	reviewDetailTrashedCacheKey     = "review_detail:trashed:page:%d:pageSize:%d:search:%s"
	reviewDetailByIdTrashedCacheKey = "review_detail:id_trashed:%d"

	ttlDefault = 5 * time.Minute
)

type reviewDetailCacheResponseDB struct {
	Data  []*repository.ReviewDetailResult `json:"data"`
	Total *int                             `json:"total_records"`
}

type reviewDetailActiveCacheResponseDB struct {
	Data  []*repository.ReviewDetailResult `json:"data"`
	Total *int                             `json:"total_records"`
}

type reviewDetailTrashedCacheResponseDB struct {
	Data  []*repository.ReviewDetailResult `json:"data"`
	Total *int                             `json:"total_records"`
}

type reviewDetailQueryCache struct {
	store *cache.CacheStore
}

func NewReviewDetailQueryCache(store *cache.CacheStore) *reviewDetailQueryCache {
	return &reviewDetailQueryCache{store: store}
}

func (r *reviewDetailQueryCache) GetReviewDetailAllCache(ctx context.Context, req *requests.FindAllReview) ([]*repository.ReviewDetailResult, *int, bool) {
	key := fmt.Sprintf(reviewDetailAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[reviewDetailCacheResponseDB](ctx, r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (r *reviewDetailQueryCache) SetReviewDetailAllCache(ctx context.Context, req *requests.FindAllReview, data []*repository.ReviewDetailResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.ReviewDetailResult{}
	}

	key := fmt.Sprintf(reviewDetailAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &reviewDetailCacheResponseDB{Data: data, Total: total}
	cache.SetToCache(ctx, r.store, key, payload, ttlDefault)
}

func (r *reviewDetailQueryCache) GetReviewDetailActiveCache(ctx context.Context, req *requests.FindAllReview) ([]*repository.ReviewDetailResult, *int, bool) {
	key := fmt.Sprintf(reviewDetailActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[reviewDetailActiveCacheResponseDB](ctx, r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (r *reviewDetailQueryCache) SetReviewDetailActiveCache(ctx context.Context, req *requests.FindAllReview, data []*repository.ReviewDetailResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.ReviewDetailResult{}
	}

	key := fmt.Sprintf(reviewDetailActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &reviewDetailActiveCacheResponseDB{Data: data, Total: total}
	cache.SetToCache(ctx, r.store, key, payload, ttlDefault)
}

func (r *reviewDetailQueryCache) GetReviewDetailTrashedCache(ctx context.Context, req *requests.FindAllReview) ([]*repository.ReviewDetailResult, *int, bool) {
	key := fmt.Sprintf(reviewDetailTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[reviewDetailTrashedCacheResponseDB](ctx, r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (r *reviewDetailQueryCache) SetReviewDetailTrashedCache(ctx context.Context, req *requests.FindAllReview, data []*repository.ReviewDetailResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.ReviewDetailResult{}
	}

	key := fmt.Sprintf(reviewDetailTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &reviewDetailTrashedCacheResponseDB{Data: data, Total: total}
	cache.SetToCache(ctx, r.store, key, payload, ttlDefault)
}

func (r *reviewDetailQueryCache) GetCachedReviewDetailCache(ctx context.Context, reviewID int) (*repository.ReviewDetailResult, bool) {
	key := fmt.Sprintf(reviewDetailByIdCacheKey, reviewID)
	result, found := cache.GetFromCache[repository.ReviewDetailResult](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (r *reviewDetailQueryCache) SetCachedReviewDetailCache(ctx context.Context, data *repository.ReviewDetailResult) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(reviewDetailByIdCacheKey, data.ReviewDetailID)
	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}

func (r *reviewDetailQueryCache) GetCachedReviewDetailTrashedCache(ctx context.Context, reviewID int) (*repository.ReviewDetailResult, bool) {
	key := fmt.Sprintf(reviewDetailByIdTrashedCacheKey, reviewID)
	result, found := cache.GetFromCache[repository.ReviewDetailResult](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (r *reviewDetailQueryCache) SetCachedReviewDetailTrashedCache(ctx context.Context, data *repository.ReviewDetailResult) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(reviewDetailByIdTrashedCacheKey, data.ReviewDetailID)
	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}
