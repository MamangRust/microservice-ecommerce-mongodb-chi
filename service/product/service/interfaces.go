package service

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-product/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type ProductQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllProduct) ([]*repository.ProductResult, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllProduct) ([]*repository.ProductResult, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*repository.ProductResult, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant) ([]*repository.ProductResult, *int, error)
	FindByCategory(ctx context.Context, req *requests.FindAllProductByCategory) ([]*repository.ProductResult, *int, error)
	FindByID(ctx context.Context, product_id int) (*repository.ProductResult, error)
}

type ProductCommandService interface {
	Create(ctx context.Context, req *requests.CreateProductRequest) (*models.Product, error)
	Update(ctx context.Context, req *requests.UpdateProductRequest) (*models.Product, error)
	UpdateProductCountStock(ctx context.Context, productID int, stock int) (*models.Product, error)
	AdjustProductStock(ctx context.Context, productID int, delta int, operationID string) (*models.Product, error)
	Trash(ctx context.Context, productID int) (interface{}, error)
	Restore(ctx context.Context, productID int) (interface{}, error)
	DeletePermanent(ctx context.Context, productID int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
