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

type bannerseeder struct {
	db     *mongo.Database
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewBannerSeeder(db *mongo.Database, ctx context.Context, logger logger.LoggerInterface) *bannerseeder {
	return &bannerseeder{db: db, ctx: ctx, logger: logger}
}

func (r *bannerseeder) Seed() error {
	count, err := r.db.Collection("banners").CountDocuments(r.ctx, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		r.logger.Debug("banners already seeded, skipping")
		return nil
	}

	id, err := database.GetNextID(r.ctx, r.db, "banners")
	if err != nil {
		r.logger.Error("failed to get next banner id", zap.Error(err))
		return err
	}
	now := time.Now()
	active := true
	name := "Promo Banner"
	banner := models.Banner{
		BannerID:  id,
		Name:      name,
		IsActive:  &active,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	if _, err := r.db.Collection("banners").InsertOne(r.ctx, banner); err != nil {
		r.logger.Error("failed to seed banners", zap.Error(err))
		return err
	}

	r.logger.Info("banners seeded successfully")
	return nil
}
