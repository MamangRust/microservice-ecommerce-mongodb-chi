package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	merchantaward_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant_award"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type merchantAwardQueryRepository struct {
	db *mongo.Database
}

func NewMerchantAwardQueryRepository(db *mongo.Database) MerchantAwardQueryRepository {
	return &merchantAwardQueryRepository{db: db}
}

func (r *merchantAwardQueryRepository) FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantCertResult, error) {
	return r.findPaginated(ctx, bson.M{"deleted_at": bson.M{"$eq": nil}}, req)
}

func (r *merchantAwardQueryRepository) FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantCertResult, error) {
	return r.findPaginated(ctx, bson.M{"deleted_at": bson.M{"$eq": nil}}, req)
}

func (r *merchantAwardQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantCertResult, error) {
	return r.findPaginated(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, req)
}

func (r *merchantAwardQueryRepository) findPaginated(ctx context.Context, baseFilter bson.M, req *requests.FindAllMerchant) ([]*MerchantCertResult, error) {
	offset := (req.Page - 1) * req.PageSize
	filter := bson.M{}
	for k, v := range baseFilter {
		filter[k] = v
	}
	if req.Search != "" {
		filter["title"] = bson.M{"$regex": req.Search, "$options": "i"}
	}
	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": bson.M{"created_at": -1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data":     []bson.M{{"$skip": offset}, {"$limit": req.PageSize}},
		}},
	}

	cursor, err := r.db.Collection("merchant_certifications_and_awards").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, merchantaward_errors.ErrFindAllMerchantAwards.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var facetResults []struct {
		Metadata []struct {
			TotalCount int64 `bson:"total_count"`
		} `bson:"metadata"`
		Data []bson.M `bson:"data"`
	}
	if err := cursor.All(ctx, &facetResults); err != nil {
		return nil, merchantaward_errors.ErrFindAllMerchantAwards.WithInternal(err)
	}
	if len(facetResults) == 0 {
		return []*MerchantCertResult{}, nil
	}

	var results []*MerchantCertResult
	for _, d := range facetResults[0].Data {
		item := bsonMToMerchantCertResult(d)
		if len(facetResults[0].Metadata) > 0 {
			item.TotalCount = facetResults[0].Metadata[0].TotalCount
		}
		results = append(results, item)
	}
	return results, nil
}

func (r *merchantAwardQueryRepository) FindByID(ctx context.Context, id int) (*MerchantCertResult, error) {
	var d bson.M
	err := r.db.Collection("merchant_certifications_and_awards").FindOne(ctx, bson.M{"_id": id, "deleted_at": bson.M{"$eq": nil}}).Decode(&d)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, merchantaward_errors.ErrMerchantAwardNotFound
		}
		return nil, merchantaward_errors.ErrMerchantAwardNotFound.WithInternal(err)
	}
	return bsonMToMerchantCertResult(d), nil
}

func bsonMToMerchantCertResult(m bson.M) *MerchantCertResult {
	return &MerchantCertResult{
		MerchantCertificationID: bsonMToInt32MA(m, "_id"),
		MerchantID:              bsonMToInt32MA(m, "merchant_id"),
		Title:                   bsonMToStrMA(m, "title"),
		Description:             bsonMToStrPtrMA(m, "description"),
		IssuedBy:                bsonMToStrPtrMA(m, "issued_by"),
		IssueDate:               bsonMToTimePtrMA(m, "issue_date"),
		ExpiryDate:              bsonMToTimePtrMA(m, "expiry_date"),
		CertificateUrl:          bsonMToStrPtrMA(m, "certificate_url"),
		CreatedAt:               bsonMToTimePtrMA(m, "created_at"),
		UpdatedAt:               bsonMToTimePtrMA(m, "updated_at"),
		DeletedAt:               bsonMToTimePtrMA(m, "deleted_at"),
	}
}

func bsonMToTimePtrMA(m bson.M, key string) *time.Time {
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

func bsonMToInt32MA(m bson.M, key string) int32 {
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

func bsonMToStrMA(m bson.M, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func bsonMToStrPtrMA(m bson.M, key string) *string {
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
