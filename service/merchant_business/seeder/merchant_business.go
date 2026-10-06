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

type merchantbusinessseeder struct {
	db     *mongo.Database
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewMerchantBusinessSeeder(db *mongo.Database, ctx context.Context, logger logger.LoggerInterface) *merchantbusinessseeder {
	return &merchantbusinessseeder{db: db, ctx: ctx, logger: logger}
}

func (r *merchantbusinessseeder) Seed() error {
	count, err := r.db.Collection("merchant_business_information").CountDocuments(r.ctx, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		r.logger.Debug("merchant_businesses already seeded, skipping")
		return nil
	}

	id, err := database.GetNextID(r.ctx, r.db, "merchant_business_information")
	if err != nil {
		r.logger.Error("failed to get next merchant_business id", zap.Error(err))
		return err
	}
	now := time.Now()
	businessType := "retail"
	taxID := "123456789"
	year := int32(2020)
	employees := int32(25)
	website := "https://example.com"
	item := models.MerchantBusinessInformation{
		MerchantBusinessInfoID: id,
		MerchantID:             0,
		BusinessType:           &businessType,
		TaxID:                  &taxID,
		EstablishedYear:        &year,
		NumberOfEmployees:      &employees,
		WebsiteUrl:             &website,
		CreatedAt:              &now,
		UpdatedAt:              &now,
	}

	if _, err := r.db.Collection("merchant_business_information").InsertOne(r.ctx, item); err != nil {
		r.logger.Error("failed to seed merchant_businesses", zap.Error(err))
		return err
	}

	r.logger.Info("merchant_businesses seeded successfully")
	return nil
}
