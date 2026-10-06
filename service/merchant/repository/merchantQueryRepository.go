package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	merchant_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type merchantQueryRepository struct {
	db *mongo.Database
}

func NewMerchantQueryRepository(db *mongo.Database) *merchantQueryRepository {
	return &merchantQueryRepository{db: db}
}

func (r *merchantQueryRepository) FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantResult, error) {
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
		{"$sort": bson.M{"merchant_id": -1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data":     []bson.M{{"$skip": offset}, {"$limit": req.PageSize}},
		}},
	}
	results, totalCount, err := runPaginatedMerchantQuery(ctx, r.db.Collection("merchants"), pipeline)
	if err != nil {
		return nil, merchant_errors.ErrFindAllMerchants.WithInternal(err)
	}
	for _, item := range results {
		item.TotalCount = totalCount
	}
	return results, nil
}

func (r *merchantQueryRepository) FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantResult, error) {
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
		{"$sort": bson.M{"merchant_id": -1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data":     []bson.M{{"$skip": offset}, {"$limit": req.PageSize}},
		}},
	}
	results, totalCount, err := runPaginatedMerchantQuery(ctx, r.db.Collection("merchants"), pipeline)
	if err != nil {
		return nil, merchant_errors.ErrFindActiveMerchants.WithInternal(err)
	}
	for _, item := range results {
		item.TotalCount = totalCount
	}
	return results, nil
}

func (r *merchantQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantResult, error) {
	offset := (req.Page - 1) * req.PageSize
	filter := bson.M{"deleted_at": bson.M{"$ne": nil}}
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": req.Search, "$options": "i"}},
			{"description": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}
	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": bson.M{"merchant_id": -1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data":     []bson.M{{"$skip": offset}, {"$limit": req.PageSize}},
		}},
	}
	results, totalCount, err := runPaginatedMerchantQuery(ctx, r.db.Collection("merchants"), pipeline)
	if err != nil {
		return nil, merchant_errors.ErrFindTrashedMerchants.WithInternal(err)
	}
	for _, item := range results {
		item.TotalCount = totalCount
	}
	return results, nil
}

func (r *merchantQueryRepository) FindByID(ctx context.Context, merchantID int) (*MerchantResult, error) {
	var d bson.M
	err := r.db.Collection("merchants").FindOne(ctx, bson.M{"_id": merchantID, "deleted_at": bson.M{"$eq": nil}}).Decode(&d)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, merchant_errors.ErrMerchantNotFound
		}
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	return bsonMToMerchantResult(d), nil
}

func runPaginatedMerchantQuery(ctx context.Context, coll *mongo.Collection, pipeline []bson.M) ([]*MerchantResult, int64, error) {
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
		return []*MerchantResult{}, 0, nil
	}

	var totalCount int64
	if len(facetResults[0].Metadata) > 0 {
		totalCount = facetResults[0].Metadata[0].TotalCount
	}

	var results []*MerchantResult
	for _, d := range facetResults[0].Data {
		results = append(results, bsonMToMerchantResult(d))
	}
	return results, totalCount, nil
}

func bsonMToMerchantResult(m bson.M) *MerchantResult {
	return &MerchantResult{
		MerchantID:   bsonMToInt32Merchant(m, "_id"),
		UserID:       bsonMToInt32Merchant(m, "user_id"),
		Name:         bsonMToStrMerchant(m, "name"),
		Description:  bsonMToStrPtrMerchant(m, "description"),
		Address:      bsonMToStrPtrMerchant(m, "address"),
		ContactEmail: bsonMToStrPtrMerchant(m, "contact_email"),
		ContactPhone: bsonMToStrPtrMerchant(m, "contact_phone"),
		Status:       bsonMToStrMerchant(m, "status"),
		CreatedAt:    bsonMToStrPtrMerchant(m, "created_at"),
		UpdatedAt:    bsonMToStrPtrMerchant(m, "updated_at"),
		DeletedAt:    bsonMToStrPtrMerchant(m, "deleted_at"),
	}
}

func bsonMToInt32Merchant(m bson.M, key string) int32 {
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

func bsonMToStrMerchant(m bson.M, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func bsonMToStrPtrMerchant(m bson.M, key string) *string {
	v, ok := m[key]
	if !ok || v == nil {
		return nil
	}
	switch t := v.(type) {
	case string:
		return &t
	case time.Time:
		s := t.Format("2006-01-02 15:04:05")
		return &s
	case bson.DateTime:
		s := t.Time().Format("2006-01-02 15:04:05")
		return &s
	default:
		return nil
	}
}
