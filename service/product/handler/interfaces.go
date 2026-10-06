package handler

import (
	pbproduct "github.com/MamangRust/microservice-ecommerce-grpc-pb/product"
)

type ProductQueryHandler interface {
	pbproduct.ProductQueryServiceServer
}

type ProductCommandHandler interface {
	pbproduct.ProductCommandServiceServer
}
