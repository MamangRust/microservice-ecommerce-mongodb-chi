package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	merchantpolicy_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant_policy_errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type merchantPolicyCommandRepository struct {
	db *mongo.Database
}

func NewMerchantPolicyCommandRepository(db *mongo.Database) MerchantPoliciesCommandRepository {
	return &merchantPolicyCommandRepository{db: db}
}

func (r *merchantPolicyCommandRepository) Create(ctx context.Context, req *requests.CreateMerchantPolicyRequest) (*models.MerchantPolicy, error) {
	id, err := database.GetNextID(ctx, r.db, "merchant_policies")
	if err != nil {
		return nil, merchantpolicy_errors.ErrCreateMerchantPolicy.WithInternal(err)
	}
	now := time.Now()
	policy := &models.MerchantPolicy{
		MerchantPolicyID: id,
		MerchantID:       int32(req.MerchantID),
		PolicyType:       req.PolicyType,
		Title:            req.Title,
		Description:      req.Description,
		CreatedAt:        &now,
		UpdatedAt:        &now,
	}
	if _, err := r.db.Collection("merchant_policies").InsertOne(ctx, policy); err != nil {
		return nil, merchantpolicy_errors.ErrCreateMerchantPolicy.WithInternal(err)
	}
	return policy, nil
}

func (r *merchantPolicyCommandRepository) Update(ctx context.Context, req *requests.UpdateMerchantPolicyRequest) (*models.MerchantPolicy, error) {
	coll := r.db.Collection("merchant_policies")
	var policy models.MerchantPolicy
	if err := coll.FindOne(ctx, bson.M{"_id": *req.MerchantPolicyID, "deleted_at": bson.M{"$eq": nil}}).Decode(&policy); err != nil {
		return nil, merchantpolicy_errors.ErrMerchantPolicyNotFound
	}
	now := time.Now()
	update := bson.M{
		"policy_type": req.PolicyType,
		"title":       req.Title,
		"description": req.Description,
		"updated_at":  now,
	}
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": *req.MerchantPolicyID}, bson.M{"$set": update}); err != nil {
		return nil, merchantpolicy_errors.ErrUpdateMerchantPolicy.WithInternal(err)
	}
	if err := coll.FindOne(ctx, bson.M{"_id": *req.MerchantPolicyID}).Decode(&policy); err != nil {
		return nil, merchantpolicy_errors.ErrUpdateMerchantPolicy.WithInternal(err)
	}
	return &policy, nil
}

func (r *merchantPolicyCommandRepository) Trash(ctx context.Context, id int) (*models.MerchantPolicy, error) {
	coll := r.db.Collection("merchant_policies")
	var policy models.MerchantPolicy
	if err := coll.FindOne(ctx, bson.M{"_id": id, "deleted_at": bson.M{"$eq": nil}}).Decode(&policy); err != nil {
		return nil, merchantpolicy_errors.ErrMerchantPolicyNotFound
	}
	now := time.Now()
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"deleted_at": now}}); err != nil {
		return nil, merchantpolicy_errors.ErrTrashedMerchantPolicy.WithInternal(err)
	}
	policy.DeletedAt = &now
	return &policy, nil
}

func (r *merchantPolicyCommandRepository) Restore(ctx context.Context, id int) (*models.MerchantPolicy, error) {
	coll := r.db.Collection("merchant_policies")
	res, err := coll.UpdateOne(ctx, bson.M{"_id": id, "deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return nil, merchantpolicy_errors.ErrRestoreMerchantPolicy.WithInternal(err)
	}
	if res.MatchedCount == 0 {
		return nil, merchantpolicy_errors.ErrMerchantPolicyNotFound
	}
	var policy models.MerchantPolicy
	if err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(&policy); err != nil {
		return nil, merchantpolicy_errors.ErrRestoreMerchantPolicy.WithInternal(err)
	}
	return &policy, nil
}

func (r *merchantPolicyCommandRepository) DeletePermanent(ctx context.Context, id int) (bool, error) {
	res, err := r.db.Collection("merchant_policies").DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, merchantpolicy_errors.ErrDeleteMerchantPolicyPermanent.WithInternal(err)
	}
	if res.DeletedCount == 0 {
		return false, merchantpolicy_errors.ErrMerchantPolicyNotFound
	}
	return true, nil
}

func (r *merchantPolicyCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("merchant_policies").UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, merchantpolicy_errors.ErrRestoreAllMerchantPolicies.WithInternal(err)
	}
	return true, nil
}

func (r *merchantPolicyCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("merchant_policies").DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, merchantpolicy_errors.ErrDeleteAllMerchantPoliciesPermanent.WithInternal(err)
	}
	return true, nil
}
