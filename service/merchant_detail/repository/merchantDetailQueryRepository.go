package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	merchantdetail_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant_detail"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type merchantDetailQueryRepository struct {
	db *mongo.Database
}

func NewMerchantDetailQueryRepository(db *mongo.Database) *merchantDetailQueryRepository {
	return &merchantDetailQueryRepository{db: db}
}

func (r *merchantDetailQueryRepository) FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantDetailResult, error) {
	offset := (req.Page - 1) * req.PageSize
	filter := bson.M{"deleted_at": bson.M{"$eq": nil}}
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"display_name": bson.M{"$regex": req.Search, "$options": "i"}},
			{"short_description": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}
	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": bson.M{"merchant_detail_id": -1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data":     []bson.M{{"$skip": offset}, {"$limit": req.PageSize}},
		}},
	}
	results, totalCount, err := runPaginatedMDQuery(ctx, r.db.Collection("merchant_details"), pipeline)
	if err != nil {
		return nil, merchantdetail_errors.ErrFindAllMerchantDetails.WithInternal(err)
	}
	for _, item := range results {
		item.TotalCount = totalCount
	}
	return results, nil
}

func (r *merchantDetailQueryRepository) FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantDetailResult, error) {
	return r.FindAll(ctx, req)
}

func (r *merchantDetailQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantDetailResult, error) {
	offset := (req.Page - 1) * req.PageSize
	filter := bson.M{"deleted_at": bson.M{"$ne": nil}}
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"display_name": bson.M{"$regex": req.Search, "$options": "i"}},
			{"short_description": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}
	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": bson.M{"merchant_detail_id": -1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data":     []bson.M{{"$skip": offset}, {"$limit": req.PageSize}},
		}},
	}
	results, totalCount, err := runPaginatedMDQuery(ctx, r.db.Collection("merchant_details"), pipeline)
	if err != nil {
		return nil, merchantdetail_errors.ErrFindTrashedMerchantDetails.WithInternal(err)
	}
	for _, item := range results {
		item.TotalCount = totalCount
	}
	return results, nil
}

func (r *merchantDetailQueryRepository) FindByID(ctx context.Context, merchantDetailID int) (*MerchantDetailResult, error) {
	var d bson.M
	err := r.db.Collection("merchant_details").FindOne(ctx, bson.M{"_id": merchantDetailID, "deleted_at": bson.M{"$eq": nil}}).Decode(&d)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, merchantdetail_errors.ErrMerchantDetailNotFound
		}
		return nil, merchantdetail_errors.ErrMerchantDetailNotFound.WithInternal(err)
	}
	return bsonMToMerchantDetailResult(d), nil
}

func runPaginatedMDQuery(ctx context.Context, coll *mongo.Collection, pipeline []bson.M) ([]*MerchantDetailResult, int64, error) {
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
		return []*MerchantDetailResult{}, 0, nil
	}

	var totalCount int64
	if len(facetResults[0].Metadata) > 0 {
		totalCount = facetResults[0].Metadata[0].TotalCount
	}

	var results []*MerchantDetailResult
	for _, d := range facetResults[0].Data {
		results = append(results, bsonMToMerchantDetailResult(d))
	}
	return results, totalCount, nil
}

func bsonMToMerchantDetailResult(m bson.M) *MerchantDetailResult {
	return &MerchantDetailResult{
		MerchantDetailID: bsonMToInt32MD(m, "_id"),
		MerchantID:       bsonMToInt32MD(m, "merchant_id"),
		DisplayName:      bsonMToStrPtrMD(m, "display_name"),
		CoverImageUrl:    bsonMToStrPtrMD(m, "cover_image_url"),
		LogoUrl:          bsonMToStrPtrMD(m, "logo_url"),
		ShortDescription: bsonMToStrPtrMD(m, "short_description"),
		WebsiteUrl:       bsonMToStrPtrMD(m, "website_url"),
		CreatedAt:        bsonMToStrPtrMD(m, "created_at"),
		UpdatedAt:        bsonMToStrPtrMD(m, "updated_at"),
		DeletedAt:        bsonMToStrPtrMD(m, "deleted_at"),
	}
}

func bsonMToInt32MD(m bson.M, key string) int32 {
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

func bsonMToStrPtrMD(m bson.M, key string) *string {
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
