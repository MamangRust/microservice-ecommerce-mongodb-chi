package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/banner_errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type bannerCommandRepository struct {
	db *mongo.Database
}

func NewBannerCommandRepository(db *mongo.Database) *bannerCommandRepository {
	return &bannerCommandRepository{db: db}
}

func parseDateBannerCmd(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func parseTimeBannerCmd(s string) (*string, error) {
	if s == "" {
		return nil, nil
	}
	_, err := time.Parse("15:04:05", s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *bannerCommandRepository) Create(ctx context.Context, request *requests.CreateBannerRequest) (*models.Banner, error) {
	startDate, _ := parseDateBannerCmd(request.StartDate)
	endDate, _ := parseDateBannerCmd(request.EndDate)
	startTime, _ := parseTimeBannerCmd(request.StartTime)
	endTime, _ := parseTimeBannerCmd(request.EndTime)

	id, err := database.GetNextID(ctx, r.db, "banners")
	if err != nil {
		return nil, banner_errors.ErrCreateBanner.WithInternal(err)
	}
	now := time.Now()
	banner := &models.Banner{
		BannerID:  id, Name: request.Name, StartDate: startDate, EndDate: endDate,
		StartTime: startTime, EndTime: endTime, IsActive: &request.IsActive,
		CreatedAt: &now, UpdatedAt: &now,
	}
	_, err = r.db.Collection("banners").InsertOne(ctx, banner)
	if err != nil {
		return nil, banner_errors.ErrCreateBanner.WithInternal(err)
	}
	return banner, nil
}

func (r *bannerCommandRepository) Update(ctx context.Context, request *requests.UpdateBannerRequest) (*models.Banner, error) {
	coll := r.db.Collection("banners")
	var banner models.Banner
	if err := coll.FindOne(ctx, bson.M{"_id": *request.BannerID}).Decode(&banner); err != nil {
		return nil, banner_errors.ErrBannerNotFound
	}
	startDate, _ := parseDateBannerCmd(request.StartDate)
	endDate, _ := parseDateBannerCmd(request.EndDate)
	startTime, _ := parseTimeBannerCmd(request.StartTime)
	endTime, _ := parseTimeBannerCmd(request.EndTime)
	banner.Name = request.Name
	banner.StartDate = startDate
	banner.EndDate = endDate
	banner.StartTime = startTime
	banner.EndTime = endTime
	banner.IsActive = &request.IsActive
	now := time.Now()
	banner.UpdatedAt = &now
	if _, err := coll.ReplaceOne(ctx, bson.M{"_id": *request.BannerID}, banner); err != nil {
		return nil, banner_errors.ErrUpdateBanner.WithInternal(err)
	}
	return &banner, nil
}

func (r *bannerCommandRepository) Trash(ctx context.Context, bannerID int) (*models.Banner, error) {
	coll := r.db.Collection("banners")
	var banner models.Banner
	if err := coll.FindOne(ctx, bson.M{"_id": bannerID}).Decode(&banner); err != nil {
		return nil, banner_errors.ErrBannerNotFound
	}
	now := time.Now()
	coll.UpdateOne(ctx, bson.M{"_id": bannerID}, bson.M{"$set": bson.M{"deleted_at": now}})
	banner.DeletedAt = &now
	return &banner, nil
}

func (r *bannerCommandRepository) Restore(ctx context.Context, bannerID int) (*models.Banner, error) {
	coll := r.db.Collection("banners")
	res, _ := coll.UpdateOne(ctx, bson.M{"_id": bannerID, "deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if res.MatchedCount == 0 {
		return nil, banner_errors.ErrBannerNotFound
	}
	var banner models.Banner
	coll.FindOne(ctx, bson.M{"_id": bannerID}).Decode(&banner)
	return &banner, nil
}

func (r *bannerCommandRepository) DeletePermanent(ctx context.Context, bannerID int) (bool, error) {
	res, err := r.db.Collection("banners").DeleteOne(ctx, bson.M{"_id": bannerID})
	if err != nil {
		return false, banner_errors.ErrDeleteBannerPermanent.WithInternal(err)
	}
	return res.DeletedCount > 0, nil
}

func (r *bannerCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("banners").UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, banner_errors.ErrRestoreAllBanners.WithInternal(err)
	}
	return true, nil
}

func (r *bannerCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("banners").DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, banner_errors.ErrDeleteAllBanners.WithInternal(err)
	}
	return true, nil
}
