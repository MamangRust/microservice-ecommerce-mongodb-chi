package service

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-merchant_award/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type MerchantAwardQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantCertResult, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantCertResult, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantCertResult, *int, error)
	FindByID(ctx context.Context, id int) (*repository.MerchantCertResult, error)
}

type MerchantAwardCommandService interface {
	Create(ctx context.Context, request *requests.CreateMerchantCertificationOrAwardRequest) (*models.MerchantCertificationsAndAward, error)
	Update(ctx context.Context, request *requests.UpdateMerchantCertificationOrAwardRequest) (*models.MerchantCertificationsAndAward, error)
	Trash(ctx context.Context, id int) (*models.MerchantCertificationsAndAward, error)
	Restore(ctx context.Context, id int) (*models.MerchantCertificationsAndAward, error)
	DeletePermanent(ctx context.Context, id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
