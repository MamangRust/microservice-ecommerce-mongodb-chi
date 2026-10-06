package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/banner_errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type bannerQueryRepository struct {
	db *mongo.Database
}

func NewBannerQueryRepository(db *mongo.Database) *bannerQueryRepository {
	return &bannerQueryRepository{db: db}
}

func (r *bannerQueryRepository) FindAll(ctx context.Context, req *requests.FindAllBanner) ([]*BannerResult, error) {
	offset := (req.Page - 1) * req.PageSize
	filter := bson.M{"deleted_at": bson.M{"$eq": nil}}
	if req.Search != "" {
		filter["name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}
	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": bson.M{"banner_id": -1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data":     []bson.M{{"$skip": offset}, {"$limit": req.PageSize}},
		}},
	}
	results, totalCount, err := runPaginatedBannerQuery(ctx, r.db.Collection("banners"), pipeline)
	if err != nil {
		return nil, banner_errors.ErrFindAllBanners.WithInternal(err)
	}
	for _, item := range results {
		item.TotalCount = totalCount
	}
	return results, nil
}

func (r *bannerQueryRepository) FindActive(ctx context.Context, req *requests.FindAllBanner) ([]*BannerResult, error) {
	offset := (req.Page - 1) * req.PageSize
	filter := bson.M{"deleted_at": bson.M{"$eq": nil}, "is_active": true}
	if req.Search != "" {
		filter["name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}
	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": bson.M{"banner_id": -1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data":     []bson.M{{"$skip": offset}, {"$limit": req.PageSize}},
		}},
	}
	results, totalCount, err := runPaginatedBannerQuery(ctx, r.db.Collection("banners"), pipeline)
	if err != nil {
		return nil, banner_errors.ErrFindActiveBanners.WithInternal(err)
	}
	for _, item := range results {
		item.TotalCount = totalCount
	}
	return results, nil
}

func (r *bannerQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllBanner) ([]*BannerResult, error) {
	offset := (req.Page - 1) * req.PageSize
	filter := bson.M{"deleted_at": bson.M{"$ne": nil}}
	if req.Search != "" {
		filter["name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}
	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": bson.M{"banner_id": -1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data":     []bson.M{{"$skip": offset}, {"$limit": req.PageSize}},
		}},
	}
	results, totalCount, err := runPaginatedBannerQuery(ctx, r.db.Collection("banners"), pipeline)
	if err != nil {
		return nil, banner_errors.ErrFindTrashedBanners.WithInternal(err)
	}
	for _, item := range results {
		item.TotalCount = totalCount
	}
	return results, nil
}

func (r *bannerQueryRepository) FindByID(ctx context.Context, bannerID int) (*BannerResult, error) {
	var d bson.M
	err := r.db.Collection("banners").FindOne(ctx, bson.M{"_id": bannerID, "deleted_at": bson.M{"$eq": nil}}).Decode(&d)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, banner_errors.ErrBannerNotFound
		}
		return nil, banner_errors.ErrFindAllBanners.WithInternal(err)
	}
	return bsonMToBannerResult(d), nil
}

func runPaginatedBannerQuery(ctx context.Context, coll *mongo.Collection, pipeline []bson.M) ([]*BannerResult, int64, error) {
	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var facetResults []struct {
		Metadata []struct{ TotalCount int64 `bson:"total_count"` } `bson:"metadata"`
		Data     []bson.M `bson:"data"`
	}
	if err := cursor.All(ctx, &facetResults); err != nil {
		return nil, 0, err
	}
	if len(facetResults) == 0 {
		return []*BannerResult{}, 0, nil
	}

	var totalCount int64
	if len(facetResults[0].Metadata) > 0 {
		totalCount = facetResults[0].Metadata[0].TotalCount
	}

	var results []*BannerResult
	for _, d := range facetResults[0].Data {
		results = append(results, bsonMToBannerResult(d))
	}
	return results, totalCount, nil
}

func bsonMToBannerResult(m bson.M) *BannerResult {
	return &BannerResult{
		BannerID:   bsonMToInt32Banner(m, "_id"),
		Name:       bsonMToStrBanner(m, "name"),
		StartDate:  bsonMToStrPtrBanner(m, "start_date"),
		EndDate:    bsonMToStrPtrBanner(m, "end_date"),
		StartTime:  bsonMToStrPtrBanner(m, "start_time"),
		EndTime:    bsonMToStrPtrBanner(m, "end_time"),
		IsActive:   bsonMToBoolPtrBanner(m, "is_active"),
		CreatedAt:  bsonMToStrPtrBanner(m, "created_at"),
		UpdatedAt:  bsonMToStrPtrBanner(m, "updated_at"),
		DeletedAt:  bsonMToStrPtrBanner(m, "deleted_at"),
	}
}

func bsonMToInt32Banner(m bson.M, key string) int32 {
	v, ok := m[key]
	if !ok {
		return 0
	}
	switch t := v.(type) {
	case int32:
		return t
	case int64:
		return int32(t)
	case float64:
		return int32(t)
	default:
		return 0
	}
}

func bsonMToStrBanner(m bson.M, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func bsonMToStrPtrBanner(m bson.M, key string) *string {
	v, ok := m[key]
	if !ok || v == nil {
		return nil
	}
	switch t := v.(type) {
	case string:
		return &t
	case time.Time:
		s := t.Format("2006-01-02 15:04:05.000")
		return &s
	case bson.DateTime:
		s := t.Time().Format("2006-01-02 15:04:05.000")
		return &s
	default:
		return nil
	}
}

func bsonMToBoolPtrBanner(m bson.M, key string) *bool {
	v, ok := m[key]
	if !ok || v == nil {
		return nil
	}
	if b, ok := v.(bool); ok {
		return &b
	}
	return nil
}

func timePtrBanner(t time.Time) *time.Time { return &t }
