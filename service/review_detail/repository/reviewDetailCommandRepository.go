package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	review_detail_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/review_detail"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type reviewDetailCommandRepository struct {
	db *mongo.Database
}

func NewReviewDetailCommandRepository(db *mongo.Database) *reviewDetailCommandRepository {
	return &reviewDetailCommandRepository{db: db}
}

func (r *reviewDetailCommandRepository) Create(ctx context.Context, request *requests.CreateReviewDetailRequest) (*models.ReviewDetail, error) {
	id, err := database.GetNextID(ctx, r.db, "review_details")
	if err != nil {
		return nil, review_detail_errors.ErrCreateReviewDetail.WithInternal(err)
	}
	caption := request.Caption
	now := time.Now()
	detail := &models.ReviewDetail{
		ReviewDetailID: id,
		ReviewID:       int32(request.ReviewID),
		Type:           request.Type,
		Url:            request.Url,
		Caption:        &caption,
		CreatedAt:      &now,
		UpdatedAt:      &now,
	}
	if _, err := r.db.Collection("review_details").InsertOne(ctx, detail); err != nil {
		return nil, review_detail_errors.ErrCreateReviewDetail.WithInternal(err)
	}
	return detail, nil
}

func (r *reviewDetailCommandRepository) Update(ctx context.Context, request *requests.UpdateReviewDetailRequest) (*models.ReviewDetail, error) {
	coll := r.db.Collection("review_details")
	var detail models.ReviewDetail
	if err := coll.FindOne(ctx, bson.M{"_id": *request.ReviewDetailID}).Decode(&detail); err != nil {
		return nil, review_detail_errors.ErrReviewDetailNotFound
	}
	caption := request.Caption
	now := time.Now()
	update := bson.M{
		"type":       request.Type,
		"url":        request.Url,
		"caption":    &caption,
		"updated_at": now,
	}
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": *request.ReviewDetailID}, bson.M{"$set": update}); err != nil {
		return nil, review_detail_errors.ErrUpdateReviewDetail.WithInternal(err)
	}
	if err := coll.FindOne(ctx, bson.M{"_id": *request.ReviewDetailID}).Decode(&detail); err != nil {
		return nil, review_detail_errors.ErrUpdateReviewDetail.WithInternal(err)
	}
	return &detail, nil
}

func (r *reviewDetailCommandRepository) Trash(ctx context.Context, reviewDetailID int) (*models.ReviewDetail, error) {
	coll := r.db.Collection("review_details")
	var detail models.ReviewDetail
	if err := coll.FindOne(ctx, bson.M{"_id": reviewDetailID}).Decode(&detail); err != nil {
		return nil, review_detail_errors.ErrReviewDetailNotFound
	}
	now := time.Now()
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": reviewDetailID}, bson.M{"$set": bson.M{"deleted_at": now}}); err != nil {
		return nil, review_detail_errors.ErrTrashedReviewDetail.WithInternal(err)
	}
	detail.DeletedAt = &now
	return &detail, nil
}

func (r *reviewDetailCommandRepository) Restore(ctx context.Context, reviewDetailID int) (*models.ReviewDetail, error) {
	coll := r.db.Collection("review_details")
	res, err := coll.UpdateOne(ctx, bson.M{"_id": reviewDetailID, "deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return nil, review_detail_errors.ErrRestoreReviewDetail.WithInternal(err)
	}
	if res.MatchedCount == 0 {
		return nil, review_detail_errors.ErrReviewDetailNotFound
	}
	var detail models.ReviewDetail
	if err := coll.FindOne(ctx, bson.M{"_id": reviewDetailID}).Decode(&detail); err != nil {
		return nil, review_detail_errors.ErrRestoreReviewDetail.WithInternal(err)
	}
	return &detail, nil
}

func (r *reviewDetailCommandRepository) DeletePermanent(ctx context.Context, reviewDetailID int) (bool, error) {
	res, err := r.db.Collection("review_details").DeleteOne(ctx, bson.M{"_id": reviewDetailID})
	if err != nil {
		return false, review_detail_errors.ErrDeleteReviewDetailPermanent.WithInternal(err)
	}
	if res.DeletedCount == 0 {
		return false, review_detail_errors.ErrReviewDetailNotFound
	}
	return true, nil
}

func (r *reviewDetailCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("review_details").UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, review_detail_errors.ErrRestoreAllReviewDetails.WithInternal(err)
	}
	return true, nil
}

func (r *reviewDetailCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("review_details").DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, review_detail_errors.ErrDeleteAllReviewDetails.WithInternal(err)
	}
	return true, nil
}
