package service

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-category/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type CategoryQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllCategory) ([]*repository.CategoryResult, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllCategory) ([]*repository.CategoryResult, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllCategory) ([]*repository.CategoryResult, *int, error)
	FindByID(ctx context.Context, categoryID int) (*models.Category, error)
}

type CategoryCommandService interface {
	Create(ctx context.Context, req *requests.CreateCategoryRequest) (*models.Category, error)
	Update(ctx context.Context, req *requests.UpdateCategoryRequest) (*models.Category, error)
	Trash(ctx context.Context, categoryID int) (*models.Category, error)
	Restore(ctx context.Context, categoryID int) (*models.Category, error)
	DeletePermanent(ctx context.Context, categoryID int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
