// Package stats_writer_test exercises the ecommerce stats-writer backfill end
// to end: it seeds the owning MongoDB collections (catalog + sales contexts),
// runs the real backfiller — which opens Mongo itself and reads the OLTP
// documents directly — against a real ClickHouse, and asserts the materialized
// rows.
package stats_writer_test

import (
	"testing"
	"time"

	chDriver "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/MamangRust/microservice-ecommerce-grpc-stats-writer/backfill"
	statsrepo "github.com/MamangRust/microservice-ecommerce-grpc-stats-writer/repository"
	pkgclickhouse "github.com/MamangRust/microservice-ecommerce-pkg/clickhouse"
	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	tests "github.com/MamangRust/microservice-ecommerce-test"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	categoryAName = "Cat A"
	categoryBName = "Cat B"

	productOnePrice = 1000
	productTwoPrice = 2000

	orderOneTotal = 3000
	orderTwoTotal = 5000

	txOneAmount = 3000
	txTwoAmount = 5000
)

type BackfillSuite struct {
	tests.BaseTestSuite

	chConn chDriver.Conn
	repo   statsrepo.Repository

	orderOne int32
	orderTwo int32
}

func (s *BackfillSuite) SetupSuite() {
	// Orders, order items and transactions live in the sales context, which is
	// the suite's primary database; catalog data is seeded cross-context through
	// Coll/DBForCollection.
	s.PrimaryContext = database.SalesDB
	s.BaseTestSuite.SetupSuite()
	s.SetupStatsEnv()

	conn, err := pkgclickhouse.NewClient(s.Log)
	s.Require().NoError(err)
	s.Require().NoError(pkgclickhouse.ApplySchema(s.Ctx, conn, s.Log))

	s.chConn = conn
	s.repo = statsrepo.NewClickhouseRepository(conn, s.Log)

	s.seedMongo()
}

func (s *BackfillSuite) TearDownSuite() {
	if s.repo != nil {
		_ = s.repo.Close()
	}
	if s.chConn != nil {
		_ = s.chConn.Close()
	}
	s.BaseTestSuite.TearDownSuite()
}

// seedMongo writes the fixture rows into the database that owns each
// collection: categories/products into the catalog context, and
// orders/order_items/transactions into the sales context. The backfiller opens
// its own Mongo connection and reads these collections directly, so this is
// purely a data setup step.
func (s *BackfillSuite) seedMongo() {
	ctx := s.Ctx
	now := time.Now().UTC()

	// Catalog context: categories → products.
	catColl := s.Coll("categories")
	categoryA, err := database.GetNextID(ctx, s.DBForCollection("categories"), "categories")
	s.Require().NoError(err)
	categoryB, err := database.GetNextID(ctx, s.DBForCollection("categories"), "categories")
	s.Require().NoError(err)
	_, err = catColl.InsertMany(ctx, []interface{}{
		bson.M{"_id": categoryA, "name": categoryAName, "description": "a", "deleted_at": nil, "created_at": now, "updated_at": now},
		bson.M{"_id": categoryB, "name": categoryBName, "description": "b", "deleted_at": nil, "created_at": now, "updated_at": now},
	})
	s.Require().NoError(err)

	prodColl := s.Coll("products")
	productOne, err := database.GetNextID(ctx, s.DBForCollection("products"), "products")
	s.Require().NoError(err)
	productTwo, err := database.GetNextID(ctx, s.DBForCollection("products"), "products")
	s.Require().NoError(err)
	_, err = prodColl.InsertMany(ctx, []interface{}{
		bson.M{"_id": productOne, "merchant_id": int32(1), "category_id": categoryA, "name": "P1", "price": int32(productOnePrice), "count_in_stock": int32(100), "deleted_at": nil, "created_at": now, "updated_at": now},
		bson.M{"_id": productTwo, "merchant_id": int32(1), "category_id": categoryB, "name": "P2", "price": int32(productTwoPrice), "count_in_stock": int32(100), "deleted_at": nil, "created_at": now, "updated_at": now},
	})
	s.Require().NoError(err)

	// Sales context: two orders, one per merchant.
	orderColl := s.Coll("orders")
	s.orderOne, err = database.GetNextID(ctx, s.DBForCollection("orders"), "orders")
	s.Require().NoError(err)
	s.orderTwo, err = database.GetNextID(ctx, s.DBForCollection("orders"), "orders")
	s.Require().NoError(err)
	_, err = orderColl.InsertMany(ctx, []interface{}{
		bson.M{"_id": s.orderOne, "user_id": int32(1), "merchant_id": int32(1), "total_price": int32(orderOneTotal), "deleted_at": nil, "created_at": now, "updated_at": now},
		bson.M{"_id": s.orderTwo, "user_id": int32(2), "merchant_id": int32(2), "total_price": int32(orderTwoTotal), "deleted_at": nil, "created_at": now, "updated_at": now},
	})
	s.Require().NoError(err)

	// Sales context: one order item per order, pointing at the two products so
	// both catalog categories are exercised.
	itemColl := s.Coll("order_items")
	itemOne, err := database.GetNextID(ctx, s.DBForCollection("order_items"), "order_items")
	s.Require().NoError(err)
	itemTwo, err := database.GetNextID(ctx, s.DBForCollection("order_items"), "order_items")
	s.Require().NoError(err)
	_, err = itemColl.InsertMany(ctx, []interface{}{
		bson.M{"_id": itemOne, "order_id": s.orderOne, "product_id": productOne, "quantity": int32(2), "price": int32(productOnePrice), "deleted_at": nil, "created_at": now, "updated_at": now},
		bson.M{"_id": itemTwo, "order_id": s.orderTwo, "product_id": productTwo, "quantity": int32(1), "price": int32(productTwoPrice), "deleted_at": nil, "created_at": now, "updated_at": now},
	})
	s.Require().NoError(err)

	// Sales context: one successful transaction per order.
	txColl := s.Coll("transactions")
	txOne, err := database.GetNextID(ctx, s.DBForCollection("transactions"), "transactions")
	s.Require().NoError(err)
	txTwo, err := database.GetNextID(ctx, s.DBForCollection("transactions"), "transactions")
	s.Require().NoError(err)
	_, err = txColl.InsertMany(ctx, []interface{}{
		bson.M{"_id": txOne, "order_id": s.orderOne, "merchant_id": int32(1), "payment_method": "cash", "amount": int32(txOneAmount), "payment_status": "success", "deleted_at": nil, "created_at": now, "updated_at": now},
		bson.M{"_id": txTwo, "order_id": s.orderTwo, "merchant_id": int32(2), "payment_method": "transfer", "amount": int32(txTwoAmount), "payment_status": "success", "deleted_at": nil, "created_at": now, "updated_at": now},
	})
	s.Require().NoError(err)
}

func (s *BackfillSuite) TestBackfillMaterializesEverySource() {
	bf, err := backfill.New(s.Log, s.repo)
	s.Require().NoError(err)
	s.Require().NoError(bf.Run(s.Ctx))

	s.Equal(uint64(2), s.count(`SELECT count() FROM order_events`))
	s.Equal(int64(orderOneTotal+orderTwoTotal), s.sum(`SELECT sum(total_price) FROM order_events`))

	s.Equal(uint64(2), s.count(`SELECT count() FROM order_item_events`))
	s.Equal(uint64(0), s.count(
		`SELECT count() FROM order_item_events WHERE quantity <= 0 OR price <= 0`),
		"item quantity and price must be carried over")

	// Category must be denormalized from the catalog onto the order item.
	s.Equal(uint64(1), s.count(`SELECT count() FROM order_item_events WHERE category_name = 'Cat A'`))
	s.Equal(uint64(1), s.count(`SELECT count() FROM order_item_events WHERE category_name = 'Cat B'`))

	// Merchant must be recovered from the owning order.
	s.Equal(uint64(1), s.count(`SELECT count() FROM order_item_events WHERE merchant_id = 1`))
	s.Equal(uint64(1), s.count(`SELECT count() FROM order_item_events WHERE merchant_id = 2`))

	s.Equal(uint64(2), s.count(`SELECT count() FROM transaction_events`))
	s.Equal(uint64(2), s.count(`SELECT count() FROM transaction_events WHERE status = 'success'`))
	s.Equal(int64(txOneAmount+txTwoAmount), s.sum(`SELECT sum(amount) FROM transaction_events`))
	s.Equal(uint64(1), s.count(`SELECT count() FROM transaction_events WHERE payment_method = 'cash'`))
	s.Equal(uint64(1), s.count(`SELECT count() FROM transaction_events WHERE payment_method = 'transfer'`))
}

func (s *BackfillSuite) count(query string) uint64 {
	var n uint64
	s.Require().NoError(s.chConn.QueryRow(s.Ctx, query).Scan(&n))
	return n
}

func (s *BackfillSuite) sum(query string) int64 {
	var n int64
	s.Require().NoError(s.chConn.QueryRow(s.Ctx, query).Scan(&n))
	return n
}

func TestBackfillSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test requires Docker")
	}
	suite.Run(t, new(BackfillSuite))
}
