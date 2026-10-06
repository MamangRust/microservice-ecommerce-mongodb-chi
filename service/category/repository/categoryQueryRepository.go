package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/category_errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type categoryQueryRepository struct {
	db *mongo.Database
}

func NewCategoryQueryRepository(db *mongo.Database) *categoryQueryRepository {
	return &categoryQueryRepository{db: db}
}

func (r *categoryQueryRepository) FindAll(ctx context.Context, req *requests.FindAllCategory) ([]*CategoryResult, error) {
	offset := (req.Page - 1) * req.PageSize
	filter := bson.M{"deleted_at": bson.M{"$eq": nil}}
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": req.Search, "$options": "i"}},
			{"description": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}

	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": bson.M{"category_id": -1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data": []bson.M{{"$skip": offset}, {"$limit": req.PageSize}},
		}},
	}
	cursor, err := r.db.Collection("categories").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, category_errors.ErrFindAllCategory.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var facetResults []struct {
		Metadata []struct{ TotalCount int64 `bson:"total_count"` } `bson:"metadata"`
		Data     []bson.M `bson:"data"`
	}
	if err := cursor.All(ctx, &facetResults); err != nil {
		return nil, category_errors.ErrFindAllCategory.WithInternal(err)
	}
	if len(facetResults) == 0 {
		return []*CategoryResult{}, nil
	}

	var results []*CategoryResult
	for _, d := range facetResults[0].Data {
		item := &CategoryResult{
			CategoryID:    bsonMToInt32Cat(d, "_id"),
			Name:          bsonMToStrCat(d, "name"),
			Description:   bsonMToStrPtrCat(d, "description"),
			SlugCategory:  bsonMToStrPtrCat(d, "slug_category"),
			ImageCategory: bsonMToStrPtrCat(d, "image_category"),
			CreatedAt:     bsonMToTimePtrCat(d, "created_at"),
			UpdatedAt:     bsonMToTimePtrCat(d, "updated_at"),
			DeletedAt:     bsonMToTimePtrCat(d, "deleted_at"),
		}
		if len(facetResults[0].Metadata) > 0 {
			item.TotalCount = facetResults[0].Metadata[0].TotalCount
		}
		results = append(results, item)
	}
	return results, nil
}

func (r *categoryQueryRepository) FindByID(ctx context.Context, categoryID int) (*models.Category, error) {
	var result models.Category
	err := r.db.Collection("categories").FindOne(ctx, bson.M{"_id": categoryID, "deleted_at": bson.M{"$eq": nil}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, category_errors.ErrCategoryNotFound
		}
		return nil, category_errors.ErrFindCategoryById.WithInternal(err)
	}
	return &result, nil
}

func (r *categoryQueryRepository) FindByIDTrashed(ctx context.Context, categoryID int) (*models.Category, error) {
	var result models.Category
	err := r.db.Collection("categories").FindOne(ctx, bson.M{"_id": categoryID, "deleted_at": bson.M{"$ne": nil}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, category_errors.ErrCategoryNotFound
		}
		return nil, category_errors.ErrFindCategoryByIdTrashed.WithInternal(err)
	}
	return &result, nil
}

func (r *categoryQueryRepository) FindActive(ctx context.Context, req *requests.FindAllCategory) ([]*CategoryResult, error) {
	offset := (req.Page - 1) * req.PageSize
	filter := bson.M{"deleted_at": bson.M{"$eq": nil}}
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": req.Search, "$options": "i"}},
			{"description": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}

	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": bson.M{"category_id": -1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data": []bson.M{{"$skip": offset}, {"$limit": req.PageSize}},
		}},
	}
	cursor, err := r.db.Collection("categories").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, category_errors.ErrFindByActiveCategory.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var facetResults []struct {
		Metadata []struct{ TotalCount int64 `bson:"total_count"` } `bson:"metadata"`
		Data     []bson.M `bson:"data"`
	}
	if err := cursor.All(ctx, &facetResults); err != nil {
		return nil, category_errors.ErrFindByActiveCategory.WithInternal(err)
	}
	if len(facetResults) == 0 {
		return []*CategoryResult{}, nil
	}

	var results []*CategoryResult
	for _, d := range facetResults[0].Data {
		item := &CategoryResult{
			CategoryID:    bsonMToInt32Cat(d, "_id"),
			Name:          bsonMToStrCat(d, "name"),
			Description:   bsonMToStrPtrCat(d, "description"),
			SlugCategory:  bsonMToStrPtrCat(d, "slug_category"),
			ImageCategory: bsonMToStrPtrCat(d, "image_category"),
			CreatedAt:     bsonMToTimePtrCat(d, "created_at"),
			UpdatedAt:     bsonMToTimePtrCat(d, "updated_at"),
			DeletedAt:     bsonMToTimePtrCat(d, "deleted_at"),
		}
		if len(facetResults[0].Metadata) > 0 {
			item.TotalCount = facetResults[0].Metadata[0].TotalCount
		}
		results = append(results, item)
	}
	return results, nil
}

func (r *categoryQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllCategory) ([]*CategoryResult, error) {
	offset := (req.Page - 1) * req.PageSize
	filter := bson.M{"deleted_at": bson.M{"$ne": nil}}
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}

	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": bson.M{"category_id": -1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data": []bson.M{{"$skip": offset}, {"$limit": req.PageSize}},
		}},
	}
	cursor, err := r.db.Collection("categories").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, category_errors.ErrFindByTrashedCategory.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var facetResults []struct {
		Metadata []struct{ TotalCount int64 `bson:"total_count"` } `bson:"metadata"`
		Data     []bson.M `bson:"data"`
	}
	if err := cursor.All(ctx, &facetResults); err != nil {
		return nil, category_errors.ErrFindByTrashedCategory.WithInternal(err)
	}
	if len(facetResults) == 0 {
		return []*CategoryResult{}, nil
	}

	var results []*CategoryResult
	for _, d := range facetResults[0].Data {
		item := &CategoryResult{
			CategoryID:    bsonMToInt32Cat(d, "_id"),
			Name:          bsonMToStrCat(d, "name"),
			Description:   bsonMToStrPtrCat(d, "description"),
			SlugCategory:  bsonMToStrPtrCat(d, "slug_category"),
			ImageCategory: bsonMToStrPtrCat(d, "image_category"),
			CreatedAt:     bsonMToTimePtrCat(d, "created_at"),
			UpdatedAt:     bsonMToTimePtrCat(d, "updated_at"),
			DeletedAt:     bsonMToTimePtrCat(d, "deleted_at"),
		}
		if len(facetResults[0].Metadata) > 0 {
			item.TotalCount = facetResults[0].Metadata[0].TotalCount
		}
		results = append(results, item)
	}
	return results, nil
}

// helpers for category package
func bsonMToTimePtrCat(m bson.M, key string) *time.Time {
	v, ok := m[key]
	if !ok || v == nil {
		return nil
	}
	switch t := v.(type) {
	case time.Time:
		return &t
	case bson.DateTime:
		ts := t.Time()
		return &ts
	default:
		return nil
	}
}
func bsonMToInt32Cat(m bson.M, key string) int32 {
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

func bsonMToStrCat(m bson.M, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func bsonMToStrPtrCat(m bson.M, key string) *string {
	v, ok := m[key]
	if !ok || v == nil {
		return nil
	}
	if s, ok := v.(string); ok {
		return &s
	}
	return nil
}
