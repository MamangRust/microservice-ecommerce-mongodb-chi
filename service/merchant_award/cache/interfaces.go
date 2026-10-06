package cache

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-merchant_award/repository"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type MerchantAwardQueryCache interface {
	GetCachedMerchantAwardAll(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantCertResult, *int, bool)
	SetCachedMerchantAwardAll(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantCertResult, total *int)
	GetCachedMerchantAwardActive(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantCertResult, *int, bool)
	SetCachedMerchantAwardActive(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantCertResult, total *int)
	GetCachedMerchantAwardTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantCertResult, *int, bool)
	SetCachedMerchantAwardTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantCertResult, total *int)
	GetCachedMerchantAward(ctx context.Context, id int) (*repository.MerchantCertResult, bool)
	SetCachedMerchantAward(ctx context.Context, data *repository.MerchantCertResult)
}

type MerchantAwardCommandCache interface {
	DeleteMerchantAwardCache(ctx context.Context, merchantID int)
}
