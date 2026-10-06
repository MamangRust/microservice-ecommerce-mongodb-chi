package cache

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-review-detail/repository"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type ReviewDetailQueryCache interface {
	GetReviewDetailAllCache(ctx context.Context, req *requests.FindAllReview) ([]*repository.ReviewDetailResult, *int, bool)
	SetReviewDetailAllCache(ctx context.Context, req *requests.FindAllReview, data []*repository.ReviewDetailResult, total *int)
	GetReviewDetailActiveCache(ctx context.Context, req *requests.FindAllReview) ([]*repository.ReviewDetailResult, *int, bool)
	SetReviewDetailActiveCache(ctx context.Context, req *requests.FindAllReview, data []*repository.ReviewDetailResult, total *int)
	GetReviewDetailTrashedCache(ctx context.Context, req *requests.FindAllReview) ([]*repository.ReviewDetailResult, *int, bool)
	SetReviewDetailTrashedCache(ctx context.Context, req *requests.FindAllReview, data []*repository.ReviewDetailResult, total *int)
	GetCachedReviewDetailCache(ctx context.Context, reviewID int) (*repository.ReviewDetailResult, bool)
	SetCachedReviewDetailCache(ctx context.Context, data *repository.ReviewDetailResult)
}

type ReviewDetailCommandCache interface {
	DeleteReviewDetailCache(ctx context.Context, reviewID int)
}
