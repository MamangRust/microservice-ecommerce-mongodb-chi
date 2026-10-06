package service

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-review/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type ReviewQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllReview) ([]*repository.ReviewResult, *int, error)
	FindByProduct(ctx context.Context, req *requests.FindAllReviewByProduct) ([]*repository.ReviewResult, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllReviewByMerchant) ([]*repository.ReviewResult, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllReview) ([]*repository.ReviewResult, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllReview) ([]*repository.ReviewResult, *int, error)
	FindByID(ctx context.Context, id int) (*repository.ReviewResult, error)
}

type ReviewCommandService interface {
	Create(ctx context.Context, request *requests.CreateReviewRequest) (*models.Review, error)
	Update(ctx context.Context, request *requests.UpdateReviewRequest) (*models.Review, error)
	Trash(ctx context.Context, reviewID int) (*models.Review, error)
	Restore(ctx context.Context, reviewID int) (*models.Review, error)
	DeletePermanent(ctx context.Context, reviewID int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
