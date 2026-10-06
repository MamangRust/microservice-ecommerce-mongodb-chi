package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/slider_errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type sliderQueryRepository struct {
	db *mongo.Database
}

func NewSliderQueryRepository(db *mongo.Database) *sliderQueryRepository {
	return &sliderQueryRepository{db: db}
}

func (r *sliderQueryRepository) FindAll(ctx context.Context, req *requests.FindAllSlider) ([]*SliderResult, error) {
	return r.runQuery(ctx, bson.M{"deleted_at": bson.M{"$eq": nil}}, req, slider_errors.ErrFindAllSliders)
}

func (r *sliderQueryRepository) FindActive(ctx context.Context, req *requests.FindAllSlider) ([]*SliderResult, error) {
	return r.runQuery(ctx, bson.M{"deleted_at": bson.M{"$eq": nil}}, req, slider_errors.ErrFindActiveSliders)
}

func (r *sliderQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllSlider) ([]*SliderResult, error) {
	return r.runQuery(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, req, slider_errors.ErrFindTrashedSliders)
}

func (r *sliderQueryRepository) FindByID(ctx context.Context, sliderID int) (*SliderResult, error) {
	var d bson.M
	err := r.db.Collection("sliders").FindOne(ctx, bson.M{"_id": sliderID, "deleted_at": bson.M{"$eq": nil}}).Decode(&d)
	if err != nil {
		return nil, slider_errors.ErrFindSliderByID
	}
	return bsonMToSliderResult(d), nil
}

func (r *sliderQueryRepository) runQuery(ctx context.Context, filter bson.M, req *requests.FindAllSlider, errFactory *errors.AppError) ([]*SliderResult, error) {
	offset := (req.Page - 1) * req.PageSize
	if req.Search != "" {
		filter["name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}
	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": bson.M{"slider_id": -1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data":     []bson.M{{"$skip": offset}, {"$limit": req.PageSize}},
		}},
	}
	cursor, err := r.db.Collection("sliders").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errFactory.WithInternal(err)
	}
	defer cursor.Close(ctx)
	var facetResults []struct {
		Metadata []struct {
			TotalCount int64 `bson:"total_count"`
		} `bson:"metadata"`
		Data []bson.M `bson:"data"`
	}
	if err := cursor.All(ctx, &facetResults); err != nil {
		return nil, errFactory.WithInternal(err)
	}
	if len(facetResults) == 0 {
		return []*SliderResult{}, nil
	}
	var results []*SliderResult
	for _, d := range facetResults[0].Data {
		item := bsonMToSliderResult(d)
		if len(facetResults[0].Metadata) > 0 {
			item.TotalCount = facetResults[0].Metadata[0].TotalCount
		}
		results = append(results, item)
	}
	return results, nil
}

func bsonMToSliderResult(m bson.M) *SliderResult {
	return &SliderResult{
		SliderID:  bsonMToInt32Sl(m, "_id"),
		Name:      bsonMToStrSl(m, "name"),
		Image:     bsonMToStrSl(m, "image"),
		CreatedAt: bsonMToTimeStrPtrSl(m, "created_at"),
		UpdatedAt: bsonMToTimeStrPtrSl(m, "updated_at"),
		DeletedAt: bsonMToTimeStrPtrSl(m, "deleted_at"),
	}
}

func bsonMToInt32Sl(m bson.M, key string) int32 {
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

func bsonMToStrSl(m bson.M, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func bsonMToTimeStrPtrSl(m bson.M, key string) *string {
	v, ok := m[key]
	if !ok || v == nil {
		return nil
	}
	switch t := v.(type) {
	case time.Time:
		s := t.Format("2006-01-02 15:04:05.000")
		return &s
	case bson.DateTime:
		s := t.Time().Format("2006-01-02 15:04:05.000")
		return &s
	case string:
		return &t
	default:
		return nil
	}
}
