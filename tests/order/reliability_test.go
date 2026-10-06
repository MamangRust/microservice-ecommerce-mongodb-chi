package order_test

import (
	"context"
	pbproduct "github.com/MamangRust/microservice-ecommerce-grpc-pb/product"
	"time"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (s *OrderServiceTestSuite) orderRequest(userID, merchID, prodID int) *requests.CreateOrderRequest {
	return &requests.CreateOrderRequest{
		UserID:     userID,
		MerchantID: merchID,
		TotalPrice: 20000,
		Items: []requests.CreateOrderItemRequest{
			{ProductID: prodID, Quantity: 1, Price: 10000},
		},
		ShippingAddress: requests.CreateShippingAddressRequest{
			Alamat:         "Test Address",
			Provinsi:       "West Java",
			Kota:           "Bandung",
			Courier:        "JNE",
			ShippingMethod: "REG",
			ShippingCost:   10000,
			Negara:         "Indonesia",
		},
	}
}

func (s *OrderServiceTestSuite) stockOf(ctx context.Context, prodID int) int32 {
	prodRes, err := pbproduct.NewProductQueryServiceClient(s.Conns["product"]).FindById(ctx, &pbproduct.FindByIdProductRequest{Id: int32(prodID)})
	s.Require().NoError(err)
	s.Require().NotNil(prodRes.Data)
	return prodRes.Data.CountInStock
}

// TestFailureInjectionInsufficientStock verifies that creating an order with a
// quantity beyond the available stock fails, rolls back any partial reservation,
// and leaves inventory untouched.
func (s *OrderServiceTestSuite) TestFailureInjectionInsufficientStock() {
	ctx := context.Background()

	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	prodID := s.SeedProduct(ctx, merchID, catID)

	s.Require().Equal(int32(100), s.stockOf(ctx, prodID))

	// Demand far more than the 100 units in stock.
	req := s.orderRequest(userID, merchID, prodID)
	req.Items[0].Quantity = 1000

	_, err := s.svc.OrderCommand.Create(ctx, req)
	s.Require().Error(err, "order creation beyond stock must fail")

	s.Equal(int32(100), s.stockOf(ctx, prodID), "stock must be untouched after a failed create")

	db, err := s.MongoDB()
	s.Require().NoError(err)

	// No orphan reservation rows (or orders) may remain for the failed create.
	orderIDs := s.userOrderIDs(ctx, db, userID)
	reservations, err := db.Collection("order_stock_reservations").CountDocuments(ctx,
		bson.M{"order_id": bson.M{"$in": orderIDs}})
	s.Require().NoError(err)
	s.Zero(reservations, "no reservations may survive a rolled-back create")

	orders, err := db.Collection("orders").CountDocuments(ctx, bson.M{"user_id": int32(userID)})
	s.Require().NoError(err)
	s.Zero(orders, "no order rows may survive a rolled-back create")
}

// TestReconcileDriftAndCleanupIdempotency verifies durable reconciliation repairs
// reservation/order drift in both directions (with correct stock deltas) and that
// the retention policy purges only rows older than the cutoff.
func (s *OrderServiceTestSuite) TestReconcileDriftAndCleanupIdempotency() {
	ctx := context.Background()

	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	prodID := s.SeedProduct(ctx, merchID, catID)

	db, err := s.MongoDB()
	s.Require().NoError(err)

	created, err := s.svc.OrderCommand.Create(ctx, s.orderRequest(userID, merchID, prodID))
	s.Require().NoError(err)
	orderID := int(created.OrderID)
	s.Equal(int32(99), s.stockOf(ctx, prodID), "create must reserve one unit of stock")

	reservationStatus := func() string {
		var doc struct {
			Status string `bson:"status"`
		}
		err := db.Collection("order_stock_reservations").FindOne(ctx,
			bson.M{"order_id": int32(orderID), "product_id": int32(prodID)}).Decode(&doc)
		s.Require().NoError(err)
		return doc.Status
	}
	s.Equal("reserved", reservationStatus(), "create must reserve stock")

	// Drift A: reservation wrongly marked released while the order is still active.
	// Durable reconciliation must re-reserve it (decrement stock again).
	_, err = db.Collection("order_stock_reservations").UpdateOne(ctx,
		bson.M{"order_id": int32(orderID), "product_id": int32(prodID)},
		bson.M{"$set": bson.M{"status": "released"}})
	s.Require().NoError(err)

	res, err := s.svc.OrderCommand.ReconcileStockReservations(ctx)
	s.Require().NoError(err)
	s.GreaterOrEqual(res.ReReserved, 1)
	s.Equal("reserved", reservationStatus(), "reconciliation must re-reserve stock for active orders")
	s.Equal(int32(98), s.stockOf(ctx, prodID), "re-reconciliation must decrement stock again")

	// Drift B: trash releases stock; simulate a reservation stuck as reserved.
	_, err = s.svc.OrderCommand.Trash(ctx, orderID)
	s.Require().NoError(err)
	s.Equal("released", reservationStatus(), "trash must release stock")
	s.Equal(int32(99), s.stockOf(ctx, prodID), "trash must return stock to inventory")

	_, err = db.Collection("order_stock_reservations").UpdateOne(ctx,
		bson.M{"order_id": int32(orderID), "product_id": int32(prodID)},
		bson.M{"$set": bson.M{"status": "reserved"}})
	s.Require().NoError(err)

	res2, err := s.svc.OrderCommand.ReconcileStockReservations(ctx)
	s.Require().NoError(err)
	s.GreaterOrEqual(res2.Released, 1)
	s.Equal("released", reservationStatus(), "reconciliation must release stock for trashed orders")
	s.Equal(int32(100), s.stockOf(ctx, prodID), "reconciliation must restore released stock")

	// Cleanup: age the released reservation and insert an old adjustment row,
	// then confirm the retention policy purges only rows older than the cutoff.
	_, err = db.Collection("order_stock_reservations").UpdateOne(ctx,
		bson.M{"order_id": int32(orderID), "product_id": int32(prodID)},
		bson.M{"$set": bson.M{"updated_at": time.Now().AddDate(0, 0, -10)}})
	s.Require().NoError(err)

	// product_stock_adjustments now lives in the catalog DB (ec_catalog) after
	// the 6-DB split; inserting it into the order DB leaves the retention job
	// with nothing to purge.
	catalogDB := s.DBForService("product")
	_, err = catalogDB.Collection("product_stock_adjustments").InsertOne(ctx, bson.M{
		"operation_id": "test-old-adjustment",
		"product_id":   int32(prodID),
		"delta":        int32(1),
		"created_at":   time.Now().AddDate(0, 0, -10),
	})
	s.Require().NoError(err)

	cleanup, err := s.svc.OrderCommand.CleanupIdempotencyRecords(ctx, 1)
	s.Require().NoError(err)
	s.GreaterOrEqual(cleanup.ReleasedReservationsRemoved, int64(1))
	s.GreaterOrEqual(cleanup.AdjustmentsRemoved, int64(1))

	// The aged trashed reservation is gone.
	stale, err := db.Collection("order_stock_reservations").CountDocuments(ctx,
		bson.M{"order_id": int32(orderID), "product_id": int32(prodID)})
	s.Require().NoError(err)
	s.Zero(stale, "aged released reservation must be purged")

	// Fresh reservations created after cleanup must survive the retention policy.
	freshOrder, err := s.svc.OrderCommand.Create(ctx, s.orderRequest(userID, merchID, prodID))
	s.Require().NoError(err)
	freshOrderID := int(freshOrder.OrderID)

	freshReservations, err := db.Collection("order_stock_reservations").CountDocuments(ctx,
		bson.M{"order_id": int32(freshOrderID), "product_id": int32(prodID), "status": "reserved"})
	s.Require().NoError(err)
	s.Equal(int64(1), freshReservations, "fresh reservation must survive the retention cleanup")
}

// userOrderIDs returns the order ids belonging to a user.
func (s *OrderServiceTestSuite) userOrderIDs(ctx context.Context, db *mongo.Database, userID int) []int32 {
	cursor, err := db.Collection("orders").Find(ctx, bson.M{"user_id": int32(userID)})
	s.Require().NoError(err)
	defer cursor.Close(ctx)

	var docs []struct {
		ID int32 `bson:"_id"`
	}
	s.Require().NoError(cursor.All(ctx, &docs))

	ids := make([]int32, 0, len(docs))
	for _, d := range docs {
		ids = append(ids, d.ID)
	}
	return ids
}
