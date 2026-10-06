package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/slider_errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type sliderCommandRepository struct {
	db *mongo.Database
}

func NewSliderCommandRepository(db *mongo.Database) *sliderCommandRepository {
	return &sliderCommandRepository{db: db}
}

func (r *sliderCommandRepository) Create(ctx context.Context, request *requests.CreateSliderRequest) (*models.Slider, error) {
	id, err := database.GetNextID(ctx, r.db, "sliders")
	if err != nil {
		return nil, slider_errors.ErrCreateSlider
	}
	now := time.Now()
	slider := &models.Slider{
		SliderID:  id,
		Name:      request.Nama,
		Image:     request.FilePath,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
	if _, err := r.db.Collection("sliders").InsertOne(ctx, slider); err != nil {
		return nil, slider_errors.ErrCreateSlider
	}
	return slider, nil
}

func (r *sliderCommandRepository) Update(ctx context.Context, request *requests.UpdateSliderRequest) (*models.Slider, error) {
	coll := r.db.Collection("sliders")
	var slider models.Slider
	if err := coll.FindOne(ctx, bson.M{"_id": *request.ID}).Decode(&slider); err != nil {
		return nil, slider_errors.ErrSliderNotFound
	}
	now := time.Now()
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": *request.ID}, bson.M{"$set": bson.M{
		"name":       request.Nama,
		"image":      request.FilePath,
		"updated_at": now,
	}}); err != nil {
		return nil, slider_errors.ErrUpdateSlider
	}
	slider.Name = request.Nama
	slider.Image = request.FilePath
	slider.UpdatedAt = &now
	return &slider, nil
}

func (r *sliderCommandRepository) Trash(ctx context.Context, sliderID int) (*models.Slider, error) {
	coll := r.db.Collection("sliders")
	var slider models.Slider
	if err := coll.FindOne(ctx, bson.M{"_id": sliderID}).Decode(&slider); err != nil {
		return nil, slider_errors.ErrSliderNotFound
	}
	now := time.Now()
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": sliderID}, bson.M{"$set": bson.M{"deleted_at": now}}); err != nil {
		return nil, slider_errors.ErrTrashSlider
	}
	slider.DeletedAt = &now
	return &slider, nil
}

func (r *sliderCommandRepository) Restore(ctx context.Context, sliderID int) (*models.Slider, error) {
	coll := r.db.Collection("sliders")
	res, err := coll.UpdateOne(ctx, bson.M{"_id": sliderID, "deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return nil, slider_errors.ErrRestoreSlider
	}
	if res.MatchedCount == 0 {
		return nil, slider_errors.ErrSliderNotFound
	}
	var slider models.Slider
	if err := coll.FindOne(ctx, bson.M{"_id": sliderID}).Decode(&slider); err != nil {
		return nil, slider_errors.ErrRestoreSlider
	}
	return &slider, nil
}

func (r *sliderCommandRepository) DeletePermanent(ctx context.Context, sliderID int) (bool, error) {
	res, err := r.db.Collection("sliders").DeleteOne(ctx, bson.M{"_id": sliderID})
	if err != nil {
		return false, slider_errors.ErrDeletePermanentSlider
	}
	if res.DeletedCount == 0 {
		return false, slider_errors.ErrSliderNotFound
	}
	return true, nil
}

func (r *sliderCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("sliders").UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, slider_errors.ErrRestoreAllSlider
	}
	return true, nil
}

func (r *sliderCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("sliders").DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, slider_errors.ErrDeleteAllPermanentSlider
	}
	return true, nil
}
