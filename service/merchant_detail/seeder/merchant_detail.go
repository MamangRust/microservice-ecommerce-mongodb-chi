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

type merchantdetailseeder struct {
	db     *mongo.Database
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewMerchantDetailSeeder(db *mongo.Database, ctx context.Context, logger logger.LoggerInterface) *merchantdetailseeder {
	return &merchantdetailseeder{db: db, ctx: ctx, logger: logger}
}

func (r *merchantdetailseeder) Seed() error {
	count, err := r.db.Collection("merchant_details").CountDocuments(r.ctx, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		r.logger.Debug("merchant_details already seeded, skipping")
		return nil
	}

	id, err := database.GetNextID(r.ctx, r.db, "merchant_details")
	if err != nil {
		r.logger.Error("failed to get next merchant_detail id", zap.Error(err))
		return err
	}
	now := time.Now()
	displayName := "Electronics Store"
	cover := "cover.png"
	logo := "logo.png"
	shortDesc := "Best electronics store"
	website := "https://example.com"
	detail := models.MerchantDetail{
		MerchantDetailID: id,
		MerchantID:       0,
		DisplayName:      &displayName,
		CoverImageUrl:    &cover,
		LogoUrl:          &logo,
		ShortDescription: &shortDesc,
		WebsiteUrl:       &website,
		CreatedAt:        &now,
		UpdatedAt:        &now,
	}

	if _, err := r.db.Collection("merchant_details").InsertOne(r.ctx, detail); err != nil {
		r.logger.Error("failed to seed merchant_details", zap.Error(err))
		return err
	}

	r.logger.Info("merchant_details seeded successfully")
	return nil
}
