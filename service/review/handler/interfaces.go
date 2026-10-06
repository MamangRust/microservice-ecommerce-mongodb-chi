package handler

import (
	pbreview "github.com/MamangRust/microservice-ecommerce-grpc-pb/review"
)

type ReviewHandleGrpc interface {
	pbreview.ReviewQueryServiceServer
	pbreview.ReviewCommandServiceServer
}

type ReviewQueryHandler interface {
	pbreview.ReviewQueryServiceServer
}

type ReviewCommandHandler interface {
	pbreview.ReviewCommandServiceServer
}
