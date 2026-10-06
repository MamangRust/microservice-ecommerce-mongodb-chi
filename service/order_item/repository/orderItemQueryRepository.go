package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	orderitem_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/order_item_errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type orderItemQueryRepository struct {
	db *mongo.Database
}

func NewOrderItemQueryRepository(db *mongo.Database) *orderItemQueryRepository {
	return &orderItemQueryRepository{db: db}
}

func (r *orderItemQueryRepository) FindAll(ctx context.Context, req *requests.FindAllOrderItems) ([]*OrderItemResult, error) {
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
	results, totalCount, err := runPaginatedOIQuery(ctx, r.db.Collection("order_items"), pipeline)
	if err != nil {
		return nil, orderitem_errors.ErrFindAllOrderItems.WithInternal(err)
	}
	for _, item := range results {
		item.TotalCount = totalCount
	}
	return results, nil
}

func (r *orderItemQueryRepository) FindActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*OrderItemResult, error) {
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
	results, totalCount, err := runPaginatedOIQuery(ctx, r.db.Collection("order_items"), pipeline)
	if err != nil {
		return nil, orderitem_errors.ErrFindByActive.WithInternal(err)
	}
	for _, item := range results {
		item.TotalCount = totalCount
	}
	return results, nil
}

func (r *orderItemQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*OrderItemResult, error) {
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
	results, totalCount, err := runPaginatedOIQuery(ctx, r.db.Collection("order_items"), pipeline)
	if err != nil {
		return nil, orderitem_errors.ErrFindByTrashed.WithInternal(err)
	}
	for _, item := range results {
		item.TotalCount = totalCount
	}
	return results, nil
}

func (r *orderItemQueryRepository) FindOrderItemByOrder(ctx context.Context, orderID int) ([]*OrderItemResult, error) {
	filter := bson.M{"order_id": orderID, "deleted_at": bson.M{"$eq": nil}}
	cursor, err := r.db.Collection("order_items").Find(ctx, filter)
	if err != nil {
		return nil, orderitem_errors.ErrFindOrderItemByOrder.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var items []bson.M
	if err := cursor.All(ctx, &items); err != nil {
		return nil, orderitem_errors.ErrFindOrderItemByOrder.WithInternal(err)
	}

	var results []*OrderItemResult
	for _, d := range items {
		results = append(results, bsonMToOrderItemResult(d))
	}
	return results, nil
}

func runPaginatedOIQuery(ctx context.Context, coll *mongo.Collection, pipeline []bson.M) ([]*OrderItemResult, int64, error) {
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
		return []*OrderItemResult{}, 0, nil
	}

	var totalCount int64
	if len(facetResults[0].Metadata) > 0 {
		totalCount = facetResults[0].Metadata[0].TotalCount
	}

	var results []*OrderItemResult
	for _, d := range facetResults[0].Data {
		results = append(results, bsonMToOrderItemResult(d))
	}
	return results, totalCount, nil
}

func bsonMToOrderItemResult(m bson.M) *OrderItemResult {
	return &OrderItemResult{
		OrderItemID: bsonMToInt32OI(m, "_id"),
		OrderID:     bsonMToInt32OI(m, "order_id"),
		ProductID:   bsonMToInt32OI(m, "product_id"),
		Quantity:    bsonMToInt32OI(m, "quantity"),
		Price:       bsonMToInt32OI(m, "price"),
		CreatedAt:   bsonMToTimePtrOI(m, "created_at"),
		UpdatedAt:   bsonMToTimePtrOI(m, "updated_at"),
		DeletedAt:   bsonMToTimePtrOI(m, "deleted_at"),
	}
}

func bsonMToInt32OI(m bson.M, key string) int32 {
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

func bsonMToTimePtrOI(m bson.M, key string) *time.Time {
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
