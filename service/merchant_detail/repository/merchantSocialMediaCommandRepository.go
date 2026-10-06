package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	merchant_social_link_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant_social_link_errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type merchantSocialLinkCommandRepository struct {
	db *mongo.Database
}

func NewMerchantSocialLinkCommandRepository(db *mongo.Database) *merchantSocialLinkCommandRepository {
	return &merchantSocialLinkCommandRepository{db: db}
}

func (r *merchantSocialLinkCommandRepository) Create(ctx context.Context, req *requests.CreateMerchantSocialRequest) (*models.MerchantSocialMediaLink, error) {
	id, err := database.GetNextID(ctx, r.db, "merchant_social_media_links")
	if err != nil {
		return nil, merchant_social_link_errors.ErrCreateMerchantSocialLink.WithInternal(err)
	}
	now := time.Now()
	link := &models.MerchantSocialMediaLink{
		MerchantSocialID: id,
		MerchantDetailID: int32(*req.MerchantDetailID),
		Platform:         req.Platform,
		Url:              req.Url,
		CreatedAt:        &now,
		UpdatedAt:        &now,
	}
	if _, err := r.db.Collection("merchant_social_media_links").InsertOne(ctx, link); err != nil {
		return nil, merchant_social_link_errors.ErrCreateMerchantSocialLink.WithInternal(err)
	}
	return link, nil
}

func (r *merchantSocialLinkCommandRepository) Update(ctx context.Context, req *requests.UpdateMerchantSocialRequest) (*models.MerchantSocialMediaLink, error) {
	coll := r.db.Collection("merchant_social_media_links")
	var link models.MerchantSocialMediaLink
	if err := coll.FindOne(ctx, bson.M{"_id": req.ID}).Decode(&link); err != nil {
		return nil, merchant_social_link_errors.ErrMerchantSocialLinkNotFound
	}
	now := time.Now()
	update := bson.M{
		"platform":    req.Platform,
		"url":         req.Url,
		"updated_at":  now,
	}
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": req.ID}, bson.M{"$set": update}); err != nil {
		return nil, merchant_social_link_errors.ErrUpdateMerchantSocialLink.WithInternal(err)
	}
	if err := coll.FindOne(ctx, bson.M{"_id": req.ID}).Decode(&link); err != nil {
		return nil, merchant_social_link_errors.ErrUpdateMerchantSocialLink.WithInternal(err)
	}
	return &link, nil
}

func (r *merchantSocialLinkCommandRepository) Trash(ctx context.Context, socialID int) (bool, error) {
	res, err := r.db.Collection("merchant_social_media_links").UpdateOne(ctx,
		bson.M{"_id": socialID, "deleted_at": bson.M{"$eq": nil}},
		bson.M{"$set": bson.M{"deleted_at": time.Now()}},
	)
	if err != nil {
		return false, merchant_social_link_errors.ErrTrashMerchantSocialLink.WithInternal(err)
	}
	if res.MatchedCount == 0 {
		return false, merchant_social_link_errors.ErrMerchantSocialLinkNotFound
	}
	return true, nil
}

func (r *merchantSocialLinkCommandRepository) Restore(ctx context.Context, socialID int) (bool, error) {
	res, err := r.db.Collection("merchant_social_media_links").UpdateOne(ctx,
		bson.M{"_id": socialID, "deleted_at": bson.M{"$ne": nil}},
		bson.M{"$set": bson.M{"deleted_at": nil}},
	)
	if err != nil {
		return false, merchant_social_link_errors.ErrRestoreMerchantSocialLink.WithInternal(err)
	}
	if res.MatchedCount == 0 {
		return false, merchant_social_link_errors.ErrMerchantSocialLinkNotFound
	}
	return true, nil
}

func (r *merchantSocialLinkCommandRepository) DeletePermanent(ctx context.Context, socialID int) (bool, error) {
	res, err := r.db.Collection("merchant_social_media_links").DeleteOne(ctx, bson.M{"_id": socialID})
	if err != nil {
		return false, merchant_social_link_errors.ErrDeletePermanentMerchantSocialLink.WithInternal(err)
	}
	if res.DeletedCount == 0 {
		return false, merchant_social_link_errors.ErrMerchantSocialLinkNotFound
	}
	return true, nil
}

func (r *merchantSocialLinkCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("merchant_social_media_links").UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, merchant_social_link_errors.ErrRestoreAllMerchantSocialLinks.WithInternal(err)
	}
	return true, nil
}

func (r *merchantSocialLinkCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("merchant_social_media_links").DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, merchant_social_link_errors.ErrDeleteAllPermanentMerchantSocialLinks.WithInternal(err)
	}
	return true, nil
}
