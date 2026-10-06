package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type stockReservationRepository struct {
	db        *mongo.Database // ec_sales: order_stock_reservations
	catalogDB *mongo.Database // ec_catalog: product_stock_adjustments
}

func NewStockReservationRepository(db, catalogDB *mongo.Database) StockReservationRepository {
	return &stockReservationRepository{db: db, catalogDB: catalogDB}
}

func (r *stockReservationRepository) GetByOrder(ctx context.Context, orderID int) ([]*models.OrderStockReservation, error) {
	cursor, err := r.db.Collection("order_stock_reservations").Find(ctx, bson.M{"order_id": orderID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var reservations []*models.OrderStockReservation
	if err := cursor.All(ctx, &reservations); err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *stockReservationRepository) Upsert(ctx context.Context, orderID, productID, quantity int) (*models.OrderStockReservation, error) {
	coll := r.db.Collection("order_stock_reservations")
	now := time.Now()

	var existing models.OrderStockReservation
	err := coll.FindOne(ctx, bson.M{"order_id": orderID, "product_id": productID}).Decode(&existing)
	if err == mongo.ErrNoDocuments {
		id, err := database.GetNextID(ctx, r.db, "order_stock_reservations")
		if err != nil {
			return nil, err
		}
		reservation := &models.OrderStockReservation{
			ReservationID: id,
			OrderID:       int32(orderID),
			ProductID:     int32(productID),
			Quantity:      int32(quantity),
			Status:        "reserved",
			CreatedAt:     &now,
			UpdatedAt:     &now,
		}
		_, err = coll.InsertOne(ctx, reservation)
		if err != nil {
			return nil, err
		}
		return reservation, nil
	} else if err != nil {
		return nil, err
	}

	existing.Quantity = int32(quantity)
	existing.Status = "reserved"
	existing.UpdatedAt = &now
	_, err = coll.ReplaceOne(ctx, bson.M{"_id": existing.ReservationID}, existing)
	if err != nil {
		return nil, err
	}
	return &existing, nil
}

func (r *stockReservationRepository) UpdateQuantity(ctx context.Context, orderID, productID, quantity int) (*models.OrderStockReservation, error) {
	coll := r.db.Collection("order_stock_reservations")
	var reservation models.OrderStockReservation
	err := coll.FindOne(ctx, bson.M{"order_id": orderID, "product_id": productID}).Decode(&reservation)
	if err != nil {
		return nil, err
	}
	reservation.Quantity = int32(quantity)
	now := time.Now()
	reservation.UpdatedAt = &now
	_, err = coll.ReplaceOne(ctx, bson.M{"_id": reservation.ReservationID}, reservation)
	return &reservation, err
}

func (r *stockReservationRepository) Release(ctx context.Context, orderID, productID int) (*models.OrderStockReservation, error) {
	return r.updateStatus(ctx, orderID, productID, "reserved", "released")
}

func (r *stockReservationRepository) Reserve(ctx context.Context, orderID, productID int) (*models.OrderStockReservation, error) {
	return r.updateStatus(ctx, orderID, productID, "released", "reserved")
}

func (r *stockReservationRepository) updateStatus(ctx context.Context, orderID, productID int, fromStatus, toStatus string) (*models.OrderStockReservation, error) {
	coll := r.db.Collection("order_stock_reservations")
	now := time.Now()
	res, err := coll.UpdateOne(ctx,
		bson.M{"order_id": orderID, "product_id": productID, "status": fromStatus},
		bson.M{"$set": bson.M{"status": toStatus, "updated_at": now}},
	)
	if err != nil {
		return nil, err
	}
	if res.MatchedCount == 0 {
		return nil, mongo.ErrNoDocuments
	}
	var reservation models.OrderStockReservation
	if err := coll.FindOne(ctx, bson.M{"order_id": orderID, "product_id": productID}).Decode(&reservation); err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (r *stockReservationRepository) ordersFilter(ctx context.Context, trashed bool) bson.M {
	// Find orders that are trashed (deleted_at != nil) or active (deleted_at == nil)
	delFilter := bson.M{"$ne": nil}
	if !trashed {
		delFilter = bson.M{"$eq": nil}
	}
	cursor, err := r.db.Collection("orders").Find(ctx, bson.M{"deleted_at": delFilter}, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return bson.M{"$in": []int32{}}
	}
	defer cursor.Close(ctx)
	var orders []bson.M
	if err := cursor.All(ctx, &orders); err != nil {
		return bson.M{"$in": []int32{}}
	}
	ids := make([]int32, 0, len(orders))
	for _, o := range orders {
		ids = append(ids, bsonMToInt32Order(o, "_id"))
	}
	return bson.M{"$in": ids}
}

func (r *stockReservationRepository) GetReservedForTrashedOrders(ctx context.Context) ([]*models.OrderStockReservation, error) {
	return r.findByStatusAndOrders(ctx, "reserved", r.ordersFilter(ctx, true))
}

func (r *stockReservationRepository) GetReleasedForTrashedOrders(ctx context.Context) ([]*models.OrderStockReservation, error) {
	return r.findByStatusAndOrders(ctx, "released", r.ordersFilter(ctx, true))
}

func (r *stockReservationRepository) GetReleasedForActiveOrders(ctx context.Context) ([]*models.OrderStockReservation, error) {
	return r.findByStatusAndOrders(ctx, "released", r.ordersFilter(ctx, false))
}

func (r *stockReservationRepository) findByStatusAndOrders(ctx context.Context, status string, orderIDFilter bson.M) ([]*models.OrderStockReservation, error) {
	filter := bson.M{"status": status, "order_id": orderIDFilter}
	cursor, err := r.db.Collection("order_stock_reservations").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var reservations []*models.OrderStockReservation
	if err := cursor.All(ctx, &reservations); err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *stockReservationRepository) DeleteByOrder(ctx context.Context, orderID int) error {
	_, err := r.db.Collection("order_stock_reservations").DeleteMany(ctx, bson.M{"order_id": orderID})
	return err
}

func (r *stockReservationRepository) DeleteByOrderProduct(ctx context.Context, orderID, productID int) error {
	_, err := r.db.Collection("order_stock_reservations").DeleteMany(ctx, bson.M{"order_id": orderID, "product_id": productID})
	return err
}

func (r *stockReservationRepository) DeleteAllForTrashedOrders(ctx context.Context) error {
	_, err := r.db.Collection("order_stock_reservations").DeleteMany(ctx, bson.M{"order_id": r.ordersFilter(ctx, true)})
	return err
}

func (r *stockReservationRepository) DeleteOldReleasedReservations(ctx context.Context, cutoff time.Time) (int64, error) {
	result, err := r.db.Collection("order_stock_reservations").DeleteMany(ctx,
		bson.M{"status": "released", "updated_at": bson.M{"$lt": cutoff}, "order_id": r.ordersFilter(ctx, true)},
	)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

func (r *stockReservationRepository) DeleteOldProductStockAdjustments(ctx context.Context, cutoff time.Time) (int64, error) {
	result, err := r.catalogDB.Collection("product_stock_adjustments").DeleteMany(ctx,
		bson.M{"created_at": bson.M{"$lt": cutoff}},
	)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
