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

type productseeder struct {
	db     *mongo.Database
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewProductSeeder(db *mongo.Database, ctx context.Context, logger logger.LoggerInterface) *productseeder {
	return &productseeder{db: db, ctx: ctx, logger: logger}
}

func (r *productseeder) Seed() error {
	count, err := r.db.Collection("products").CountDocuments(r.ctx, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		r.logger.Debug("products already seeded, skipping")
		return nil
	}

	now := time.Now()
	for i := 1; i <= 10; i++ {
		id, err := database.GetNextID(r.ctx, r.db, "products")
		if err != nil {
			r.logger.Error("failed to get next product id", zap.Error(err))
			return err
		}
		desc := "Sample product"
		brand := "Generic"
		weight := int32(500)
		slug := "sample-product"
		img := "product.png"
		product := models.Product{
			ProductID:    id,
			MerchantID:   1,
			CategoryID:   1,
			Name:         "Sample Product",
			Description:  &desc,
			Price:        100000,
			CountInStock: 50,
			Brand:        &brand,
			Weight:       &weight,
			SlugProduct:  &slug,
			ImageProduct: &img,
			CreatedAt:    &now,
			UpdatedAt:    &now,
		}
		if _, err := r.db.Collection("products").InsertOne(r.ctx, product); err != nil {
			r.logger.Error("failed to seed product", zap.Error(err))
			return err
		}
	}

	r.logger.Info("products seeded successfully")
	return nil
}
