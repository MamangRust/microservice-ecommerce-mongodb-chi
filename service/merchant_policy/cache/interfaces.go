package cache

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-merchant_policy/repository"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type MerchantPoliciesQueryCache interface {
	GetCachedMerchantPolicyAll(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantPolicyResult, *int, bool)
	SetCachedMerchantPolicyAll(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantPolicyResult, total *int)

	GetCachedMerchantPolicyActive(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantPolicyResult, *int, bool)
	SetCachedMerchantPolicyActive(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantPolicyResult, total *int)

	GetCachedMerchantPolicyTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantPolicyResult, *int, bool)
	SetCachedMerchantPolicyTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantPolicyResult, total *int)

	GetCachedMerchantPolicy(ctx context.Context, id int) (*repository.MerchantPolicyResult, bool)
	SetCachedMerchantPolicy(ctx context.Context, data *repository.MerchantPolicyResult)
}

type MerchantPoliciesCommandCache interface {
	DeleteMerchantPolicyCache(ctx context.Context, merchantID int)
}
