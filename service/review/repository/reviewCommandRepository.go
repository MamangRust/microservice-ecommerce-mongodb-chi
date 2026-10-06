package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	review_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/review"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type reviewCommandRepository struct {
	db *mongo.Database
}

func NewReviewCommandRepository(db *mongo.Database) *reviewCommandRepository {
	return &reviewCommandRepository{db: db}
}

func (r *reviewCommandRepository) Create(ctx context.Context, request *requests.CreateReviewRequest) (*models.Review, error) {
	id, err := database.GetNextID(ctx, r.db, "reviews")
	if err != nil {
		return nil, review_errors.ErrCreateReview.WithInternal(err)
	}
	now := time.Now()
	review := &models.Review{
		ReviewID:  id,
		UserID:    int32(request.UserID),
		ProductID: int32(request.ProductID),
		Comment:   request.Comment,
		Rating:    int32(request.Rating),
		CreatedAt: &now,
		UpdatedAt: &now,
	}
	if _, err := r.db.Collection("reviews").InsertOne(ctx, review); err != nil {
		return nil, review_errors.ErrCreateReview.WithInternal(err)
	}
	return review, nil
}

func (r *reviewCommandRepository) Update(ctx context.Context, request *requests.UpdateReviewRequest) (*models.Review, error) {
	coll := r.db.Collection("reviews")
	var review models.Review
	if err := coll.FindOne(ctx, bson.M{"_id": *request.ReviewID}).Decode(&review); err != nil {
		return nil, review_errors.ErrReviewNotFound
	}
	now := time.Now()
	update := bson.M{
		"name":       request.Name,
		"rating":     int32(request.Rating),
		"comment":    request.Comment,
		"updated_at": now,
	}
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": *request.ReviewID}, bson.M{"$set": update}); err != nil {
		return nil, review_errors.ErrUpdateReview.WithInternal(err)
	}
	if err := coll.FindOne(ctx, bson.M{"_id": *request.ReviewID}).Decode(&review); err != nil {
		return nil, review_errors.ErrUpdateReview.WithInternal(err)
	}
	return &review, nil
}

func (r *reviewCommandRepository) Trash(ctx context.Context, reviewID int) (*models.Review, error) {
	coll := r.db.Collection("reviews")
	var review models.Review
	if err := coll.FindOne(ctx, bson.M{"_id": reviewID}).Decode(&review); err != nil {
		return nil, review_errors.ErrReviewNotFound
	}
	now := time.Now()
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": reviewID}, bson.M{"$set": bson.M{"deleted_at": now}}); err != nil {
		return nil, review_errors.ErrTrashReview.WithInternal(err)
	}
	review.DeletedAt = &now
	return &review, nil
}

func (r *reviewCommandRepository) Restore(ctx context.Context, reviewID int) (*models.Review, error) {
	coll := r.db.Collection("reviews")
	res, err := coll.UpdateOne(ctx, bson.M{"_id": reviewID, "deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return nil, review_errors.ErrRestoreReview.WithInternal(err)
	}
	if res.MatchedCount == 0 {
		return nil, review_errors.ErrReviewNotFound
	}
	var review models.Review
	if err := coll.FindOne(ctx, bson.M{"_id": reviewID}).Decode(&review); err != nil {
		return nil, review_errors.ErrRestoreReview.WithInternal(err)
	}
	return &review, nil
}

func (r *reviewCommandRepository) DeletePermanent(ctx context.Context, reviewID int) (bool, error) {
	res, err := r.db.Collection("reviews").DeleteOne(ctx, bson.M{"_id": reviewID})
	if err != nil {
		return false, review_errors.ErrDeleteReviewPermanent.WithInternal(err)
	}
	if res.DeletedCount == 0 {
		return false, review_errors.ErrReviewNotFound
	}
	return true, nil
}

func (r *reviewCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("reviews").UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, review_errors.ErrRestoreAllReviews.WithInternal(err)
	}
	return true, nil
}

func (r *reviewCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("reviews").DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, review_errors.ErrDeleteAllPermanentReview.WithInternal(err)
	}
	return true, nil
}
