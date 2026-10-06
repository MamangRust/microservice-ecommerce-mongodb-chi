package service

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-review-detail/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type ReviewDetailQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllReview) ([]*repository.ReviewDetailResult, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllReview) ([]*repository.ReviewDetailResult, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllReview) ([]*repository.ReviewDetailResult, *int, error)
	FindByID(ctx context.Context, id int) (*repository.ReviewDetailResult, error)
}

type ReviewDetailCommandService interface {
	Create(ctx context.Context, request *requests.CreateReviewDetailRequest) (*models.ReviewDetail, error)
	Update(ctx context.Context, request *requests.UpdateReviewDetailRequest) (*models.ReviewDetail, error)
	Trash(ctx context.Context, reviewDetailID int) (*models.ReviewDetail, error)
	Restore(ctx context.Context, reviewDetailID int) (*models.ReviewDetail, error)
	DeletePermanent(ctx context.Context, reviewDetailID int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
