package repository

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type ReviewResult struct {
	ReviewID      int32
	UserID        int32
	ProductID     int32
	Name          string
	Comment       string
	Rating        int32
	CreatedAt     *string
	UpdatedAt     *string
	DeletedAt     *string
	TotalCount    int64
	ReviewDetails interface{}
}

type ReviewQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllReview) ([]*ReviewResult, error)
	FindByProduct(ctx context.Context, req *requests.FindAllReviewByProduct) ([]*ReviewResult, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllReviewByMerchant) ([]*ReviewResult, error)
	FindActive(ctx context.Context, req *requests.FindAllReview) ([]*ReviewResult, error)
	FindTrashed(ctx context.Context, req *requests.FindAllReview) ([]*ReviewResult, error)
	FindByID(ctx context.Context, id int) (*ReviewResult, error)
}

type ReviewCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateReviewRequest) (*models.Review, error)
	Update(ctx context.Context, request *requests.UpdateReviewRequest) (*models.Review, error)
	Trash(ctx context.Context, review_id int) (*models.Review, error)
	Restore(ctx context.Context, review_id int) (*models.Review, error)
	DeletePermanent(ctx context.Context, review_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
