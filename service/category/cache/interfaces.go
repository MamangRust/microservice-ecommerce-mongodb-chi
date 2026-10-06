package cache

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-category/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type CategoryQueryCache interface {
	GetCachedCategoriesCache(ctx context.Context, req *requests.FindAllCategory) ([]*repository.CategoryResult, *int, bool)
	SetCachedCategoriesCache(ctx context.Context, req *requests.FindAllCategory, data []*repository.CategoryResult, total *int)

	GetCachedCategoryActiveCache(ctx context.Context, req *requests.FindAllCategory) ([]*repository.CategoryResult, *int, bool)
	SetCachedCategoryActiveCache(ctx context.Context, req *requests.FindAllCategory, data []*repository.CategoryResult, total *int)

	GetCachedCategoryTrashedCache(ctx context.Context, req *requests.FindAllCategory) ([]*repository.CategoryResult, *int, bool)
	SetCachedCategoryTrashedCache(ctx context.Context, req *requests.FindAllCategory, data []*repository.CategoryResult, total *int)

	GetCachedCategoryCache(ctx context.Context, id int) (*models.Category, bool)
	SetCachedCategoryCache(ctx context.Context, data *models.Category)
}

type CategoryCommandCache interface {
	DeleteCachedCategoryCache(ctx context.Context, id int)
}
