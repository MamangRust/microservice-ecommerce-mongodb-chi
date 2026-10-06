package repository

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type BannerResult struct {
	BannerID   int32
	Name       string
	StartDate  *string
	EndDate    *string
	StartTime  *string
	EndTime    *string
	IsActive   *bool
	CreatedAt  *string
	UpdatedAt  *string
	DeletedAt  *string
	TotalCount int64
}

type BannerQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllBanner) ([]*BannerResult, error)
	FindActive(ctx context.Context, req *requests.FindAllBanner) ([]*BannerResult, error)
	FindTrashed(ctx context.Context, req *requests.FindAllBanner) ([]*BannerResult, error)
	FindByID(ctx context.Context, banner_id int) (*BannerResult, error)
}

type BannerCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateBannerRequest) (*models.Banner, error)
	Update(ctx context.Context, request *requests.UpdateBannerRequest) (*models.Banner, error)

	Trash(ctx context.Context, banner_id int) (*models.Banner, error)
	Restore(ctx context.Context, banner_id int) (*models.Banner, error)
	DeletePermanent(ctx context.Context, banner_id int) (bool, error)

	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
