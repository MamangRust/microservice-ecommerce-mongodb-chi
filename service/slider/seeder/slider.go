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

type sliderseeder struct {
	db     *mongo.Database
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewSliderSeeder(db *mongo.Database, ctx context.Context, logger logger.LoggerInterface) *sliderseeder {
	return &sliderseeder{db: db, ctx: ctx, logger: logger}
}

func (r *sliderseeder) Seed() error {
	count, err := r.db.Collection("sliders").CountDocuments(r.ctx, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		r.logger.Debug("sliders already seeded, skipping")
		return nil
	}

	id, err := database.GetNextID(r.ctx, r.db, "sliders")
	if err != nil {
		r.logger.Error("failed to get next slider id", zap.Error(err))
		return err
	}
	now := time.Now()
	slider := models.Slider{
		SliderID:  id,
		Name:      "Promo Slider",
		Image:     "slider.png",
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	if _, err := r.db.Collection("sliders").InsertOne(r.ctx, slider); err != nil {
		r.logger.Error("failed to seed sliders", zap.Error(err))
		return err
	}

	r.logger.Info("sliders seeded successfully")
	return nil
}
