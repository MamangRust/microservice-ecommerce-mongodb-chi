package service

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-banner/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type BannerQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllBanner) ([]*repository.BannerResult, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllBanner) ([]*repository.BannerResult, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllBanner) ([]*repository.BannerResult, *int, error)
	FindByID(ctx context.Context, bannerID int) (*repository.BannerResult, error)
}

type BannerCommandService interface {
	Create(ctx context.Context, req *requests.CreateBannerRequest) (*models.Banner, error)
	Update(ctx context.Context, req *requests.UpdateBannerRequest) (*models.Banner, error)

	Trash(ctx context.Context, bannerID int) (*models.Banner, error)
	Restore(ctx context.Context, bannerID int) (*models.Banner, error)
	DeletePermanent(ctx context.Context, bannerID int) (bool, error)

	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
