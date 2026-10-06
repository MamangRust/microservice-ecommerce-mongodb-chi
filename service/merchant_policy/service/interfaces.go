package service

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-merchant_policy/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type MerchantPoliciesQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantPolicyResult, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantPolicyResult, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantPolicyResult, *int, error)
	FindByID(ctx context.Context, id int) (*repository.MerchantPolicyResult, error)
}

type MerchantPoliciesCommandService interface {
	Create(ctx context.Context, request *requests.CreateMerchantPolicyRequest) (*models.MerchantPolicy, error)
	Update(ctx context.Context, request *requests.UpdateMerchantPolicyRequest) (*models.MerchantPolicy, error)
	Trash(ctx context.Context, id int) (*models.MerchantPolicy, error)
	Restore(ctx context.Context, id int) (*models.MerchantPolicy, error)
	DeletePermanent(ctx context.Context, id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
