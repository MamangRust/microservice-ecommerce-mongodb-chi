package service

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-merchant_detail/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type MerchantDetailQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantDetailResult, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantDetailResult, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantDetailResult, *int, error)
	FindByID(ctx context.Context, user_id int) (*repository.MerchantDetailResult, error)
}

type MerchantDetailCommandService interface {
	Create(ctx context.Context, request *requests.CreateMerchantDetailRequest) (*models.MerchantDetail, error)
	Update(ctx context.Context, request *requests.UpdateMerchantDetailRequest) (*models.MerchantDetail, error)
	Trash(ctx context.Context, merchant_id int) (*models.MerchantDetail, error)
	Restore(ctx context.Context, merchant_id int) (*models.MerchantDetail, error)
	DeletePermanent(ctx context.Context, merchant_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

type MerchantSocialLinkCommandService interface {
	Create(ctx context.Context, req *requests.CreateMerchantSocialRequest) (*models.MerchantSocialMediaLink, error)
	Update(ctx context.Context, req *requests.UpdateMerchantSocialRequest) (*models.MerchantSocialMediaLink, error)
	Trash(ctx context.Context, socialID int) (bool, error)
	Restore(ctx context.Context, socialID int) (bool, error)
	DeletePermanent(ctx context.Context, socialID int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
