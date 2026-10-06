package cache

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-merchant_business/repository"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type MerchantBusinessQueryCache interface {
	GetCachedMerchantBusinessAll(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantBusinessResult, *int, bool)
	SetCachedMerchantBusinessAll(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantBusinessResult, total *int)
	GetCachedMerchantBusinessActive(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantBusinessResult, *int, bool)
	SetCachedMerchantBusinessActive(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantBusinessResult, total *int)
	GetCachedMerchantBusinessTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantBusinessResult, *int, bool)
	SetCachedMerchantBusinessTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantBusinessResult, total *int)
	GetCachedMerchantBusiness(ctx context.Context, id int) (*repository.MerchantBusinessResult, bool)
	SetCachedMerchantBusiness(ctx context.Context, data *repository.MerchantBusinessResult)
}

type MerchantBusinessCommandCache interface {
	DeleteMerchantBusinessCache(ctx context.Context, merchantID int)
}
