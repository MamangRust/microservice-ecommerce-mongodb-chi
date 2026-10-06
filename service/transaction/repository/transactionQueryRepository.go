package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/transaction_errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type transactionQueryRepository struct {
	db *mongo.Database
}

func NewTransactionQueryRepository(db *mongo.Database) *transactionQueryRepository {
	return &transactionQueryRepository{db: db}
}

func (r *transactionQueryRepository) FindAll(ctx context.Context, req *requests.FindAllTransaction) ([]*TransactionResult, error) {
	return r.runQuery(ctx, bson.M{"deleted_at": bson.M{"$eq": nil}}, req, transaction_errors.ErrFindAllTransactions)
}

func (r *transactionQueryRepository) FindActive(ctx context.Context, req *requests.FindAllTransaction) ([]*TransactionResult, error) {
	return r.runQuery(ctx, bson.M{"deleted_at": bson.M{"$eq": nil}}, req, transaction_errors.ErrFindByActive)
}

func (r *transactionQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllTransaction) ([]*TransactionResult, error) {
	return r.runQuery(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, req, transaction_errors.ErrFindByTrashed)
}

func (r *transactionQueryRepository) FindByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant) ([]*TransactionResult, error) {
	filter := bson.M{"deleted_at": bson.M{"$eq": nil}, "merchant_id": req.MerchantID}
	offset := (req.Page - 1) * req.PageSize
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"payment_method": bson.M{"$regex": req.Search, "$options": "i"}},
			{"payment_status": bson.M{"$regex": req.Search, "$options": "i"}},
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
	cursor, err := r.db.Collection("transactions").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, transaction_errors.ErrFindByMerchant.WithInternal(err)
	}
	defer cursor.Close(ctx)
	var facetResults []struct {
		Metadata []struct {
			TotalCount int64 `bson:"total_count"`
		} `bson:"metadata"`
		Data []bson.M `bson:"data"`
	}
	if err := cursor.All(ctx, &facetResults); err != nil {
		return nil, transaction_errors.ErrFindByMerchant.WithInternal(err)
	}
	if len(facetResults) == 0 {
		return []*TransactionResult{}, nil
	}
	var results []*TransactionResult
	for _, d := range facetResults[0].Data {
		item := bsonMToTransactionResult(d)
		if len(facetResults[0].Metadata) > 0 {
			item.TotalCount = facetResults[0].Metadata[0].TotalCount
		}
		results = append(results, item)
	}
	return results, nil
}

func (r *transactionQueryRepository) FindByID(ctx context.Context, transactionID int) (*TransactionResult, error) {
	var d bson.M
	err := r.db.Collection("transactions").FindOne(ctx, bson.M{"_id": transactionID, "deleted_at": bson.M{"$eq": nil}}).Decode(&d)
	if err != nil {
		return nil, transaction_errors.ErrTransactionNotFound
	}
	return bsonMToTransactionResult(d), nil
}

func (r *transactionQueryRepository) FindByOrderID(ctx context.Context, orderID int) (*TransactionResult, error) {
	var d bson.M
	err := r.db.Collection("transactions").FindOne(ctx, bson.M{"order_id": orderID, "deleted_at": bson.M{"$eq": nil}}).Decode(&d)
	if err != nil {
		return nil, transaction_errors.ErrTransactionNotFound
	}
	return bsonMToTransactionResult(d), nil
}

func (r *transactionQueryRepository) runQuery(ctx context.Context, filter bson.M, req *requests.FindAllTransaction, errFactory *errors.AppError) ([]*TransactionResult, error) {
	offset := (req.Page - 1) * req.PageSize
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"payment_method": bson.M{"$regex": req.Search, "$options": "i"}},
			{"payment_status": bson.M{"$regex": req.Search, "$options": "i"}},
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
	cursor, err := r.db.Collection("transactions").Aggregate(ctx, pipeline)
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
		return []*TransactionResult{}, nil
	}
	var results []*TransactionResult
	for _, d := range facetResults[0].Data {
		item := bsonMToTransactionResult(d)
		if len(facetResults[0].Metadata) > 0 {
			item.TotalCount = facetResults[0].Metadata[0].TotalCount
		}
		results = append(results, item)
	}
	return results, nil
}

func bsonMToTransactionResult(m bson.M) *TransactionResult {
	return &TransactionResult{
		TransactionID: bsonMToInt32Tr(m, "_id"),
		OrderID:       bsonMToInt32Tr(m, "order_id"),
		MerchantID:    bsonMToInt32Tr(m, "merchant_id"),
		PaymentMethod: bsonMToStrTr(m, "payment_method"),
		Amount:        bsonMToInt32Tr(m, "amount"),
		PaymentStatus: bsonMToStrTr(m, "payment_status"),
		CreatedAt:     bsonMToTimePtrTr(m, "created_at"),
		UpdatedAt:     bsonMToTimePtrTr(m, "updated_at"),
		DeletedAt:     bsonMToTimePtrTr(m, "deleted_at"),
	}
}

func bsonMToInt32Tr(m bson.M, key string) int32 {
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

func bsonMToStrTr(m bson.M, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func bsonMToTimePtrTr(m bson.M, key string) *time.Time {
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
