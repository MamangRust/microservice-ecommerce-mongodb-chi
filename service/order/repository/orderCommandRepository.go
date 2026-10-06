package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/order_errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type orderCommandRepository struct {
	db *mongo.Database
}

func NewOrderCommandRepository(db *mongo.Database) OrderCommandRepository {
	return &orderCommandRepository{db: db}
}

func (r *orderCommandRepository) Create(ctx context.Context, request *requests.CreateOrderRecordRequest) (*models.Order, error) {
	id, err := database.GetNextID(ctx, r.db, "orders")
	if err != nil {
		return nil, order_errors.ErrCreateOrder.WithInternal(err)
	}
	now := time.Now()
	order := &models.Order{
		OrderID:    id,
		MerchantID: int32(request.MerchantID),
		UserID:     int32(request.UserID),
		TotalPrice: int32(request.TotalPrice),
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}
	if _, err := r.db.Collection("orders").InsertOne(ctx, order); err != nil {
		return nil, order_errors.ErrCreateOrder.WithInternal(err)
	}
	return order, nil
}

func (r *orderCommandRepository) Update(ctx context.Context, request *requests.UpdateOrderRecordRequest) (*models.Order, error) {
	coll := r.db.Collection("orders")
	var order models.Order
	if err := coll.FindOne(ctx, bson.M{"_id": request.OrderID, "deleted_at": bson.M{"$eq": nil}}).Decode(&order); err != nil {
		return nil, order_errors.ErrOrderNotFound
	}
	order.TotalPrice = int32(request.TotalPrice)
	now := time.Now()
	order.UpdatedAt = &now
	if _, err := coll.ReplaceOne(ctx, bson.M{"_id": request.OrderID}, order); err != nil {
		return nil, order_errors.ErrUpdateOrder.WithInternal(err)
	}
	return &order, nil
}

func (r *orderCommandRepository) Trash(ctx context.Context, orderID int) (*models.Order, error) {
	coll := r.db.Collection("orders")
	now := time.Now()
	res, err := coll.UpdateOne(ctx, bson.M{"_id": orderID, "deleted_at": bson.M{"$eq": nil}}, bson.M{"$set": bson.M{"deleted_at": now}})
	if err != nil {
		return nil, order_errors.ErrTrashedOrder.WithInternal(err)
	}
	if res.MatchedCount == 0 {
		return nil, order_errors.ErrOrderNotFound
	}
	var order models.Order
	if err := coll.FindOne(ctx, bson.M{"_id": orderID}).Decode(&order); err != nil {
		return nil, order_errors.ErrTrashedOrder.WithInternal(err)
	}
	return &order, nil
}

func (r *orderCommandRepository) Restore(ctx context.Context, orderID int) (*models.Order, error) {
	coll := r.db.Collection("orders")
	res, err := coll.UpdateOne(ctx, bson.M{"_id": orderID, "deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return nil, order_errors.ErrRestoreOrder.WithInternal(err)
	}
	if res.MatchedCount == 0 {
		return nil, order_errors.ErrOrderNotFound
	}
	var order models.Order
	if err := coll.FindOne(ctx, bson.M{"_id": orderID}).Decode(&order); err != nil {
		return nil, order_errors.ErrRestoreOrder.WithInternal(err)
	}
	return &order, nil
}

func (r *orderCommandRepository) FindTrashedByID(ctx context.Context, orderID int) (*models.Order, error) {
	var order models.Order
	err := r.db.Collection("orders").FindOne(ctx, bson.M{"_id": orderID, "deleted_at": bson.M{"$ne": nil}}).Decode(&order)
	if err != nil {
		return nil, order_errors.ErrOrderNotFound
	}
	return &order, nil
}

func (r *orderCommandRepository) FindTrashed(ctx context.Context) ([]*models.Order, error) {
	cursor, err := r.db.Collection("orders").Find(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return nil, order_errors.ErrFindByTrashed.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var orders []*models.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, order_errors.ErrFindByTrashed.WithInternal(err)
	}
	return orders, nil
}

func (r *orderCommandRepository) DeletePermanent(ctx context.Context, orderID int) (bool, error) {
	// Delete child records first
	r.db.Collection("order_stock_reservations").DeleteMany(ctx, bson.M{"order_id": orderID})

	res, err := r.db.Collection("orders").DeleteOne(ctx, bson.M{"_id": orderID})
	if err != nil {
		return false, order_errors.ErrDeleteOrderPermanent.WithInternal(err)
	}
	if res.DeletedCount == 0 {
		return false, order_errors.ErrOrderNotFound
	}
	return true, nil
}

func (r *orderCommandRepository) DeletePermanentWithChildren(ctx context.Context, orderID int) (bool, error) {
	return r.DeletePermanent(ctx, orderID)
}

func (r *orderCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("orders").UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, order_errors.ErrRestoreAllOrder.WithInternal(err)
	}
	return true, nil
}

func (r *orderCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("orders").DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, order_errors.ErrDeleteAllOrderPermanent.WithInternal(err)
	}
	return true, nil
}
