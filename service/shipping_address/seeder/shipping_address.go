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

type shippingaddressseeder struct {
	db     *mongo.Database
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewShippingAddressSeeder(db *mongo.Database, ctx context.Context, logger logger.LoggerInterface) *shippingaddressseeder {
	return &shippingaddressseeder{db: db, ctx: ctx, logger: logger}
}

func (r *shippingaddressseeder) Seed() error {
	count, err := r.db.Collection("shipping_addresses").CountDocuments(r.ctx, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		r.logger.Debug("shipping_addresses already seeded, skipping")
		return nil
	}

	id, err := database.GetNextID(r.ctx, r.db, "shipping_addresses")
	if err != nil {
		r.logger.Error("failed to get next shipping_address id", zap.Error(err))
		return err
	}
	now := time.Now()
	addr := models.ShippingAddress{
		ShippingAddressID: id,
		OrderID:           1,
		Alamat:            "Jl. Merdeka No. 1",
		Provinsi:          "DKI Jakarta",
		Kota:              "Jakarta",
		Negara:            "Indonesia",
		Courier:           "JNE",
		ShippingMethod:    "regular",
		ShippingCost:      15000,
		CreatedAt:         &now,
		UpdatedAt:         &now,
	}

	if _, err := r.db.Collection("shipping_addresses").InsertOne(r.ctx, addr); err != nil {
		r.logger.Error("failed to seed shipping_addresses", zap.Error(err))
		return err
	}

	r.logger.Info("shipping_addresses seeded successfully")
	return nil
}
