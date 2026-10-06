package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	merchantaward_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant_award"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type merchantAwardCommandRepository struct {
	db *mongo.Database
}

func NewMerchantAwardCommandRepository(db *mongo.Database) MerchantAwardCommandRepository {
	return &merchantAwardCommandRepository{db: db}
}

func stringPtrMA(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func timePtrMA(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, _ := time.Parse("2006-01-02", s)
	return &t
}

func (r *merchantAwardCommandRepository) Create(ctx context.Context, req *requests.CreateMerchantCertificationOrAwardRequest) (*models.MerchantCertificationsAndAward, error) {
	id, err := database.GetNextID(ctx, r.db, "merchant_certifications_and_awards")
	if err != nil {
		return nil, merchantaward_errors.ErrCreateMerchantAward.WithInternal(err)
	}
	now := time.Now()
	item := &models.MerchantCertificationsAndAward{
		MerchantCertificationID: id,
		MerchantID:              int32(req.MerchantID),
		Title:                   req.Title,
		Description:             stringPtrMA(req.Description),
		IssuedBy:                stringPtrMA(req.IssuedBy),
		IssueDate:               timePtrMA(req.IssueDate),
		ExpiryDate:              timePtrMA(req.ExpiryDate),
		CertificateUrl:          stringPtrMA(req.CertificateUrl),
		CreatedAt:               &now,
		UpdatedAt:               &now,
	}
	if _, err := r.db.Collection("merchant_certifications_and_awards").InsertOne(ctx, item); err != nil {
		return nil, merchantaward_errors.ErrCreateMerchantAward.WithInternal(err)
	}
	return item, nil
}

func (r *merchantAwardCommandRepository) Update(ctx context.Context, req *requests.UpdateMerchantCertificationOrAwardRequest) (*models.MerchantCertificationsAndAward, error) {
	coll := r.db.Collection("merchant_certifications_and_awards")
	var item models.MerchantCertificationsAndAward
	if err := coll.FindOne(ctx, bson.M{"_id": *req.MerchantCertificationID, "deleted_at": bson.M{"$eq": nil}}).Decode(&item); err != nil {
		return nil, merchantaward_errors.ErrMerchantAwardNotFound
	}
	now := time.Now()
	update := bson.M{
		"title":           req.Title,
		"description":     stringPtrMA(req.Description),
		"issued_by":       stringPtrMA(req.IssuedBy),
		"issue_date":      timePtrMA(req.IssueDate),
		"expiry_date":     timePtrMA(req.ExpiryDate),
		"certificate_url": stringPtrMA(req.CertificateUrl),
		"updated_at":      now,
	}
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": *req.MerchantCertificationID}, bson.M{"$set": update}); err != nil {
		return nil, merchantaward_errors.ErrUpdateMerchantAward.WithInternal(err)
	}
	if err := coll.FindOne(ctx, bson.M{"_id": *req.MerchantCertificationID}).Decode(&item); err != nil {
		return nil, merchantaward_errors.ErrUpdateMerchantAward.WithInternal(err)
	}
	return &item, nil
}

func (r *merchantAwardCommandRepository) Trash(ctx context.Context, id int) (*models.MerchantCertificationsAndAward, error) {
	coll := r.db.Collection("merchant_certifications_and_awards")
	var item models.MerchantCertificationsAndAward
	if err := coll.FindOne(ctx, bson.M{"_id": id, "deleted_at": bson.M{"$eq": nil}}).Decode(&item); err != nil {
		return nil, merchantaward_errors.ErrMerchantAwardNotFound
	}
	now := time.Now()
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"deleted_at": now}}); err != nil {
		return nil, merchantaward_errors.ErrTrashedMerchantAward.WithInternal(err)
	}
	item.DeletedAt = &now
	return &item, nil
}

func (r *merchantAwardCommandRepository) Restore(ctx context.Context, id int) (*models.MerchantCertificationsAndAward, error) {
	coll := r.db.Collection("merchant_certifications_and_awards")
	res, err := coll.UpdateOne(ctx, bson.M{"_id": id, "deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return nil, merchantaward_errors.ErrRestoreMerchantAward.WithInternal(err)
	}
	if res.MatchedCount == 0 {
		return nil, merchantaward_errors.ErrMerchantAwardNotFound
	}
	var item models.MerchantCertificationsAndAward
	if err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(&item); err != nil {
		return nil, merchantaward_errors.ErrRestoreMerchantAward.WithInternal(err)
	}
	return &item, nil
}

func (r *merchantAwardCommandRepository) DeletePermanent(ctx context.Context, id int) (bool, error) {
	res, err := r.db.Collection("merchant_certifications_and_awards").DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, merchantaward_errors.ErrDeleteMerchantAwardPermanent.WithInternal(err)
	}
	if res.DeletedCount == 0 {
		return false, merchantaward_errors.ErrMerchantAwardNotFound
	}
	return true, nil
}

func (r *merchantAwardCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("merchant_certifications_and_awards").UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, merchantaward_errors.ErrRestoreAllMerchantAwards.WithInternal(err)
	}
	return true, nil
}

func (r *merchantAwardCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("merchant_certifications_and_awards").DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, merchantaward_errors.ErrDeleteAllMerchantAwardsPermanent.WithInternal(err)
	}
	return true, nil
}
