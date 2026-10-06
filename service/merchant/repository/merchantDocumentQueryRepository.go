package repository

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	merchant_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type merchantDocumentQueryRepository struct {
	db *mongo.Database
}

func NewMerchantDocumentQueryRepository(db *mongo.Database) *merchantDocumentQueryRepository {
	return &merchantDocumentQueryRepository{db: db}
}

func (r *merchantDocumentQueryRepository) FindAll(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*MerchantDocumentResult, *int, error) {
	return r.findPaginated(ctx, bson.M{"deleted_at": bson.M{"$eq": nil}}, req)
}

func (r *merchantDocumentQueryRepository) FindActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*MerchantDocumentResult, *int, error) {
	return r.findPaginated(ctx, bson.M{"deleted_at": bson.M{"$eq": nil}}, req)
}

func (r *merchantDocumentQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*MerchantDocumentResult, *int, error) {
	return r.findPaginated(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, req)
}

func (r *merchantDocumentQueryRepository) findPaginated(ctx context.Context, baseFilter bson.M, req *requests.FindAllMerchantDocuments) ([]*MerchantDocumentResult, *int, error) {
	offset := (req.Page - 1) * req.PageSize
	filter := bson.M{}
	for k, v := range baseFilter {
		filter[k] = v
	}
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"document_type": bson.M{"$regex": req.Search, "$options": "i"}},
			{"status": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}

	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": bson.M{"document_id": -1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data":     []bson.M{{"$skip": offset}, {"$limit": req.PageSize}},
		}},
	}

	cursor, err := r.db.Collection("merchant_documents").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var facetResults []struct {
		Metadata []struct{ TotalCount int64 `bson:"total_count"` } `bson:"metadata"`
		Data     []bson.M `bson:"data"`
	}
	if err := cursor.All(ctx, &facetResults); err != nil {
		return nil, nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}

	var totalCount int
	if len(facetResults) > 0 && len(facetResults[0].Metadata) > 0 {
		totalCount = int(facetResults[0].Metadata[0].TotalCount)
	}

	var results []*MerchantDocumentResult
	if len(facetResults) > 0 {
		for _, d := range facetResults[0].Data {
			results = append(results, bsonMToMerchantDocumentResult(d))
		}
	}
	return results, &totalCount, nil
}

func (r *merchantDocumentQueryRepository) FindByID(ctx context.Context, id int) (*MerchantDocumentResult, error) {
	var d bson.M
	err := r.db.Collection("merchant_documents").FindOne(ctx, bson.M{"_id": id}).Decode(&d)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, merchant_errors.ErrMerchantNotFound
		}
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	return bsonMToMerchantDocumentResult(d), nil
}

func bsonMToMerchantDocumentResult(m bson.M) *MerchantDocumentResult {
	return &MerchantDocumentResult{
		DocumentID:   bsonMToInt32Merchant(m, "_id"),
		MerchantID:   bsonMToInt32Merchant(m, "merchant_id"),
		DocumentType: bsonMToStrMerchant(m, "document_type"),
		DocumentUrl:  bsonMToStrMerchant(m, "document_url"),
		Status:       bsonMToStrMerchant(m, "status"),
		Note:         bsonMToStrPtrMerchant(m, "note"),
		UploadedAt:   bsonMToStrPtrMerchant(m, "uploaded_at"),
		CreatedAt:    bsonMToStrPtrMerchant(m, "created_at"),
		UpdatedAt:    bsonMToStrPtrMerchant(m, "updated_at"),
		DeletedAt:    bsonMToStrPtrMerchant(m, "deleted_at"),
	}
}
