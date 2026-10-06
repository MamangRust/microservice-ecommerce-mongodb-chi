package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type CategoryResult struct {
	CategoryID    int32
	Name          string
	Description   *string
	SlugCategory  *string
	ImageCategory *string
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
	DeletedAt     *time.Time
	TotalCount    int64
}

type CategoryQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllCategory) ([]*CategoryResult, error)
	FindActive(ctx context.Context, req *requests.FindAllCategory) ([]*CategoryResult, error)
	FindTrashed(ctx context.Context, req *requests.FindAllCategory) ([]*CategoryResult, error)
	FindByID(ctx context.Context, category_id int) (*models.Category, error)
	FindByIDTrashed(ctx context.Context, category_id int) (*models.Category, error)
}

type CategoryCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateCategoryRequest) (*models.Category, error)
	Update(ctx context.Context, request *requests.UpdateCategoryRequest) (*models.Category, error)
	Trash(ctx context.Context, category_id int) (*models.Category, error)
	Restore(ctx context.Context, category_id int) (*models.Category, error)
	DeletePermanent(ctx context.Context, category_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
