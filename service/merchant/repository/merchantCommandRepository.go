package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	shared_errors "github.com/MamangRust/microservice-ecommerce-shared/errors"
	merchant_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type merchantCommandRepository struct {
	db *mongo.Database
}

func NewMerchantCommandRepository(db *mongo.Database) *merchantCommandRepository {
	return &merchantCommandRepository{db: db}
}

func (r *merchantCommandRepository) Create(ctx context.Context, request *requests.CreateMerchantRequest) (*models.Merchant, error) {
	return r.create(ctx, r.db, request)
}

func (r *merchantCommandRepository) CreateInTx(ctx context.Context, tx *mongo.Database, request *requests.CreateMerchantRequest) (*models.Merchant, error) {
	return r.create(ctx, tx, request)
}

func (r *merchantCommandRepository) create(ctx context.Context, db *mongo.Database, request *requests.CreateMerchantRequest) (*models.Merchant, error) {
	id, err := database.GetNextID(ctx, db, "merchants")
	if err != nil {
		return nil, merchant_errors.ErrCreateMerchant.WithInternal(err)
	}
	now := time.Now()
	merchant := &models.Merchant{
		MerchantID:   id,
		UserID:       int32(request.UserID),
		Name:         request.Name,
		Status:       "active",
		Description:  stringPtrMerchant(request.Description),
		Address:      stringPtrMerchant(request.Address),
		ContactEmail: stringPtrMerchant(request.ContactEmail),
		ContactPhone: stringPtrMerchant(request.ContactPhone),
		CreatedAt:    &now,
		UpdatedAt:    &now,
	}
	if _, err := db.Collection("merchants").InsertOne(ctx, merchant); err != nil {
		return nil, merchant_errors.ErrCreateMerchant.WithInternal(err)
	}
	return merchant, nil
}

func (r *merchantCommandRepository) Update(ctx context.Context, request *requests.UpdateMerchantRequest) (*models.Merchant, error) {
	coll := r.db.Collection("merchants")
	var merchant models.Merchant
	if err := coll.FindOne(ctx, bson.M{"_id": *request.MerchantID}).Decode(&merchant); err != nil {
		return nil, merchant_errors.ErrMerchantNotFound
	}

	now := time.Now()
	update := bson.M{
		"name":          request.Name,
		"description":   stringPtrMerchant(request.Description),
		"address":       stringPtrMerchant(request.Address),
		"contact_email": stringPtrMerchant(request.ContactEmail),
		"contact_phone": stringPtrMerchant(request.ContactPhone),
		"status":        request.Status,
		"updated_at":    now,
	}
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": *request.MerchantID}, bson.M{"$set": update}); err != nil {
		return nil, merchant_errors.ErrUpdateMerchant.WithInternal(err)
	}
	if err := coll.FindOne(ctx, bson.M{"_id": *request.MerchantID}).Decode(&merchant); err != nil {
		return nil, merchant_errors.ErrUpdateMerchant.WithInternal(err)
	}
	return &merchant, nil
}

func (r *merchantCommandRepository) Trash(ctx context.Context, merchantID int) (*models.Merchant, error) {
	coll := r.db.Collection("merchants")
	var merchant models.Merchant
	if err := coll.FindOne(ctx, bson.M{"_id": merchantID}).Decode(&merchant); err != nil {
		return nil, merchant_errors.ErrMerchantNotFound
	}
	now := time.Now()
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": merchantID}, bson.M{"$set": bson.M{"deleted_at": now}}); err != nil {
		return nil, merchant_errors.ErrTrashedMerchant.WithInternal(err)
	}
	merchant.DeletedAt = &now
	return &merchant, nil
}

func (r *merchantCommandRepository) Restore(ctx context.Context, merchantID int) (*models.Merchant, error) {
	coll := r.db.Collection("merchants")
	res, err := coll.UpdateOne(ctx, bson.M{"_id": merchantID, "deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return nil, merchant_errors.ErrRestoreMerchant.WithInternal(err)
	}
	if res.MatchedCount == 0 {
		return nil, merchant_errors.ErrMerchantNotFound
	}
	var merchant models.Merchant
	if err := coll.FindOne(ctx, bson.M{"_id": merchantID}).Decode(&merchant); err != nil {
		return nil, merchant_errors.ErrRestoreMerchant.WithInternal(err)
	}
	return &merchant, nil
}

func (r *merchantCommandRepository) DeletePermanent(ctx context.Context, merchantID int) (bool, error) {
	res, err := r.db.Collection("merchants").DeleteOne(ctx, bson.M{"_id": merchantID})
	if err != nil {
		return false, shared_errors.NewConflictError("cannot permanently delete merchant while related records exist").WithInternal(err)
	}
	if res.DeletedCount == 0 {
		return false, merchant_errors.ErrMerchantNotFound
	}
	return true, nil
}

func (r *merchantCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("merchants").UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, merchant_errors.ErrRestoreAllMerchants.WithInternal(err)
	}
	return true, nil
}

func (r *merchantCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("merchants").DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, shared_errors.NewConflictError("cannot permanently delete merchants while related records exist").WithInternal(err)
	}
	return true, nil
}

func (r *merchantCommandRepository) UpdateStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*models.Merchant, error) {
	coll := r.db.Collection("merchants")
	var merchant models.Merchant
	if err := coll.FindOne(ctx, bson.M{"_id": *request.MerchantID}).Decode(&merchant); err != nil {
		return nil, merchant_errors.ErrMerchantNotFound
	}
	now := time.Now()
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": *request.MerchantID}, bson.M{"$set": bson.M{"status": request.Status, "updated_at": now}}); err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	if err := coll.FindOne(ctx, bson.M{"_id": *request.MerchantID}).Decode(&merchant); err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	return &merchant, nil
}

func (r *merchantCommandRepository) UpdateStatusInTx(ctx context.Context, tx *mongo.Database, request *requests.UpdateMerchantStatusRequest) (*models.Merchant, error) {
	coll := tx.Collection("merchants")
	var merchant models.Merchant
	if err := coll.FindOne(ctx, bson.M{"_id": *request.MerchantID}).Decode(&merchant); err != nil {
		return nil, merchant_errors.ErrMerchantNotFound
	}
	now := time.Now()
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": *request.MerchantID}, bson.M{"$set": bson.M{"status": request.Status, "updated_at": now}}); err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	if err := coll.FindOne(ctx, bson.M{"_id": *request.MerchantID}).Decode(&merchant); err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	return &merchant, nil
}

func stringPtrMerchant(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
