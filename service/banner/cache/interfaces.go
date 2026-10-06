package cache

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-banner/repository"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type BannerQueryCache interface {
	GetCachedBannersCache(ctx context.Context, req *requests.FindAllBanner) ([]*repository.BannerResult, *int, bool)
	SetCachedBannersCache(ctx context.Context, req *requests.FindAllBanner, data []*repository.BannerResult, total *int)

	GetCachedBannerActiveCache(ctx context.Context, req *requests.FindAllBanner) ([]*repository.BannerResult, *int, bool)
	SetCachedBannerActiveCache(ctx context.Context, req *requests.FindAllBanner, data []*repository.BannerResult, total *int)

	GetCachedBannerTrashedCache(ctx context.Context, req *requests.FindAllBanner) ([]*repository.BannerResult, *int, bool)
	SetCachedBannerTrashedCache(ctx context.Context, req *requests.FindAllBanner, data []*repository.BannerResult, total *int)

	GetCachedBannerCache(ctx context.Context, id int) (*repository.BannerResult, bool)
	SetCachedBannerCache(ctx context.Context, data *repository.BannerResult)
}

type BannerCommandCache interface {
	DeleteBannerCache(ctx context.Context, id int)
}
