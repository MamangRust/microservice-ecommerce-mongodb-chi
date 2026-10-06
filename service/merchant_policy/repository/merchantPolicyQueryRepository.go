package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	merchantpolicy_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant_policy_errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type merchantPolicyQueryRepository struct {
	db *mongo.Database
}

func NewMerchantPolicyQueryRepository(db *mongo.Database) MerchantPoliciesQueryRepository {
	return &merchantPolicyQueryRepository{db: db}
}

func (r *merchantPolicyQueryRepository) FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantPolicyResult, error) {
	return r.findPaginated(ctx, bson.M{"deleted_at": bson.M{"$eq": nil}}, req)
}

func (r *merchantPolicyQueryRepository) FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantPolicyResult, error) {
	return r.findPaginated(ctx, bson.M{"deleted_at": bson.M{"$eq": nil}}, req)
}

func (r *merchantPolicyQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantPolicyResult, error) {
	return r.findPaginated(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, req)
}

func (r *merchantPolicyQueryRepository) findPaginated(ctx context.Context, baseFilter bson.M, req *requests.FindAllMerchant) ([]*MerchantPolicyResult, error) {
	offset := (req.Page - 1) * req.PageSize
	filter := bson.M{}
	for k, v := range baseFilter {
		filter[k] = v
	}
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"policy_type": bson.M{"$regex": req.Search, "$options": "i"}},
			{"title": bson.M{"$regex": req.Search, "$options": "i"}},
			{"description": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}
	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": bson.M{"created_at": -1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data":     []bson.M{{"$skip": offset}, {"$limit": req.PageSize}},
		}},
	}

	cursor, err := r.db.Collection("merchant_policies").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, merchantpolicy_errors.ErrFindAllMerchantPolicies.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var facetResults []struct {
		Metadata []struct {
			TotalCount int64 `bson:"total_count"`
		} `bson:"metadata"`
		Data []bson.M `bson:"data"`
	}
	if err := cursor.All(ctx, &facetResults); err != nil {
		return nil, merchantpolicy_errors.ErrFindAllMerchantPolicies.WithInternal(err)
	}
	if len(facetResults) == 0 {
		return []*MerchantPolicyResult{}, nil
	}

	var results []*MerchantPolicyResult
	for _, d := range facetResults[0].Data {
		item := bsonMToMerchantPolicyResult(d)
		if len(facetResults[0].Metadata) > 0 {
			item.TotalCount = facetResults[0].Metadata[0].TotalCount
		}
		results = append(results, item)
	}
	return results, nil
}

func (r *merchantPolicyQueryRepository) FindByID(ctx context.Context, id int) (*MerchantPolicyResult, error) {
	var d bson.M
	err := r.db.Collection("merchant_policies").FindOne(ctx, bson.M{"_id": id, "deleted_at": bson.M{"$eq": nil}}).Decode(&d)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, merchantpolicy_errors.ErrMerchantPolicyNotFound
		}
		return nil, merchantpolicy_errors.ErrMerchantPolicyNotFound.WithInternal(err)
	}
	return bsonMToMerchantPolicyResult(d), nil
}

func bsonMToMerchantPolicyResult(m bson.M) *MerchantPolicyResult {
	return &MerchantPolicyResult{
		MerchantPolicyID: bsonMToInt32MP(m, "_id"),
		MerchantID:       bsonMToInt32MP(m, "merchant_id"),
		PolicyType:       bsonMToStrMP(m, "policy_type"),
		Title:            bsonMToStrMP(m, "title"),
		Description:      bsonMToStrMP(m, "description"),
		CreatedAt:        bsonMToTimePtrMP(m, "created_at"),
		UpdatedAt:        bsonMToTimePtrMP(m, "updated_at"),
		DeletedAt:        bsonMToTimePtrMP(m, "deleted_at"),
	}
}

func bsonMToInt32MP(m bson.M, key string) int32 {
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

func bsonMToStrMP(m bson.M, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func bsonMToTimePtrMP(m bson.M, key string) *time.Time {
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
