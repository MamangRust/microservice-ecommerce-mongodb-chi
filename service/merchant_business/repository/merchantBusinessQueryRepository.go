package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors"
	merchantbusiness_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant_business"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type merchantBusinessQueryRepository struct {
	db *mongo.Database
}

func NewMerchantBusinessQueryRepository(db *mongo.Database) MerchantBusinessQueryRepository {
	return &merchantBusinessQueryRepository{db: db}
}

func (r *merchantBusinessQueryRepository) FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantBusinessResult, error) {
	return r.findPaginated(ctx, bson.M{"deleted_at": bson.M{"$eq": nil}}, req, merchantbusiness_errors.ErrMerchantBusinessNotFound)
}

func (r *merchantBusinessQueryRepository) FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantBusinessResult, error) {
	return r.findPaginated(ctx, bson.M{"deleted_at": bson.M{"$eq": nil}}, req, merchantbusiness_errors.ErrFindActiveMerchantBusinesses)
}

func (r *merchantBusinessQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantBusinessResult, error) {
	return r.findPaginated(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, req, merchantbusiness_errors.ErrFindTrashedMerchantBusinesses)
}

func (r *merchantBusinessQueryRepository) findPaginated(ctx context.Context, baseFilter bson.M, req *requests.FindAllMerchant, errFactory *errors.AppError) ([]*MerchantBusinessResult, error) {
	offset := (req.Page - 1) * req.PageSize
	filter := bson.M{}
	for k, v := range baseFilter {
		filter[k] = v
	}
	if req.Search != "" {
		filter["business_type"] = bson.M{"$regex": req.Search, "$options": "i"}
	}
	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": bson.M{"created_at": -1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data":     []bson.M{{"$skip": offset}, {"$limit": req.PageSize}},
		}},
	}

	cursor, err := r.db.Collection("merchant_business_information").Aggregate(ctx, pipeline)
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
		return []*MerchantBusinessResult{}, nil
	}

	var results []*MerchantBusinessResult
	for _, d := range facetResults[0].Data {
		item := bsonMToMerchantBusinessResult(d)
		if len(facetResults[0].Metadata) > 0 {
			item.TotalCount = facetResults[0].Metadata[0].TotalCount
		}
		results = append(results, item)
	}
	return results, nil
}

func (r *merchantBusinessQueryRepository) FindByID(ctx context.Context, id int) (*MerchantBusinessResult, error) {
	var d bson.M
	err := r.db.Collection("merchant_business_information").FindOne(ctx, bson.M{"_id": id, "deleted_at": bson.M{"$eq": nil}}).Decode(&d)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, merchantbusiness_errors.ErrMerchantBusinessNotFound
		}
		return nil, merchantbusiness_errors.ErrMerchantBusinessNotFound.WithInternal(err)
	}
	return bsonMToMerchantBusinessResult(d), nil
}

func bsonMToMerchantBusinessResult(m bson.M) *MerchantBusinessResult {
	return &MerchantBusinessResult{
		MerchantBusinessInfoID: bsonMToInt32MB(m, "_id"),
		MerchantID:             bsonMToInt32MB(m, "merchant_id"),
		BusinessType:           bsonMToStrPtrMB(m, "business_type"),
		TaxID:                  bsonMToStrPtrMB(m, "tax_id"),
		EstablishedYear:        bsonMToInt32PtrMB(m, "established_year"),
		NumberOfEmployees:      bsonMToInt32PtrMB(m, "number_of_employees"),
		WebsiteUrl:             bsonMToStrPtrMB(m, "website_url"),
		CreatedAt:              bsonMToTimePtrMB(m, "created_at"),
		UpdatedAt:              bsonMToTimePtrMB(m, "updated_at"),
		DeletedAt:              bsonMToTimePtrMB(m, "deleted_at"),
	}
}

func bsonMToInt32MB(m bson.M, key string) int32 {
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

func bsonMToInt32PtrMB(m bson.M, key string) *int32 {
	v, ok := m[key]
	if !ok || v == nil {
		return nil
	}
	switch t := v.(type) {
	case int32:
		return &t
	case int64:
		i := int32(t)
		return &i
	case float64:
		i := int32(t)
		return &i
	default:
		return nil
	}
}

func bsonMToStrPtrMB(m bson.M, key string) *string {
	v, ok := m[key]
	if !ok || v == nil {
		return nil
	}
	if s, ok := v.(string); ok {
		return &s
	}
	return nil
}

func bsonMToTimePtrMB(m bson.M, key string) *time.Time {
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
