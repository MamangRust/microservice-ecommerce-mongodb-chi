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

type orderseeder struct {
	db     *mongo.Database
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewOrderSeeder(db *mongo.Database, ctx context.Context, logger logger.LoggerInterface) *orderseeder {
	return &orderseeder{db: db, ctx: ctx, logger: logger}
}

func (r *orderseeder) Seed() error {
	count, err := r.db.Collection("orders").CountDocuments(r.ctx, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		r.logger.Debug("orders already seeded, skipping")
		return nil
	}

	// Grab first merchant & user ids if present
	var merchantID, userID int32
	var mDoc, uDoc bson.M
	if err := r.db.Collection("merchants").FindOne(r.ctx, bson.M{}).Decode(&mDoc); err == nil {
		merchantID = bsonMToInt32Seed(mDoc, "_id")
	}
	if err := r.db.Collection("users").FindOne(r.ctx, bson.M{}).Decode(&uDoc); err == nil {
		userID = bsonMToInt32Seed(uDoc, "_id")
	}
	if merchantID == 0 || userID == 0 {
		r.logger.Debug("no merchants or users found, skipping order seeding")
		return nil
	}

	id, err := database.GetNextID(r.ctx, r.db, "orders")
	if err != nil {
		r.logger.Error("failed to get next order id", zap.Error(err))
		return err
	}
	now := time.Now()
	order := models.Order{
		OrderID:    id,
		UserID:     userID,
		MerchantID: merchantID,
		TotalPrice: 100000,
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}

	if _, err := r.db.Collection("orders").InsertOne(r.ctx, order); err != nil {
		r.logger.Error("failed to seed orders", zap.Error(err))
		return err
	}

	r.logger.Info("orders seeded successfully")
	return nil
}

func bsonMToInt32Seed(m bson.M, key string) int32 {
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
