package cache

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-product/repository"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type ProductQueryCache interface {
	GetCachedProducts(ctx context.Context, req *requests.FindAllProduct) ([]*repository.ProductResult, *int, bool)
	SetCachedProducts(ctx context.Context, req *requests.FindAllProduct, data []*repository.ProductResult, total *int)

	GetCachedProductsByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant) ([]*repository.ProductResult, *int, bool)
	SetCachedProductsByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant, data []*repository.ProductResult, total *int)

	GetCachedProductsByCategory(ctx context.Context, req *requests.FindAllProductByCategory) ([]*repository.ProductResult, *int, bool)
	SetCachedProductsByCategory(ctx context.Context, req *requests.FindAllProductByCategory, data []*repository.ProductResult, total *int)

	GetCachedProductActive(ctx context.Context, req *requests.FindAllProduct) ([]*repository.ProductResult, *int, bool)
	SetCachedProductActive(ctx context.Context, req *requests.FindAllProduct, data []*repository.ProductResult, total *int)

	GetCachedProductTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*repository.ProductResult, *int, bool)
	SetCachedProductTrashed(ctx context.Context, req *requests.FindAllProduct, data []*repository.ProductResult, total *int)

	GetCachedProduct(ctx context.Context, productID int) (*repository.ProductResult, bool)
	SetCachedProduct(ctx context.Context, data *repository.ProductResult)
}

type ProductCommandCache interface {
	DeleteCachedProduct(ctx context.Context, productID int)
}
