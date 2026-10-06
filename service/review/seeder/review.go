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

type reviewseeder struct {
	db     *mongo.Database
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewReviewSeeder(db *mongo.Database, ctx context.Context, logger logger.LoggerInterface) *reviewseeder {
	return &reviewseeder{db: db, ctx: ctx, logger: logger}
}

func (r *reviewseeder) Seed() error {
	count, err := r.db.Collection("reviews").CountDocuments(r.ctx, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		r.logger.Debug("reviews already seeded, skipping")
		return nil
	}

	now := time.Now()
	for i := 1; i <= 5; i++ {
		id, err := database.GetNextID(r.ctx, r.db, "reviews")
		if err != nil {
			r.logger.Error("failed to get next review id", zap.Error(err))
			return err
		}
		review := models.Review{
			ReviewID:  id,
			UserID:    1,
			ProductID: 1,
			Name:      "Customer",
			Comment:   "Great product!",
			Rating:    5,
			CreatedAt: &now,
			UpdatedAt: &now,
		}
		if _, err := r.db.Collection("reviews").InsertOne(r.ctx, review); err != nil {
			r.logger.Error("failed to seed review", zap.Error(err))
			return err
		}
	}

	r.logger.Info("reviews seeded successfully")
	return nil
}
