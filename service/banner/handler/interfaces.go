package handler

import (
	pbbanner "github.com/MamangRust/microservice-ecommerce-grpc-pb/banner"
)

type BannerQueryHandler interface {
	pbbanner.BannerQueryServiceServer
}

type BannerCommandHandler interface {
	pbbanner.BannerCommandServiceServer
}
