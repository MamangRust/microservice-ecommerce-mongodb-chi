package seeder

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

type merchantpolicyseeder struct {
	db     *mongo.Database
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewMerchantPolicySeeder(db *mongo.Database, ctx context.Context, logger logger.LoggerInterface) *merchantpolicyseeder {
	return &merchantpolicyseeder{db: db, ctx: ctx, logger: logger}
}

func (r *merchantpolicyseeder) Seed() error {
	count, err := r.db.Collection("merchant_policies").CountDocuments(r.ctx, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		r.logger.Debug("merchant_policies already seeded, skipping")
		return nil
	}

	id, err := database.GetNextID(r.ctx, r.db, "merchant_policies")
	if err != nil {
		r.logger.Error("failed to get next merchant_policy id", zap.Error(err))
		return err
	}
	now := time.Now()
	policy := models.MerchantPolicy{
		MerchantPolicyID: id,
		MerchantID:       0,
		PolicyType:       "return",
		Title:            "Return Policy",
		Description:      "30-day return policy",
		CreatedAt:        &now,
		UpdatedAt:        &now,
	}

	if _, err := r.db.Collection("merchant_policies").InsertOne(r.ctx, policy); err != nil {
		r.logger.Error("failed to seed merchant_policies", zap.Error(err))
		return err
	}

	r.logger.Info("merchant_policies seeded successfully")
	return nil
}
