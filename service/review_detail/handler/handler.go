package handler

import (
	pbreview_detail "github.com/MamangRust/microservice-ecommerce-grpc-pb/review_detail"
	"github.com/MamangRust/microservice-ecommerce-grpc-review-detail/service"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	ReviewDetail        ReviewDetailHandleGrpc
	ReviewDetailQuery   pbreview_detail.ReviewDetailQueryServiceServer
	ReviewDetailCommand pbreview_detail.ReviewDetailCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		ReviewDetailQuery:   NewReviewDetailQueryHandler(deps.Service.ReviewDetailQuery, deps.Logger),
		ReviewDetailCommand: NewReviewDetailCommandHandler(deps.Service.ReviewDetailCommand, deps.Logger),
	}
}
