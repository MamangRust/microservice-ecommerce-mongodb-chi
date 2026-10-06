package cache

import (
	"context"
	"fmt"

	"github.com/MamangRust/microservice-ecommerce-grpc-user/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type userCacheResponseAll struct {
	Data         []*repository.UserResult `json:"data"`
	TotalRecords *int                     `json:"total_records"`
}

type userCacheResponseActive struct {
	Data         []*repository.UserResult `json:"data"`
	TotalRecords *int                     `json:"total_records"`
}

type userCacheResponseTrashed struct {
	Data         []*repository.UserResult `json:"data"`
	TotalRecords *int                     `json:"total_records"`
}

type userQueryCache struct {
	store *cache.CacheStore
}

func NewUserQueryCache(store *cache.CacheStore) UserQueryCache {
	return &userQueryCache{store: store}
}

func (s *userQueryCache) GetCachedUsersCache(ctx context.Context, req *requests.FindAllUsers) ([]*repository.UserResult, *int, bool) {
	key := fmt.Sprintf(userAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[userCacheResponseAll](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *userQueryCache) SetCachedUsersCache(ctx context.Context, req *requests.FindAllUsers, data []*repository.UserResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.UserResult{}
	}

	key := fmt.Sprintf(userAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &userCacheResponseAll{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *userQueryCache) GetCachedUserActiveCache(ctx context.Context, req *requests.FindAllUsers) ([]*repository.UserResult, *int, bool) {
	key := fmt.Sprintf(userActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[userCacheResponseActive](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *userQueryCache) SetCachedUserActiveCache(ctx context.Context, req *requests.FindAllUsers, data []*repository.UserResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.UserResult{}
	}

	key := fmt.Sprintf(userActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &userCacheResponseActive{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *userQueryCache) GetCachedUserTrashedCache(ctx context.Context, req *requests.FindAllUsers) ([]*repository.UserResult, *int, bool) {
	key := fmt.Sprintf(userTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[userCacheResponseTrashed](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *userQueryCache) SetCachedUserTrashedCache(ctx context.Context, req *requests.FindAllUsers, data []*repository.UserResult, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*repository.UserResult{}
	}

	key := fmt.Sprintf(userTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &userCacheResponseTrashed{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *userQueryCache) GetCachedUserCache(ctx context.Context, id int) (*models.User, bool) {
	key := fmt.Sprintf(userByIdCacheKey, id)

	result, found := cache.GetFromCache[*models.User](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *userQueryCache) SetCachedUserCache(ctx context.Context, data *models.User) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(userByIdCacheKey, data.UserID)
	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}
