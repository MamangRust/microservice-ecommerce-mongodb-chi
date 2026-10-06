package cache

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-review/repository"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type ReviewQueryCache interface {
	GetReviewAllCache(ctx context.Context, req *requests.FindAllReview) ([]*repository.ReviewResult, *int, bool)
	SetReviewAllCache(ctx context.Context, req *requests.FindAllReview, data []*repository.ReviewResult, total *int)
	GetReviewByProductCache(ctx context.Context, req *requests.FindAllReviewByProduct) ([]*repository.ReviewResult, *int, bool)
	SetReviewByProductCache(ctx context.Context, req *requests.FindAllReviewByProduct, data []*repository.ReviewResult, total *int)
	GetReviewByMerchantCache(ctx context.Context, req *requests.FindAllReviewByMerchant) ([]*repository.ReviewResult, *int, bool)
	SetReviewByMerchantCache(ctx context.Context, req *requests.FindAllReviewByMerchant, data []*repository.ReviewResult, total *int)
	GetReviewActiveCache(ctx context.Context, req *requests.FindAllReview) ([]*repository.ReviewResult, *int, bool)
	SetReviewActiveCache(ctx context.Context, req *requests.FindAllReview, data []*repository.ReviewResult, total *int)
	GetReviewTrashedCache(ctx context.Context, req *requests.FindAllReview) ([]*repository.ReviewResult, *int, bool)
	SetReviewTrashedCache(ctx context.Context, req *requests.FindAllReview, data []*repository.ReviewResult, total *int)
	GetReviewByIdCache(ctx context.Context, id int) (*repository.ReviewResult, bool)
	SetReviewByIdCache(ctx context.Context, data *repository.ReviewResult)
}

type ReviewCommandCache interface {
	DeleteReviewCache(ctx context.Context, reviewID int)
}
