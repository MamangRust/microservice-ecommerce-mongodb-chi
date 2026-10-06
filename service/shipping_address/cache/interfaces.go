package cache

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-shipping-address/repository"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type ShippingAddressQueryCache interface {
	GetShippingAddressAllCache(ctx context.Context, req *requests.FindAllShippingAddress) ([]*repository.ShippingAddressResult, *int, bool)
	SetShippingAddressAllCache(ctx context.Context, req *requests.FindAllShippingAddress, res []*repository.ShippingAddressResult, total *int)

	GetShippingAddressTrashedCache(ctx context.Context, req *requests.FindAllShippingAddress) ([]*repository.ShippingAddressResult, *int, bool)
	SetShippingAddressTrashedCache(ctx context.Context, req *requests.FindAllShippingAddress, res []*repository.ShippingAddressResult, total *int)

	GetShippingAddressActiveCache(ctx context.Context, req *requests.FindAllShippingAddress) ([]*repository.ShippingAddressResult, *int, bool)
	SetShippingAddressActiveCache(ctx context.Context, req *requests.FindAllShippingAddress, res []*repository.ShippingAddressResult, total *int)

	GetCachedShippingAddressCache(ctx context.Context, shipping_id int) (*repository.ShippingAddressResult, bool)
	SetCachedShippingAddressCache(ctx context.Context, data *repository.ShippingAddressResult)

	GetCachedShippingAddressByOrderCache(ctx context.Context, order_id int) (*repository.ShippingAddressResult, bool)
	SetCachedShippingAddressByOrderCache(ctx context.Context, data *repository.ShippingAddressResult)
}

type ShippingAddressCommandCache interface {
	DeleteShippingAddressCache(ctx context.Context, shipping_id int)
	InvalidateShippingAddressCache(ctx context.Context)
}
