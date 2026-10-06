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

type merchantseeder struct {
	db     *mongo.Database
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewMerchantSeeder(db *mongo.Database, ctx context.Context, logger logger.LoggerInterface) *merchantseeder {
	return &merchantseeder{db: db, ctx: ctx, logger: logger}
}

func (r *merchantseeder) Seed() error {
	count, err := r.db.Collection("merchants").CountDocuments(r.ctx, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		r.logger.Debug("merchants already seeded, skipping")
		return nil
	}

	// Grab first user id if present
	var userID int32
	var uDoc bson.M
	if err := r.db.Collection("users").FindOne(r.ctx, bson.M{}).Decode(&uDoc); err == nil {
		userID = bsonMToInt32Mer(uDoc, "_id")
	}

	id, err := database.GetNextID(r.ctx, r.db, "merchants")
	if err != nil {
		r.logger.Error("failed to get next merchant id", zap.Error(err))
		return err
	}
	now := time.Now()
	name := "Electronics Store"
	desc := "Best electronics"
	addr := "Jakarta"
	email := "merchant@example.com"
	phone := "08123456789"
	merchant := models.Merchant{
		MerchantID:   id,
		UserID:       userID,
		Name:         name,
		Description:  &desc,
		Address:      &addr,
		ContactEmail: &email,
		ContactPhone: &phone,
		Status:       "active",
		CreatedAt:    &now,
		UpdatedAt:    &now,
	}

	if _, err := r.db.Collection("merchants").InsertOne(r.ctx, merchant); err != nil {
		r.logger.Error("failed to seed merchants", zap.Error(err))
		return err
	}

	r.logger.Info("merchants seeded successfully")
	return nil
}

func bsonMToInt32Mer(m bson.M, key string) int32 {
	v, ok := m[key]
	if !ok {
		return 0
	}
	switch t := v.(type) {
	case int32:
		return t
	case int64:
		return int32(t)
	case float64:
		return int32(t)
	default:
		return 0
	}
}
