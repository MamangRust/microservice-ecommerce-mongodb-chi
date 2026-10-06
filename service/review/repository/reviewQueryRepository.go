package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors"
	review_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/review"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type reviewQueryRepository struct {
	db        *mongo.Database // ec_experience: reviews, review_details
	catalogDB *mongo.Database // ec_catalog: products
}

func NewReviewQueryRepository(db, catalogDB *mongo.Database) *reviewQueryRepository {
	return &reviewQueryRepository{db: db, catalogDB: catalogDB}
}

func (r *reviewQueryRepository) FindAll(ctx context.Context, req *requests.FindAllReview) ([]*ReviewResult, error) {
	offset := (req.Page - 1) * req.PageSize
	filter := bson.M{"deleted_at": bson.M{"$eq": nil}}
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": req.Search, "$options": "i"}},
			{"comment": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}
	return r.runQuery(ctx, filter, offset, req.PageSize, review_errors.ErrFindAllReviews)
}

func (r *reviewQueryRepository) FindByProduct(ctx context.Context, req *requests.FindAllReviewByProduct) ([]*ReviewResult, error) {
	offset := (req.Page - 1) * req.PageSize
	filter := bson.M{"deleted_at": bson.M{"$eq": nil}, "product_id": req.ProductID}
	if req.Rating != 0 {
		filter["rating"] = req.Rating
	}
	return r.runQuery(ctx, filter, offset, req.PageSize, review_errors.ErrFindReviewsByProduct)
}

func (r *reviewQueryRepository) FindByMerchant(ctx context.Context, req *requests.FindAllReviewByMerchant) ([]*ReviewResult, error) {
	offset := (req.Page - 1) * req.PageSize
	cursor, err := r.catalogDB.Collection("products").Find(ctx, bson.M{"merchant_id": req.MerchantID})
	if err != nil {
		return nil, review_errors.ErrFindReviewsByMerchant.WithInternal(err)
	}
	var products []bson.M
	if err := cursor.All(ctx, &products); err != nil {
		cursor.Close(ctx)
		return nil, review_errors.ErrFindReviewsByMerchant.WithInternal(err)
	}
	cursor.Close(ctx)
	productIDs := make([]int32, 0, len(products))
	for _, p := range products {
		productIDs = append(productIDs, bsonMToInt32Rev(p, "_id"))
	}
	filter := bson.M{"deleted_at": bson.M{"$eq": nil}, "product_id": bson.M{"$in": productIDs}}
	if req.Rating != 0 {
		filter["rating"] = req.Rating
	}
	return r.runQuery(ctx, filter, offset, req.PageSize, review_errors.ErrFindReviewsByMerchant)
}

func (r *reviewQueryRepository) FindActive(ctx context.Context, req *requests.FindAllReview) ([]*ReviewResult, error) {
	return r.FindAll(ctx, req)
}

func (r *reviewQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllReview) ([]*ReviewResult, error) {
	offset := (req.Page - 1) * req.PageSize
	filter := bson.M{"deleted_at": bson.M{"$ne": nil}}
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": req.Search, "$options": "i"}},
			{"comment": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}
	return r.runQuery(ctx, filter, offset, req.PageSize, review_errors.ErrFindTrashedReviews)
}

func (r *reviewQueryRepository) FindByID(ctx context.Context, id int) (*ReviewResult, error) {
	var d bson.M
	err := r.db.Collection("reviews").FindOne(ctx, bson.M{"_id": id, "deleted_at": bson.M{"$eq": nil}}).Decode(&d)
	if err != nil {
		return nil, review_errors.ErrReviewNotFound
	}
	return bsonMToReviewResult(d), nil
}

func (r *reviewQueryRepository) runQuery(ctx context.Context, filter bson.M, offset, pageSize int, errFactory *errors.AppError) ([]*ReviewResult, error) {
	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": bson.M{"review_id": -1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data":     []bson.M{{"$skip": offset}, {"$limit": pageSize}},
		}},
	}
	cursor, err := r.db.Collection("reviews").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errFactory.WithInternal(err)
	}
	defer cursor.Close(ctx)
	var facetResults []struct {
		Metadata []struct{ TotalCount int64 `bson:"total_count"` } `bson:"metadata"`
		Data     []bson.M `bson:"data"`
	}
	if err := cursor.All(ctx, &facetResults); err != nil {
		return nil, errFactory.WithInternal(err)
	}
	if len(facetResults) == 0 {
		return []*ReviewResult{}, nil
	}
	var results []*ReviewResult
	for _, d := range facetResults[0].Data {
		item := bsonMToReviewResult(d)
		if len(facetResults[0].Metadata) > 0 {
			item.TotalCount = facetResults[0].Metadata[0].TotalCount
		}
		results = append(results, item)
	}
	return results, nil
}

func bsonMToReviewResult(m bson.M) *ReviewResult {
	return &ReviewResult{
		ReviewID:  bsonMToInt32Rev(m, "_id"),
		UserID:    bsonMToInt32Rev(m, "user_id"),
		ProductID: bsonMToInt32Rev(m, "product_id"),
		Name:      bsonMToStrRev(m, "name"),
		Comment:   bsonMToStrRev(m, "comment"),
		Rating:    bsonMToInt32Rev(m, "rating"),
		CreatedAt: bsonMToTimeStrPtrRev(m, "created_at"),
		UpdatedAt: bsonMToTimeStrPtrRev(m, "updated_at"),
		DeletedAt: bsonMToTimeStrPtrRev(m, "deleted_at"),
	}
}

func bsonMToInt32Rev(m bson.M, key string) int32 {
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

func bsonMToStrRev(m bson.M, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func bsonMToTimeStrPtrRev(m bson.M, key string) *string {
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
