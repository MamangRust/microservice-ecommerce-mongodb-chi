package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	merchantbusiness_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant_business"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type merchantBusinessCommandRepository struct {
	db *mongo.Database
}

func NewMerchantBusinessCommandRepository(db *mongo.Database) MerchantBusinessCommandRepository {
	return &merchantBusinessCommandRepository{db: db}
}

func int32PtrMB(v int) *int32 {
	i := int32(v)
	return &i
}

func (r *merchantBusinessCommandRepository) Create(ctx context.Context, req *requests.CreateMerchantBusinessInformationRequest) (*models.MerchantBusinessInformation, error) {
	id, err := database.GetNextID(ctx, r.db, "merchant_business_information")
	if err != nil {
		return nil, merchantbusiness_errors.ErrCreateMerchantBusiness.WithInternal(err)
	}
	now := time.Now()
	item := &models.MerchantBusinessInformation{
		MerchantBusinessInfoID: id,
		MerchantID:             int32(req.MerchantID),
		BusinessType:           &req.BusinessType,
		TaxID:                  &req.TaxID,
		EstablishedYear:        int32PtrMB(req.EstablishedYear),
		NumberOfEmployees:      int32PtrMB(req.NumberOfEmployees),
		WebsiteUrl:             &req.WebsiteUrl,
		CreatedAt:              &now,
		UpdatedAt:              &now,
	}
	if _, err := r.db.Collection("merchant_business_information").InsertOne(ctx, item); err != nil {
		return nil, merchantbusiness_errors.ErrCreateMerchantBusiness.WithInternal(err)
	}
	return item, nil
}

func (r *merchantBusinessCommandRepository) Update(ctx context.Context, req *requests.UpdateMerchantBusinessInformationRequest) (*models.MerchantBusinessInformation, error) {
	coll := r.db.Collection("merchant_business_information")
	var item models.MerchantBusinessInformation
	if err := coll.FindOne(ctx, bson.M{"_id": *req.MerchantBusinessInfoID, "deleted_at": bson.M{"$eq": nil}}).Decode(&item); err != nil {
		return nil, merchantbusiness_errors.ErrMerchantBusinessNotFound
	}
	now := time.Now()
	update := bson.M{
		"business_type":       &req.BusinessType,
		"tax_id":              &req.TaxID,
		"established_year":    int32PtrMB(req.EstablishedYear),
		"number_of_employees": int32PtrMB(req.NumberOfEmployees),
		"website_url":         &req.WebsiteUrl,
		"updated_at":          now,
	}
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": *req.MerchantBusinessInfoID}, bson.M{"$set": update}); err != nil {
		return nil, merchantbusiness_errors.ErrUpdateMerchantBusiness.WithInternal(err)
	}
	if err := coll.FindOne(ctx, bson.M{"_id": *req.MerchantBusinessInfoID}).Decode(&item); err != nil {
		return nil, merchantbusiness_errors.ErrUpdateMerchantBusiness.WithInternal(err)
	}
	return &item, nil
}

func (r *merchantBusinessCommandRepository) Trash(ctx context.Context, id int) (*models.MerchantBusinessInformation, error) {
	coll := r.db.Collection("merchant_business_information")
	var item models.MerchantBusinessInformation
	if err := coll.FindOne(ctx, bson.M{"_id": id, "deleted_at": bson.M{"$eq": nil}}).Decode(&item); err != nil {
		return nil, merchantbusiness_errors.ErrMerchantBusinessNotFound
	}
	now := time.Now()
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"deleted_at": now}}); err != nil {
		return nil, merchantbusiness_errors.ErrTrashMerchantBusiness.WithInternal(err)
	}
	item.DeletedAt = &now
	return &item, nil
}

func (r *merchantBusinessCommandRepository) Restore(ctx context.Context, id int) (*models.MerchantBusinessInformation, error) {
	coll := r.db.Collection("merchant_business_information")
	res, err := coll.UpdateOne(ctx, bson.M{"_id": id, "deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return nil, merchantbusiness_errors.ErrRestoreMerchantBusiness.WithInternal(err)
	}
	if res.MatchedCount == 0 {
		return nil, merchantbusiness_errors.ErrMerchantBusinessNotFound
	}
	var item models.MerchantBusinessInformation
	if err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(&item); err != nil {
		return nil, merchantbusiness_errors.ErrRestoreMerchantBusiness.WithInternal(err)
	}
	return &item, nil
}

func (r *merchantBusinessCommandRepository) DeletePermanent(ctx context.Context, id int) (bool, error) {
	res, err := r.db.Collection("merchant_business_information").DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, merchantbusiness_errors.ErrDeletePermanentMerchantBusiness.WithInternal(err)
	}
	if res.DeletedCount == 0 {
		return false, merchantbusiness_errors.ErrMerchantBusinessNotFound
	}
	return true, nil
}

func (r *merchantBusinessCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("merchant_business_information").UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, merchantbusiness_errors.ErrRestoreAllMerchantBusinesses.WithInternal(err)
	}
	return true, nil
}

func (r *merchantBusinessCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("merchant_business_information").DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, merchantbusiness_errors.ErrDeleteAllPermanentMerchantBusinesses.WithInternal(err)
	}
	return true, nil
}
