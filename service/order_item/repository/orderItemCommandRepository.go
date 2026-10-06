package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	orderitem_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/order_item_errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type orderItemCommandRepository struct {
	db *mongo.Database
}

func NewOrderItemCommandRepository(db *mongo.Database) *orderItemCommandRepository {
	return &orderItemCommandRepository{db: db}
}

func (r *orderItemCommandRepository) Create(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*models.OrderItem, error) {
	id, err := database.GetNextID(ctx, r.db, "order_items")
	if err != nil {
		return nil, orderitem_errors.ErrFailedCreateOrderItem.WithInternal(err)
	}
	now := time.Now()
	item := &models.OrderItem{
		OrderItemID: id,
		OrderID:     int32(req.OrderID),
		ProductID:   int32(req.ProductID),
		Quantity:    int32(req.Quantity),
		Price:       int32(req.Price),
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}
	if _, err := r.db.Collection("order_items").InsertOne(ctx, item); err != nil {
		return nil, orderitem_errors.ErrFailedCreateOrderItem.WithInternal(err)
	}
	return item, nil
}

func (r *orderItemCommandRepository) Update(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*models.OrderItem, error) {
	coll := r.db.Collection("order_items")
	var item models.OrderItem
	if err := coll.FindOne(ctx, bson.M{"_id": req.OrderItemID, "deleted_at": bson.M{"$eq": nil}}).Decode(&item); err != nil {
		return nil, orderitem_errors.ErrOrderItemNotFound
	}
	item.Quantity = int32(req.Quantity)
	item.Price = int32(req.Price)
	now := time.Now()
	item.UpdatedAt = &now
	if _, err := coll.ReplaceOne(ctx, bson.M{"_id": req.OrderItemID}, item); err != nil {
		return nil, orderitem_errors.ErrFailedUpdateOrderItem.WithInternal(err)
	}
	return &item, nil
}

func (r *orderItemCommandRepository) Trash(ctx context.Context, orderItemID int) (*models.OrderItem, error) {
	coll := r.db.Collection("order_items")
	var item models.OrderItem
	if err := coll.FindOne(ctx, bson.M{"_id": orderItemID, "deleted_at": bson.M{"$eq": nil}}).Decode(&item); err != nil {
		return nil, orderitem_errors.ErrOrderItemNotFound
	}
	now := time.Now()
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": orderItemID}, bson.M{"$set": bson.M{"deleted_at": now}}); err != nil {
		return nil, orderitem_errors.ErrTrashedOrderItem.WithInternal(err)
	}
	item.DeletedAt = &now
	return &item, nil
}

func (r *orderItemCommandRepository) Restore(ctx context.Context, orderItemID int) (*models.OrderItem, error) {
	coll := r.db.Collection("order_items")
	res, err := coll.UpdateOne(ctx, bson.M{"_id": orderItemID, "deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return nil, orderitem_errors.ErrRestoreOrderItem.WithInternal(err)
	}
	if res.MatchedCount == 0 {
		return nil, orderitem_errors.ErrOrderItemNotFound
	}
	var item models.OrderItem
	if err := coll.FindOne(ctx, bson.M{"_id": orderItemID}).Decode(&item); err != nil {
		return nil, orderitem_errors.ErrRestoreOrderItem.WithInternal(err)
	}
	return &item, nil
}

func (r *orderItemCommandRepository) DeletePermanent(ctx context.Context, orderItemID int) (bool, error) {
	res, err := r.db.Collection("order_items").DeleteOne(ctx, bson.M{"_id": orderItemID})
	if err != nil {
		return false, orderitem_errors.ErrDeleteOrderItemPermanent.WithInternal(err)
	}
	if res.DeletedCount == 0 {
		return false, orderitem_errors.ErrOrderItemNotFound
	}
	return true, nil
}

func (r *orderItemCommandRepository) DeleteOrderItemByOrderPermanent(ctx context.Context, orderID int) (bool, error) {
	_, err := r.db.Collection("order_items").DeleteMany(ctx, bson.M{"order_id": orderID})
	if err != nil {
		return false, orderitem_errors.ErrDeleteOrderItemPermanent.WithInternal(err)
	}
	return true, nil
}

func (r *orderItemCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("order_items").UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, orderitem_errors.ErrRestoreAllOrderItem.WithInternal(err)
	}
	return true, nil
}

func (r *orderItemCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("order_items").DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, orderitem_errors.ErrDeleteAllOrderPermanent.WithInternal(err)
	}
	return true, nil
}

func (r *orderItemCommandRepository) CalculateTotalPrice(ctx context.Context, orderID int) (int, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"order_id": orderID, "deleted_at": bson.M{"$eq": nil}}},
		{"$group": bson.M{
			"_id": nil,
			"total": bson.M{"$sum": bson.M{"$multiply": []interface{}{"$price", "$quantity"}}},
		}},
	}
	cursor, err := r.db.Collection("order_items").Aggregate(ctx, pipeline)
	if err != nil {
		return 0, orderitem_errors.ErrCalculateTotalPrice.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var results []struct {
		Total int64 `bson:"total"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return 0, orderitem_errors.ErrCalculateTotalPrice.WithInternal(err)
	}
	if len(results) == 0 {
		return 0, nil
	}
	return int(results[0].Total), nil
}
