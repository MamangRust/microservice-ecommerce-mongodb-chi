package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type MerchantPolicyResult struct {
	MerchantPolicyID int32
	MerchantID       int32
	PolicyType       string
	Title            string
	Description      string
	MerchantName     string
	CreatedAt        *time.Time
	UpdatedAt        *time.Time
	DeletedAt        *time.Time
	TotalCount       int64
}

type MerchantPoliciesQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantPolicyResult, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantPolicyResult, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantPolicyResult, error)
	FindByID(ctx context.Context, id int) (*MerchantPolicyResult, error)
}

type MerchantPoliciesCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateMerchantPolicyRequest) (*models.MerchantPolicy, error)
	Update(ctx context.Context, request *requests.UpdateMerchantPolicyRequest) (*models.MerchantPolicy, error)
	Trash(ctx context.Context, id int) (*models.MerchantPolicy, error)
	Restore(ctx context.Context, id int) (*models.MerchantPolicy, error)
	DeletePermanent(ctx context.Context, id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

type MerchantQueryRepository interface {
	FindByID(ctx context.Context, id int) (string, error)
}
