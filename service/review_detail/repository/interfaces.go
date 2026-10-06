package repository

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type ReviewDetailResult struct {
	ReviewDetailID int32
	ReviewID       int32
	Type           string
	Url            string
	Caption        *string
	CreatedAt      *string
	UpdatedAt      *string
	DeletedAt      *string
	TotalCount     int64
}

type ReviewDetailQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllReview) ([]*ReviewDetailResult, error)
	FindActive(ctx context.Context, req *requests.FindAllReview) ([]*ReviewDetailResult, error)
	FindTrashed(ctx context.Context, req *requests.FindAllReview) ([]*ReviewDetailResult, error)
	FindByID(ctx context.Context, id int) (*ReviewDetailResult, error)
}

type ReviewDetailCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateReviewDetailRequest) (*models.ReviewDetail, error)
	Update(ctx context.Context, request *requests.UpdateReviewDetailRequest) (*models.ReviewDetail, error)
	Trash(ctx context.Context, reviewDetailID int) (*models.ReviewDetail, error)
	Restore(ctx context.Context, reviewDetailID int) (*models.ReviewDetail, error)
	DeletePermanent(ctx context.Context, reviewDetailID int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
