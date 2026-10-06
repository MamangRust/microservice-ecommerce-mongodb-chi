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

type categoryseeder struct {
	db     *mongo.Database
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewCategorySeeder(db *mongo.Database, ctx context.Context, logger logger.LoggerInterface) *categoryseeder {
	return &categoryseeder{db: db, ctx: ctx, logger: logger}
}

func (r *categoryseeder) Seed() error {
	count, err := r.db.Collection("categories").CountDocuments(r.ctx, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		r.logger.Debug("categories already seeded, skipping")
		return nil
	}

	now := time.Now()
	slug := "electronics"
	desc := "Electronic products"
	img := "electronics.png"
	category := models.Category{
		Name:          "Electronics",
		Description:   &desc,
		SlugCategory:  &slug,
		ImageCategory: &img,
		CreatedAt:     &now,
		UpdatedAt:     &now,
	}

	id, err := database.GetNextID(r.ctx, r.db, "categories")
	if err != nil {
		r.logger.Error("failed to get next category id", zap.Error(err))
		return err
	}
	category.CategoryID = id

	if _, err := r.db.Collection("categories").InsertOne(r.ctx, category); err != nil {
		r.logger.Error("failed to seed categories", zap.Error(err))
		return err
	}

	r.logger.Info("categories seeded successfully")
	return nil
}
