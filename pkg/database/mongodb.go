package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"go.uber.org/zap"
)

// NewMongoClient connects to MongoDB using the MONGO_* env vars and returns a
// *mongo.Client. Pemilihan database dilakukan terpisah lewat NewMongoDatabase,
// supaya satu koneksi bisa dipakai untuk beberapa database bila memang perlu.
func NewMongoClient(log logger.LoggerInterface) (*mongo.Client, error) {
	return NewMongoClientWithPrefix(log, "MONGO")
}

// NewMongoClientForContext connects to the MongoDB server that owns the given
// bounded context (lihat ContextPrefix). Setiap database punya server sendiri
// ("1 database = 1 MongoDB"), jadi pemanggil yang butuh lebih dari satu
// database harus membuka lebih dari satu client.
//
// Bila key per-context tidak ada, rantai fallback di NewMongoClientWithPrefix
// jatuh ke MONGO_URI / MONGO_HOST — sehingga deployment instance-tunggal lama
// (termasuk harness integration test) tetap jalan tanpa perubahan env.
func NewMongoClientForContext(log logger.LoggerInterface, databaseName string) (*mongo.Client, error) {
	return NewMongoClientWithPrefix(log, ContextPrefix(databaseName))
}

// NewMongoClientWithPrefix connects to MongoDB configured via the given prefix
// keys (e.g. DB_EMAIL_HOST) with fallback to the base MONGO_* keys.
//
// Rantai fallback tiga lapis ini menyamakan pola dengan pointofsale dan
// payment-gateway: prefix yang env-nya kosong jatuh ke MONGO_*, bukan langsung
// ke localhost.
func NewMongoClientWithPrefix(log logger.LoggerInterface, prefix string) (*mongo.Client, error) {
	if prefix == "" {
		prefix = "MONGO"
	}

	host := viper.GetString(fmt.Sprintf("%s_HOST", prefix))
	if host == "" {
		host = viper.GetString("MONGO_HOST")
	}
	if host == "" {
		host = "localhost"
	}

	port := viper.GetString(fmt.Sprintf("%s_PORT", prefix))
	if port == "" {
		port = viper.GetString("MONGO_PORT")
	}
	if port == "" {
		port = "27017"
	}

	user := viper.GetString(fmt.Sprintf("%s_USERNAME", prefix))
	if user == "" {
		user = viper.GetString("MONGO_USERNAME")
	}

	pass := viper.GetString(fmt.Sprintf("%s_PASSWORD", prefix))
	if pass == "" {
		pass = viper.GetString("MONGO_PASSWORD")
	}

	uri := fmt.Sprintf("mongodb://%s:%s", host, port)
	if user != "" && pass != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s", user, pass, host, port)
	}

	// Allow explicit URI override (mis. dari ConfigMap K8s yang menyertakan
	// ?replicaSet=rs0). Sebelumnya MONGO_URI di-set di ConfigMap/Secret/compose
	// tapi tidak pernah dibaca kode.
	if envURI := viper.GetString(fmt.Sprintf("%s_URI", prefix)); envURI != "" {
		uri = envURI
	} else if envURI := viper.GetString("MONGO_URI"); envURI != "" {
		uri = envURI
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)
	clientOpts.SetMaxPoolSize(100)
	clientOpts.SetMinPoolSize(10)
	clientOpts.SetMaxConnIdleTime(30 * time.Minute)
	clientOpts.SetConnectTimeout(10 * time.Second)
	clientOpts.SetServerSelectionTimeout(5 * time.Second)

	client, err := mongo.Connect(clientOpts)
	if err != nil {
		log.Error("Failed to create MongoDB client", zap.Error(err))
		return nil, fmt.Errorf("failed to create MongoDB client: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Error("Failed to ping MongoDB", zap.Error(err))
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Info("MongoDB connection established",
		zap.String("prefix", prefix),
		zap.String("uri", redactMongoURI(uri)),
	)

	return client, nil
}

// redactMongoURI menyembunyikan kredensial di connection string sebelum
// di-log. Tanpa ini, log "connection established" menampilkan host/port
// fallback (MONGO_HOST) alih-alih target yang benar-benar dipakai, sehingga
// salah satu context bisa tampak terhubung ke server yang salah padahal tidak.
func redactMongoURI(uri string) string {
	at := strings.LastIndex(uri, "@")
	if at == -1 {
		return uri
	}
	scheme := strings.Index(uri, "://")
	if scheme == -1 {
		return "mongodb://***"
	}
	return uri[:scheme+3] + "***:***" + uri[at:]
}

// NewMongoDatabase returns the named database from the given client.
//
// Nama database wajib eksplisit — tidak ada fallback ke env lagi. Pemanggil
// bertanggung jawab memastikan nama tidak kosong (pkg/server memvalidasinya
// lewat Config.DBName; migrate/seeder/email memakai konstanta di names.go yang
// tidak mungkin kosong).
func NewMongoDatabase(client *mongo.Client, database string) *mongo.Database {
	return client.Database(database)
}

// GetNextID generates a sequential int32 ID using a counters collection.
// This is used to maintain compatibility with int32 primary keys in protobuf definitions.
func GetNextID(ctx context.Context, db *mongo.Database, collection string) (int32, error) {
	counters := db.Collection("counters")

	filter := bson.M{"_id": collection}
	update := bson.M{
		"$inc": bson.M{"seq": int32(1)},
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result struct {
		Seq int32 `bson:"seq"`
	}

	err := counters.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		return 0, fmt.Errorf("failed to get next ID for %s: %w", collection, err)
	}

	return result.Seq, nil
}
