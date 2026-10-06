package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	merchantdetail_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant_detail"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type merchantDetailCommandRepository struct {
	db *mongo.Database
}

func NewMerchantDetailCommandRepository(db *mongo.Database) *merchantDetailCommandRepository {
	return &merchantDetailCommandRepository{db: db}
}

func (r *merchantDetailCommandRepository) Create(ctx context.Context, request *requests.CreateMerchantDetailRequest) (*models.MerchantDetail, error) {
	id, err := database.GetNextID(ctx, r.db, "merchant_details")
	if err != nil {
		return nil, merchantdetail_errors.ErrCreateMerchantDetail.WithInternal(err)
	}
	now := time.Now()
	detail := &models.MerchantDetail{
		MerchantDetailID: id,
		MerchantID:       int32(request.MerchantID),
		DisplayName:      &request.DisplayName,
		CoverImageUrl:    &request.CoverImageUrl,
		LogoUrl:          &request.LogoUrl,
		ShortDescription: &request.ShortDescription,
		WebsiteUrl:       &request.WebsiteUrl,
		CreatedAt:        &now,
		UpdatedAt:        &now,
	}
	if _, err := r.db.Collection("merchant_details").InsertOne(ctx, detail); err != nil {
		return nil, merchantdetail_errors.ErrCreateMerchantDetail.WithInternal(err)
	}
	return detail, nil
}

func (r *merchantDetailCommandRepository) Update(ctx context.Context, request *requests.UpdateMerchantDetailRequest) (*models.MerchantDetail, error) {
	coll := r.db.Collection("merchant_details")
	var detail models.MerchantDetail
	if err := coll.FindOne(ctx, bson.M{"_id": *request.MerchantDetailID}).Decode(&detail); err != nil {
		return nil, merchantdetail_errors.ErrMerchantDetailNotFound
	}
	now := time.Now()
	update := bson.M{
		"display_name":      &request.DisplayName,
		"cover_image_url":   &request.CoverImageUrl,
		"logo_url":          &request.LogoUrl,
		"short_description": &request.ShortDescription,
		"website_url":       &request.WebsiteUrl,
		"updated_at":        now,
	}
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": *request.MerchantDetailID}, bson.M{"$set": update}); err != nil {
		return nil, merchantdetail_errors.ErrUpdateMerchantDetail.WithInternal(err)
	}
	if err := coll.FindOne(ctx, bson.M{"_id": *request.MerchantDetailID}).Decode(&detail); err != nil {
		return nil, merchantdetail_errors.ErrUpdateMerchantDetail.WithInternal(err)
	}
	return &detail, nil
}

func (r *merchantDetailCommandRepository) Trash(ctx context.Context, merchantDetailID int) (*models.MerchantDetail, error) {
	coll := r.db.Collection("merchant_details")
	var detail models.MerchantDetail
	if err := coll.FindOne(ctx, bson.M{"_id": merchantDetailID}).Decode(&detail); err != nil {
		return nil, merchantdetail_errors.ErrMerchantDetailNotFound
	}
	now := time.Now()
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": merchantDetailID}, bson.M{"$set": bson.M{"deleted_at": now}}); err != nil {
		return nil, merchantdetail_errors.ErrTrashMerchantDetail.WithInternal(err)
	}
	detail.DeletedAt = &now
	return &detail, nil
}

func (r *merchantDetailCommandRepository) Restore(ctx context.Context, merchantDetailID int) (*models.MerchantDetail, error) {
	coll := r.db.Collection("merchant_details")
	res, err := coll.UpdateOne(ctx, bson.M{"_id": merchantDetailID, "deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return nil, merchantdetail_errors.ErrRestoreMerchantDetail.WithInternal(err)
	}
	if res.MatchedCount == 0 {
		return nil, merchantdetail_errors.ErrMerchantDetailNotFound
	}
	var detail models.MerchantDetail
	if err := coll.FindOne(ctx, bson.M{"_id": merchantDetailID}).Decode(&detail); err != nil {
		return nil, merchantdetail_errors.ErrRestoreMerchantDetail.WithInternal(err)
	}
	return &detail, nil
}

func (r *merchantDetailCommandRepository) DeletePermanent(ctx context.Context, merchantDetailID int) (bool, error) {
	res, err := r.db.Collection("merchant_details").DeleteOne(ctx, bson.M{"_id": merchantDetailID})
	if err != nil {
		return false, merchantdetail_errors.ErrDeletePermanentMerchantDetail.WithInternal(err)
	}
	if res.DeletedCount == 0 {
		return false, merchantdetail_errors.ErrMerchantDetailNotFound
	}
	return true, nil
}

func (r *merchantDetailCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("merchant_details").UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, merchantdetail_errors.ErrRestoreAllMerchantDetails.WithInternal(err)
	}
	return true, nil
}

func (r *merchantDetailCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("merchant_details").DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, merchantdetail_errors.ErrDeleteAllPermanentMerchantDetails.WithInternal(err)
	}
	return true, nil
}
