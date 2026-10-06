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

type transactionseeder struct {
	db     *mongo.Database
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewTransactionSeeder(db *mongo.Database, ctx context.Context, logger logger.LoggerInterface) *transactionseeder {
	return &transactionseeder{db: db, ctx: ctx, logger: logger}
}

func (r *transactionseeder) Seed() error {
	count, err := r.db.Collection("transactions").CountDocuments(r.ctx, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		r.logger.Debug("transactions already seeded, skipping")
		return nil
	}

	id, err := database.GetNextID(r.ctx, r.db, "transactions")
	if err != nil {
		r.logger.Error("failed to get next transaction id", zap.Error(err))
		return err
	}
	now := time.Now()
	txn := models.Transaction{
		TransactionID: id,
		OrderID:       1,
		MerchantID:    1,
		PaymentMethod: "bank_transfer",
		Amount:        100000,
		PaymentStatus: "paid",
		CreatedAt:     &now,
		UpdatedAt:     &now,
	}

	if _, err := r.db.Collection("transactions").InsertOne(r.ctx, txn); err != nil {
		r.logger.Error("failed to seed transactions", zap.Error(err))
		return err
	}

	r.logger.Info("transactions seeded successfully")
	return nil
}
