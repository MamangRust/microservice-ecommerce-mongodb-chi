package database

import (
	"context"
	"fmt"

	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
)

// indexPlan mengelompokkan koleksi-koleksi milik satu database (bounded context)
// beserta model index-nya.
type indexPlan struct {
	db    string
	index map[string][]mongo.IndexModel
}

// Index plan per bounded context. Penempatan koleksi mengikuti mapping yang
// disepakati (lihat pkg/database/names.go):
//   - ec_identity   : auth, user, role
//   - ec_catalog    : category, product (+ product_stock_adjustments)
//   - ec_merchant   : merchant, merchant_award, merchant_business,
//                     merchant_detail, merchant_policy
//   - ec_sales      : order, order_item, transaction
//   - ec_experience : cart, shipping_address, banner, slider, review, review_detail
//   - ec_email      : outbox + consumer inbox (dikonsumsi email service)
var indexPlans = []indexPlan{
	{
		db: IdentityDB,
		index: map[string][]mongo.IndexModel{
			"users": {
				{Keys: bson.D{{Key: "email", Value: 1}}, Options: indexUnique()},
				{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
			},
			"roles": {
				{Keys: bson.D{{Key: "role_name", Value: 1}}, Options: indexUnique()},
				{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
			},
			"user_roles": {
				{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "role_id", Value: 1}}, Options: indexUnique()},
				{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
			},
			"refresh_tokens": {
				{Keys: bson.D{{Key: "token", Value: 1}}, Options: indexUnique()},
				{Keys: bson.D{{Key: "user_id", Value: 1}}},
			},
			"reset_tokens": {
				{Keys: bson.D{{Key: "token", Value: 1}}, Options: indexUnique()},
				{Keys: bson.D{{Key: "user_id", Value: 1}}},
			},
		},
	},
	{
		db: CatalogDB,
		index: map[string][]mongo.IndexModel{
			"categories": {
				{Keys: bson.D{{Key: "name", Value: 1}}, Options: indexUnique()},
				{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
			},
			"products": {
				{Keys: bson.D{{Key: "merchant_id", Value: 1}}},
				{Keys: bson.D{{Key: "category_id", Value: 1}}},
				{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
			},
			"product_stock_adjustments": {
				{Keys: bson.D{{Key: "product_id", Value: 1}}},
			},
		},
	},
	{
		db: MerchantDB,
		index: map[string][]mongo.IndexModel{
			"merchants": {
				{Keys: bson.D{{Key: "user_id", Value: 1}}},
				{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
				{Keys: bson.D{{Key: "status", Value: 1}}},
			},
			"merchant_documents": {
				{Keys: bson.D{{Key: "merchant_id", Value: 1}}},
				{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
			},
			"merchant_details": {
				{Keys: bson.D{{Key: "merchant_id", Value: 1}}},
				{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
			},
			"merchant_social_media_links": {
				{Keys: bson.D{{Key: "merchant_detail_id", Value: 1}}},
				{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
			},
			"merchant_business_information": {
				{Keys: bson.D{{Key: "merchant_id", Value: 1}}},
				{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
			},
			"merchant_policies": {
				{Keys: bson.D{{Key: "merchant_id", Value: 1}}},
				{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
			},
			"merchant_certifications_and_awards": {
				{Keys: bson.D{{Key: "merchant_id", Value: 1}}},
				{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
			},
		},
	},
	{
		db: SalesDB,
		index: map[string][]mongo.IndexModel{
			"orders": {
				{Keys: bson.D{{Key: "user_id", Value: 1}}},
				{Keys: bson.D{{Key: "merchant_id", Value: 1}}},
				{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
			},
			"order_items": {
				{Keys: bson.D{{Key: "order_id", Value: 1}}},
				{Keys: bson.D{{Key: "product_id", Value: 1}}},
				{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
			},
			"order_stock_reservations": {
				{Keys: bson.D{{Key: "order_id", Value: 1}}},
				{Keys: bson.D{{Key: "product_id", Value: 1}}},
			},
			"transactions": {
				{Keys: bson.D{{Key: "order_id", Value: 1}}},
				{Keys: bson.D{{Key: "merchant_id", Value: 1}}},
				{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
			},
		},
	},
	{
		db: ExperienceDB,
		index: map[string][]mongo.IndexModel{
			"banners": {
				{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
			},
			"sliders": {
				{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
			},
			"shipping_addresses": {
				{Keys: bson.D{{Key: "order_id", Value: 1}}},
				{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
			},
			"carts": {
				{Keys: bson.D{{Key: "user_id", Value: 1}}},
				{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
			},
			"reviews": {
				{Keys: bson.D{{Key: "user_id", Value: 1}}},
				{Keys: bson.D{{Key: "product_id", Value: 1}}},
				{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
			},
			"review_details": {
				{Keys: bson.D{{Key: "review_id", Value: 1}}},
				{Keys: bson.D{{Key: "deleted_at", Value: 1}}},
			},
		},
	},
	{
		db: EmailDB,
		index: map[string][]mongo.IndexModel{
			"outbox_events": {
				{Keys: bson.D{{Key: "status", Value: 1}, {Key: "next_attempt_at", Value: 1}}},
			},
			"consumer_inbox": {
				{Keys: bson.D{{Key: "consumer_name", Value: 1}, {Key: "event_key", Value: 1}}, Options: indexUnique()},
				{Keys: bson.D{{Key: "status", Value: 1}, {Key: "lease_until", Value: 1}}},
			},
		},
	},
}

// indexUnique returns index options that enforce a unique constraint.
func indexUnique() *options.IndexOptionsBuilder {
	return options.Index().SetUnique(true)
}

// EnsureIndexes creates all indexes for every bounded-context database.
//
// Setiap bounded context punya server MongoDB sendiri ("1 database = 1
// MongoDB"), jadi fungsi ini membuka satu koneksi per context lewat
// NewMongoClientForContext dan menutupnya lagi setelah context itu selesai.
// Tidak ada lagi satu client bersama untuk enam database.
//
// Counters TIDAK lagi diinisialisasi di sini: GetNextID memakai $inc + upsert,
// sehingga counter tercipta otomatis saat ID pertama diminta — tidak perlu
// predeclaration.
//
// Pemanggilan ini hanya dilakukan oleh service/migrate (bukan saat startup
// masing-masing service), supaya pembuatan index terpusat dan idempoten.
func EnsureIndexes(log logger.LoggerInterface) error {
	for _, plan := range indexPlans {
		if err := ensureIndexesForContext(log, plan); err != nil {
			return err
		}
	}

	log.Info("MongoDB indexes created successfully across all bounded-context databases")
	return nil
}

// ensureIndexesForContext connects to the server owning plan.db, creates the
// indexes, and always disconnects before returning.
func ensureIndexesForContext(log logger.LoggerInterface, plan indexPlan) error {
	client, err := NewMongoClientForContext(log, plan.db)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB for %s: %w", plan.db, err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Error("Failed to close MongoDB connection",
				zap.String("db", plan.db),
				zap.Error(err))
		}
	}()

	return createIndexesForPlan(context.Background(), log, client.Database(plan.db), plan)
}

// EnsureIndexesOnClient creates the same indexes using a single shared client.
//
// Dipakai oleh harness integration test, yang menjalankan keenam bounded context
// di satu container MongoDB sehingga tidak perlu (dan tidak bisa) membuka enam
// koneksi. Topologi produksi memakai EnsureIndexes.
func EnsureIndexesOnClient(log logger.LoggerInterface, client *mongo.Client) error {
	ctx := context.Background()

	for _, plan := range indexPlans {
		if err := createIndexesForPlan(ctx, log, client.Database(plan.db), plan); err != nil {
			return err
		}
	}

	log.Info("MongoDB indexes created successfully across all bounded-context databases")
	return nil
}

func createIndexesForPlan(
	ctx context.Context,
	log logger.LoggerInterface,
	db *mongo.Database,
	plan indexPlan,
) error {
	for collName, models := range plan.index {
		coll := db.Collection(collName)
		if _, err := coll.Indexes().CreateMany(ctx, models); err != nil {
			log.Error("Failed to create indexes",
				zap.String("db", plan.db),
				zap.String("collection", collName),
				zap.Error(err))
			return err
		}
		log.Info("Created indexes",
			zap.String("db", plan.db),
			zap.String("collection", collName),
			zap.Int("count", len(models)))
	}

	return nil
}
