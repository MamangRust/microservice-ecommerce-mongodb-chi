package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/cart_errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type cartQueryRepository struct {
	db *mongo.Database
}

func NewCartQueryRepository(db *mongo.Database) CartQueryRepository {
	return &cartQueryRepository{db: db}
}

func (r *cartQueryRepository) FindCarts(ctx context.Context, req *requests.FindAllCarts) ([]*CartResult, error) {
	filter := bson.M{"user_id": req.UserID, "deleted_at": bson.M{"$eq": nil}}
	if req.Search != "" {
		filter["name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}
	cursor, err := r.db.Collection("carts").Find(ctx, filter)
	if err != nil {
		return nil, cart_errors.ErrFindAllCarts.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var docs []bson.M
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, cart_errors.ErrFindAllCarts.WithInternal(err)
	}

	total := int64(len(docs))
	var out []*CartResult
	for _, d := range docs {
		out = append(out, &CartResult{
			CartID:     bsonMToInt32Cart(d, "_id"),
			UserID:     bsonMToInt32Cart(d, "user_id"),
			ProductID:  bsonMToInt32Cart(d, "product_id"),
			Name:       bsonMToStrCart(d, "name"),
			Price:      bsonMToInt32Cart(d, "price"),
			Image:      bsonMToStrCart(d, "image"),
			Quantity:   bsonMToInt32Cart(d, "quantity"),
			Weight:     bsonMToInt32Cart(d, "weight"),
			CreatedAt:  bsonMToTimeStrCart(d, "created_at"),
			UpdatedAt:  bsonMToTimeStrCart(d, "updated_at"),
			TotalCount: total,
		})
	}
	return out, nil
}

func bsonMToInt32Cart(m bson.M, key string) int32 {
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

func bsonMToStrCart(m bson.M, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func bsonMToTimeStrCart(m bson.M, key string) string {
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	switch t := v.(type) {
	case time.Time:
		return t.Format("2006-01-02 15:04:05")
	case bson.DateTime:
		return t.Time().Format("2006-01-02 15:04:05")
	case string:
		return t
	default:
		return ""
	}
}
