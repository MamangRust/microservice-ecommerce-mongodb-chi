package handler

import (
	pbcart "github.com/MamangRust/microservice-ecommerce-grpc-pb/cart"
)

type CartQueryHandler interface {
	pbcart.CartQueryServiceServer
}

type CartCommandHandler interface {
	pbcart.CartCommandServiceServer
}
