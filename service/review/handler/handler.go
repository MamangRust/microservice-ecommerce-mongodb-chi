package handler

import (
	pbreview "github.com/MamangRust/microservice-ecommerce-grpc-pb/review"
	"github.com/MamangRust/microservice-ecommerce-grpc-review/service"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	ReviewQuery   pbreview.ReviewQueryServiceServer
	ReviewCommand pbreview.ReviewCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		ReviewQuery:   NewReviewQueryHandler(deps.Service.ReviewQuery, deps.Logger),
		ReviewCommand: NewReviewCommandHandler(deps.Service.ReviewCommand, deps.Logger),
	}
}

type reviewHandleGrpc struct {
	// Dummy struct for mapping receiver
}
