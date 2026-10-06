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

type merchantawardseeder struct {
	db     *mongo.Database
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewMerchantAwardSeeder(db *mongo.Database, ctx context.Context, logger logger.LoggerInterface) *merchantawardseeder {
	return &merchantawardseeder{db: db, ctx: ctx, logger: logger}
}

func (r *merchantawardseeder) Seed() error {
	count, err := r.db.Collection("merchant_certifications_and_awards").CountDocuments(r.ctx, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		r.logger.Debug("merchant_awards already seeded, skipping")
		return nil
	}

	id, err := database.GetNextID(r.ctx, r.db, "merchant_certifications_and_awards")
	if err != nil {
		r.logger.Error("failed to get next merchant_award id", zap.Error(err))
		return err
	}
	now := time.Now()
	desc := "ISO certification"
	issuedBy := "ISO"
	certURL := "https://example.com/cert.png"
	item := models.MerchantCertificationsAndAward{
		MerchantCertificationID: id,
		MerchantID:              0,
		Title:                   "ISO 9001",
		Description:             &desc,
		IssuedBy:                &issuedBy,
		IssueDate:               &now,
		CertificateUrl:          &certURL,
		CreatedAt:               &now,
		UpdatedAt:               &now,
	}

	if _, err := r.db.Collection("merchant_certifications_and_awards").InsertOne(r.ctx, item); err != nil {
		r.logger.Error("failed to seed merchant_awards", zap.Error(err))
		return err
	}

	r.logger.Info("merchant_awards seeded successfully")
	return nil
}
