package cache

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-merchant/repository"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type MerchantQueryCache interface {
	GetMerchantAllCache(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantResult, *int, bool)
	SetMerchantAllCache(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantResult, total *int)
	GetMerchantActiveCache(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantResult, *int, bool)
	SetMerchantActiveCache(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantResult, total *int)
	GetMerchantTrashedCache(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantResult, *int, bool)
	SetMerchantTrashedCache(ctx context.Context, req *requests.FindAllMerchant, data []*repository.MerchantResult, total *int)
	GetCachedMerchantCache(ctx context.Context, id int) (*repository.MerchantResult, bool)
	SetCachedMerchantCache(ctx context.Context, data *repository.MerchantResult)
	InvalidateMerchantCache(ctx context.Context)
}

type MerchantCommandCache interface {
	DeleteMerchantCache(ctx context.Context, merchantID int)
	InvalidateMerchantCache(ctx context.Context)
}

type MerchantDocumentQueryCache interface {
	GetCachedMerchantDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*repository.MerchantDocumentResult, *int, bool)
	SetCachedMerchantDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*repository.MerchantDocumentResult, total *int)
	GetCachedMerchantDocumentsActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*repository.MerchantDocumentResult, *int, bool)
	SetCachedMerchantDocumentsActive(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*repository.MerchantDocumentResult, total *int)
	GetCachedMerchantDocumentsTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*repository.MerchantDocumentResult, *int, bool)
	SetCachedMerchantDocumentsTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*repository.MerchantDocumentResult, total *int)
	GetCachedMerchantDocument(ctx context.Context, id int) (*repository.MerchantDocumentResult, bool)
	SetCachedMerchantDocument(ctx context.Context, data *repository.MerchantDocumentResult)
}

type MerchantDocumentCommandCache interface {
	DeleteCachedMerchantDocuments(ctx context.Context, id int)
}
