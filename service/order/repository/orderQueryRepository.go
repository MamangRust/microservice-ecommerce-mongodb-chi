package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/order_errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type orderQueryRepository struct {
	db *mongo.Database
}

func NewOrderQueryRepository(db *mongo.Database) OrderQueryRepository {
	return &orderQueryRepository{db: db}
}

func (r *orderQueryRepository) FindAll(ctx context.Context, req *requests.FindAllOrder) ([]*OrderResult, error) {
	offset := (req.Page - 1) * req.PageSize
	filter := bson.M{"deleted_at": bson.M{"$eq": nil}}
	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": bson.M{"created_at": -1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data":     []bson.M{{"$skip": offset}, {"$limit": req.PageSize}},
		}},
	}
	results, totalCount, err := runPaginatedOrderQuery(ctx, r.db.Collection("orders"), pipeline)
	if err != nil {
		return nil, order_errors.ErrFindAllOrders.WithInternal(err)
	}
	for _, item := range results {
		item.TotalCount = totalCount
	}
	return results, nil
}

func (r *orderQueryRepository) FindActive(ctx context.Context, req *requests.FindAllOrder) ([]*OrderResult, error) {
	offset := (req.Page - 1) * req.PageSize
	filter := bson.M{"deleted_at": bson.M{"$eq": nil}}
	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": bson.M{"created_at": -1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data":     []bson.M{{"$skip": offset}, {"$limit": req.PageSize}},
		}},
	}
	results, totalCount, err := runPaginatedOrderQuery(ctx, r.db.Collection("orders"), pipeline)
	if err != nil {
		return nil, order_errors.ErrFindByActive.WithInternal(err)
	}
	for _, item := range results {
		item.TotalCount = totalCount
	}
	return results, nil
}

func (r *orderQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllOrder) ([]*OrderResult, error) {
	offset := (req.Page - 1) * req.PageSize
	filter := bson.M{"deleted_at": bson.M{"$ne": nil}}
	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": bson.M{"created_at": -1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data":     []bson.M{{"$skip": offset}, {"$limit": req.PageSize}},
		}},
	}
	results, totalCount, err := runPaginatedOrderQuery(ctx, r.db.Collection("orders"), pipeline)
	if err != nil {
		return nil, order_errors.ErrFindByTrashed.WithInternal(err)
	}
	for _, item := range results {
		item.TotalCount = totalCount
	}
	return results, nil
}

func (r *orderQueryRepository) FindByMerchant(ctx context.Context, req *requests.FindAllOrderByMerchant) ([]*OrderResult, error) {
	offset := (req.Page - 1) * req.PageSize
	filter := bson.M{"deleted_at": bson.M{"$eq": nil}, "merchant_id": req.MerchantID}
	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": bson.M{"created_at": -1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data":     []bson.M{{"$skip": offset}, {"$limit": req.PageSize}},
		}},
	}
	results, totalCount, err := runPaginatedOrderQuery(ctx, r.db.Collection("orders"), pipeline)
	if err != nil {
		return nil, order_errors.ErrFindByMerchant.WithInternal(err)
	}
	for _, item := range results {
		item.TotalCount = totalCount
	}
	return results, nil
}

func (r *orderQueryRepository) FindByID(ctx context.Context, id int) (*models.Order, error) {
	var order models.Order
	err := r.db.Collection("orders").FindOne(ctx, bson.M{"_id": id, "deleted_at": bson.M{"$eq": nil}}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, order_errors.ErrOrderNotFound.WithInternal(err)
		}
		return nil, order_errors.ErrFindById.WithInternal(err)
	}
	return &order, nil
}

func runPaginatedOrderQuery(ctx context.Context, coll *mongo.Collection, pipeline []bson.M) ([]*OrderResult, int64, error) {
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
		return []*OrderResult{}, 0, nil
	}

	var totalCount int64
	if len(facetResults[0].Metadata) > 0 {
		totalCount = facetResults[0].Metadata[0].TotalCount
	}

	var results []*OrderResult
	for _, d := range facetResults[0].Data {
		results = append(results, &OrderResult{
			OrderID:    bsonMToInt32Order(d, "_id"),
			UserID:     bsonMToInt32Order(d, "user_id"),
			MerchantID: bsonMToInt32Order(d, "merchant_id"),
			TotalPrice: bsonMToInt32Order(d, "total_price"),
			CreatedAt:  bsonMToTimePtrOrder(d, "created_at"),
			UpdatedAt:  bsonMToTimePtrOrder(d, "updated_at"),
			DeletedAt:  bsonMToTimePtrOrder(d, "deleted_at"),
		})
	}
	return results, totalCount, nil
}

func bsonMToInt32Order(m bson.M, key string) int32 {
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

func bsonMToTimePtrOrder(m bson.M, key string) *time.Time {
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
