// Package backfill implements the stats-writer `backfill` command: it reads
// historical OLTP documents from the per-service MongoDB databases (orders,
// order_items joined with products/categories, transactions) and materializes
// them into ClickHouse through the same batch repository used for live events.
package backfill

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/microservice-ecommerce-grpc-stats-writer/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/events"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

// backfillEventID derives a deterministic UUID per entity so re-running the
// backfill replaces the same ReplacingMergeTree key (with a newer version)
// instead of appending duplicates.
func backfillEventID(kind string, id int32) string {
	return uuid.NewSHA1(uuid.NameSpaceDNS, []byte(fmt.Sprintf("backfill:%s:%d", kind, id))).String()
}

// Backfiller reads OLTP documents and pushes them into ClickHouse.
type Backfiller struct {
	log      logger.LoggerInterface
	repo     repository.Repository
	order    *mongo.Database
	item     *mongo.Database
	product  *mongo.Database
	category *mongo.Database
	tx       *mongo.Database

	// Sumber statistik tersebar di dua bounded context yang masing-masing punya
	// server MongoDB sendiri, jadi backfiller memegang dua koneksi dan menutup
	// keduanya di Close.
	salesClient   *mongo.Client
	catalogClient *mongo.Client
}

// New opens one connection per bounded context that owns a stats source and
// returns a ready Backfiller.
func New(log logger.LoggerInterface, repo repository.Repository) (*Backfiller, error) {
	// order/order_item/transaction -> ec_sales, product/category -> ec_catalog.
	salesClient, err := database.NewMongoClientForContext(log, database.SalesDB)
	if err != nil {
		return nil, fmt.Errorf("connect mongo (%s): %w", database.SalesDB, err)
	}

	catalogClient, err := database.NewMongoClientForContext(log, database.CatalogDB)
	if err != nil {
		_ = salesClient.Disconnect(context.Background())
		return nil, fmt.Errorf("connect mongo (%s): %w", database.CatalogDB, err)
	}

	sales := database.NewMongoDatabase(salesClient, database.SalesDB)
	catalog := database.NewMongoDatabase(catalogClient, database.CatalogDB)

	return &Backfiller{
		log:           log,
		repo:          repo,
		order:         sales,
		item:          sales,
		product:       catalog,
		category:      catalog,
		tx:            sales,
		salesClient:   salesClient,
		catalogClient: catalogClient,
	}, nil
}

// Close disconnects both context connections.
func (b *Backfiller) Close() {
	ctx := context.Background()

	if b.salesClient != nil {
		if err := b.salesClient.Disconnect(ctx); err != nil {
			b.log.Error("Failed to close MongoDB connection",
				zap.String("db", database.SalesDB), zap.Error(err))
		}
	}
	if b.catalogClient != nil {
		if err := b.catalogClient.Disconnect(ctx); err != nil {
			b.log.Error("Failed to close MongoDB connection",
				zap.String("db", database.CatalogDB), zap.Error(err))
		}
	}
}

// Run streams all stats sources into ClickHouse. The event version is the
// backfill run timestamp so re-running supersedes previous rows.
func (b *Backfiller) Run(ctx context.Context) error {
	version := uint64(time.Now().Unix())
	counts := map[string]int{}

	if err := b.backfillOrders(ctx, version, counts); err != nil {
		return err
	}
	if err := b.backfillOrderItems(ctx, version, counts); err != nil {
		return err
	}
	if err := b.backfillTransactions(ctx, version, counts); err != nil {
		return err
	}

	if err := b.repo.Flush(ctx); err != nil {
		return fmt.Errorf("flush backfill batches: %w", err)
	}

	b.log.Info("backfill complete",
		zap.Int("orders", counts["order"]),
		zap.Int("order_items", counts["order_item"]),
		zap.Int("transactions", counts["transaction"]),
	)
	return nil
}

func findAll(ctx context.Context, db *mongo.Database, coll string, filter bson.M) ([]bson.M, error) {
	cursor, err := db.Collection(coll).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var docs []bson.M
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}
	return docs, nil
}

func toInt32Field(m bson.M, key string) int32 {
	switch v := m[key].(type) {
	case int32:
		return v
	case int64:
		return int32(v)
	case float64:
		return int32(v)
	default:
		return 0
	}
}

func toStrField(m bson.M, key string) string {
	if s, ok := m[key].(string); ok {
		return s
	}
	return ""
}

func toTimeField(m bson.M, key string) time.Time {
	switch v := m[key].(type) {
	case time.Time:
		return v
	case bson.DateTime:
		return v.Time()
	default:
		return time.Time{}
	}
}

func (b *Backfiller) backfillOrders(ctx context.Context, version uint64, counts map[string]int) error {
	docs, err := findAll(ctx, b.order, "orders", bson.M{"deleted_at": bson.M{"$eq": nil}})
	if err != nil {
		return fmt.Errorf("query orders: %w", err)
	}
	for _, d := range docs {
		orderID := toInt32Field(d, "_id")
		event := events.OrderEvent{
			OrderID:    orderID,
			UserID:     toInt32Field(d, "user_id"),
			MerchantID: toInt32Field(d, "merchant_id"),
			TotalPrice: toInt32Field(d, "total_price"),
			Status:     "created",
			EventTime:  toTimeField(d, "created_at").UTC().Format(time.RFC3339),
		}
		if err := b.repo.InsertOrderEvent(ctx, backfillEventID("order", orderID), version, event); err != nil {
			return fmt.Errorf("insert order %d: %w", orderID, err)
		}
		counts["order"]++
	}
	return nil
}

func (b *Backfiller) backfillOrderItems(ctx context.Context, version uint64, counts map[string]int) error {
	orderDocs, err := findAll(ctx, b.order, "orders", bson.M{"deleted_at": bson.M{"$eq": nil}})
	if err != nil {
		return fmt.Errorf("query order merchant map: %w", err)
	}
	orderMerchant := map[int32]int32{}
	for _, d := range orderDocs {
		orderMerchant[toInt32Field(d, "_id")] = toInt32Field(d, "merchant_id")
	}

	productDocs, err := findAll(ctx, b.product, "products", bson.M{"deleted_at": bson.M{"$eq": nil}})
	if err != nil {
		return fmt.Errorf("query product category map: %w", err)
	}
	productCategory := map[int32]int32{}
	for _, d := range productDocs {
		productCategory[toInt32Field(d, "_id")] = toInt32Field(d, "category_id")
	}

	categoryDocs, err := findAll(ctx, b.category, "categories", bson.M{"deleted_at": bson.M{"$eq": nil}})
	if err != nil {
		return fmt.Errorf("query category name map: %w", err)
	}
	categoryName := map[int32]string{}
	for _, d := range categoryDocs {
		categoryName[toInt32Field(d, "_id")] = toStrField(d, "name")
	}

	itemDocs, err := findAll(ctx, b.item, "order_items", bson.M{"deleted_at": bson.M{"$eq": nil}})
	if err != nil {
		return fmt.Errorf("query order items: %w", err)
	}
	for _, d := range itemDocs {
		itemID := toInt32Field(d, "_id")
		orderID := toInt32Field(d, "order_id")
		productID := toInt32Field(d, "product_id")
		catID := productCategory[productID]
		event := events.OrderItemEvent{
			OrderItemID:  itemID,
			OrderID:      orderID,
			MerchantID:   orderMerchant[orderID],
			ProductID:    productID,
			CategoryID:   catID,
			CategoryName: categoryName[catID],
			Quantity:     toInt32Field(d, "quantity"),
			Price:        toInt32Field(d, "price"),
			EventTime:    toTimeField(d, "created_at").UTC().Format(time.RFC3339),
		}
		if err := b.repo.InsertOrderItemEvent(ctx, backfillEventID("order_item", itemID), version, event); err != nil {
			return fmt.Errorf("insert order item %d: %w", itemID, err)
		}
		counts["order_item"]++
	}
	return nil
}

func (b *Backfiller) backfillTransactions(ctx context.Context, version uint64, counts map[string]int) error {
	docs, err := findAll(ctx, b.tx, "transactions", bson.M{"deleted_at": bson.M{"$eq": nil}})
	if err != nil {
		return fmt.Errorf("query transactions: %w", err)
	}
	for _, d := range docs {
		txID := toInt32Field(d, "_id")
		event := events.TransactionEvent{
			TransactionID: txID,
			OrderID:       toInt32Field(d, "order_id"),
			MerchantID:    toInt32Field(d, "merchant_id"),
			PaymentMethod: toStrField(d, "payment_method"),
			Amount:        toInt32Field(d, "amount"),
			Status:        toStrField(d, "payment_status"),
			EventTime:     toTimeField(d, "created_at").UTC().Format(time.RFC3339),
		}
		if err := b.repo.InsertTransactionEvent(ctx, backfillEventID("transaction", txID), version, event); err != nil {
			return fmt.Errorf("insert transaction %d: %w", txID, err)
		}
		counts["transaction"]++
	}
	return nil
}
