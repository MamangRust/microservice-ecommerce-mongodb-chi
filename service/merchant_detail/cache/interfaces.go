package cache

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-merchant_detail/repository"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type MerchantDetailQueryCache interface {
	GetCachedMerchantDetailsAll(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantDetailResult, *int, bool)
	SetCachedMerchantDetailsAll(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantDetailResult, total *int)
	GetCachedMerchantDetailsActive(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantDetailResult, *int, bool)
	SetCachedMerchantDetailsActive(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantDetailResult, total *int)
	GetCachedMerchantDetailsTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantDetailResult, *int, bool)
	SetCachedMerchantDetailsTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantDetailResult, total *int)
	GetCachedMerchantDetail(ctx context.Context, id int) (*repository.MerchantDetailResult, bool)
	SetCachedMerchantDetail(ctx context.Context, data *repository.MerchantDetailResult)
}

type MerchantDetailCommandCache interface {
	DeleteMerchantDetailCache(ctx context.Context, merchantID int)
}
