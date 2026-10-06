package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors"
	review_detail_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/review_detail"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type reviewDetailQueryRepository struct {
	db *mongo.Database
}

func NewReviewDetailQueryRepository(db *mongo.Database) *reviewDetailQueryRepository {
	return &reviewDetailQueryRepository{db: db}
}

func (r *reviewDetailQueryRepository) FindAll(ctx context.Context, req *requests.FindAllReview) ([]*ReviewDetailResult, error) {
	return r.runQuery(ctx, bson.M{"deleted_at": bson.M{"$eq": nil}}, req, review_detail_errors.ErrFindAllReviewDetails)
}

func (r *reviewDetailQueryRepository) FindActive(ctx context.Context, req *requests.FindAllReview) ([]*ReviewDetailResult, error) {
	return r.runQuery(ctx, bson.M{"deleted_at": bson.M{"$eq": nil}}, req, review_detail_errors.ErrFindActiveReviewDetails)
}

func (r *reviewDetailQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllReview) ([]*ReviewDetailResult, error) {
	return r.runQuery(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, req, review_detail_errors.ErrFindTrashedReviewDetails)
}

func (r *reviewDetailQueryRepository) FindByID(ctx context.Context, id int) (*ReviewDetailResult, error) {
	var d bson.M
	err := r.db.Collection("review_details").FindOne(ctx, bson.M{"_id": id}).Decode(&d)
	if err != nil {
		return nil, review_detail_errors.ErrReviewDetailNotFound
	}
	return bsonMToReviewDetailResult(d), nil
}

func (r *reviewDetailQueryRepository) runQuery(ctx context.Context, filter bson.M, req *requests.FindAllReview, errFactory *errors.AppError) ([]*ReviewDetailResult, error) {
	offset := (req.Page - 1) * req.PageSize
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"type": bson.M{"$regex": req.Search, "$options": "i"}},
			{"caption": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}
	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": bson.M{"review_detail_id": -1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data":     []bson.M{{"$skip": offset}, {"$limit": req.PageSize}},
		}},
	}
	cursor, err := r.db.Collection("review_details").Aggregate(ctx, pipeline)
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
		return []*ReviewDetailResult{}, nil
	}
	var results []*ReviewDetailResult
	for _, d := range facetResults[0].Data {
		item := bsonMToReviewDetailResult(d)
		if len(facetResults[0].Metadata) > 0 {
			item.TotalCount = facetResults[0].Metadata[0].TotalCount
		}
		results = append(results, item)
	}
	return results, nil
}

func bsonMToReviewDetailResult(m bson.M) *ReviewDetailResult {
	return &ReviewDetailResult{
		ReviewDetailID: bsonMToInt32RD(m, "_id"),
		ReviewID:       bsonMToInt32RD(m, "review_id"),
		Type:           bsonMToStrRD(m, "type"),
		Url:            bsonMToStrRD(m, "url"),
		Caption:        bsonMToStrPtrRD(m, "caption"),
		CreatedAt:      bsonMToTimeStrPtrRD(m, "created_at"),
		UpdatedAt:      bsonMToTimeStrPtrRD(m, "updated_at"),
		DeletedAt:      bsonMToTimeStrPtrRD(m, "deleted_at"),
	}
}

func bsonMToInt32RD(m bson.M, key string) int32 {
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

func bsonMToStrRD(m bson.M, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func bsonMToStrPtrRD(m bson.M, key string) *string {
	v, ok := m[key]
	if !ok || v == nil {
		return nil
	}
	if s, ok := v.(string); ok {
		return &s
	}
	return nil
}

func bsonMToTimeStrPtrRD(m bson.M, key string) *string {
	v, ok := m[key]
	if !ok || v == nil {
		return nil
	}
	switch t := v.(type) {
	case time.Time:
		s := t.Format("2006-01-02 15:04:05")
		return &s
	case bson.DateTime:
		s := t.Time().Format("2006-01-02 15:04:05")
		return &s
	case string:
		return &t
	default:
		return nil
	}
}
