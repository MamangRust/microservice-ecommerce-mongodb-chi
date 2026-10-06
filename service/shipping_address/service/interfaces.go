package service

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-shipping-address/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type ShippingAddressQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllShippingAddress) ([]*repository.ShippingAddressResult, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllShippingAddress) ([]*repository.ShippingAddressResult, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllShippingAddress) ([]*repository.ShippingAddressResult, *int, error)
	FindByOrder(ctx context.Context, shipping_id int) (*repository.ShippingAddressResult, error)
	FindByID(ctx context.Context, shipping_id int) (*repository.ShippingAddressResult, error)
}

type ShippingAddressCommandService interface {
	Create(ctx context.Context, request *requests.CreateShippingAddressRequest) (*models.ShippingAddress, error)
	Update(ctx context.Context, request *requests.UpdateShippingAddressRequest) (*models.ShippingAddress, error)
	Trash(ctx context.Context, shipping_id int) (*models.ShippingAddress, error)
	Restore(ctx context.Context, shipping_id int) (*models.ShippingAddress, error)
	DeletePermanent(ctx context.Context, shipping_id int) (bool, error)
	DeleteShippingAddressByOrderPermanent(ctx context.Context, order_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
