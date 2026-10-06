package service

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-merchant_business/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type MerchantBusinessQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantBusinessResult, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantBusinessResult, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantBusinessResult, *int, error)
	FindByID(ctx context.Context, id int) (*repository.MerchantBusinessResult, error)
}

type MerchantBusinessCommandService interface {
	Create(ctx context.Context, request *requests.CreateMerchantBusinessInformationRequest) (*models.MerchantBusinessInformation, error)
	Update(ctx context.Context, request *requests.UpdateMerchantBusinessInformationRequest) (*models.MerchantBusinessInformation, error)
	Trash(ctx context.Context, id int) (*models.MerchantBusinessInformation, error)
	Restore(ctx context.Context, id int) (*models.MerchantBusinessInformation, error)
	DeletePermanent(ctx context.Context, id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
