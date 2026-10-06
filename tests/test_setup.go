package tests

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"

	goredis "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	mongoopts "go.mongodb.org/mongo-driver/v2/mongo/options"
	"google.golang.org/grpc"
)

// contextDatabases memetakan nama context (konstanta di pkg/database/names.go)
// ke nama database fisik. TestSuite membuka satu client lalu menyiapkan keenam
// database ini, sehingga service yang menulis ke context berbeda tetap ter-cover
// dalam satu testcontainer.
var contextDatabases = []string{
	database.IdentityDB,
	database.CatalogDB,
	database.MerchantDB,
	database.SalesDB,
	database.ExperienceDB,
	database.EmailDB,
}

// collectionContext memetakan nama koleksi ke bounded context pemiliknya. Dipakai
// oleh DBForCollection/Coll agar test bisa mengakses koleksi lintas context tanpa
// harus tahu database mana yang menampungnya.
var collectionContext = map[string]string{
	"users":                              database.IdentityDB,
	"roles":                              database.IdentityDB,
	"user_roles":                         database.IdentityDB,
	"refresh_tokens":                     database.IdentityDB,
	"reset_tokens":                       database.IdentityDB,
	"categories":                         database.CatalogDB,
	"products":                           database.CatalogDB,
	"product_stock_adjustments":          database.CatalogDB,
	"merchants":                          database.MerchantDB,
	"merchant_documents":                 database.MerchantDB,
	"merchant_details":                   database.MerchantDB,
	"merchant_social_media_links":        database.MerchantDB,
	"merchant_business_information":      database.MerchantDB,
	"merchant_policies":                  database.MerchantDB,
	"merchant_certifications_and_awards": database.MerchantDB,
	"orders":                             database.SalesDB,
	"order_items":                        database.SalesDB,
	"order_stock_reservations":           database.SalesDB,
	"transactions":                       database.SalesDB,
	"banners":                            database.ExperienceDB,
	"sliders":                            database.ExperienceDB,
	"shipping_addresses":                 database.ExperienceDB,
	"carts":                              database.ExperienceDB,
	"reviews":                            database.ExperienceDB,
	"review_details":                     database.ExperienceDB,
	"outbox_events":                      database.EmailDB,
	"consumer_inbox":                     database.EmailDB,
}

// serviceContext memetakan nama service ke bounded context pemilik database-nya.
var serviceContext = map[string]string{
	"auth":              database.IdentityDB,
	"user":              database.IdentityDB,
	"role":              database.IdentityDB,
	"category":          database.CatalogDB,
	"product":           database.CatalogDB,
	"merchant":          database.MerchantDB,
	"merchant_award":    database.MerchantDB,
	"merchant_business": database.MerchantDB,
	"merchant_detail":   database.MerchantDB,
	"merchant_policy":   database.MerchantDB,
	"order":             database.SalesDB,
	"order_item":        database.SalesDB,
	"transaction":       database.SalesDB,
	"cart":              database.ExperienceDB,
	"shipping_address":  database.ExperienceDB,
	"banner":            database.ExperienceDB,
	"slider":            database.ExperienceDB,
	"review":            database.ExperienceDB,
	"review_detail":     database.ExperienceDB,
	"email":             database.EmailDB,
}

type TestSuite struct {
	MongoContainer *mongodb.MongoDBContainer
	RedisContainer *redis.RedisContainer
	MongoClient    *mongo.Client
	Databases      map[string]*mongo.Database
	PrimaryContext string
	RedisURL       string
	MongoURI       string
	Ctx            context.Context
}

// SetupTestSuite membuka satu MongoDB testcontainer, menyiapkan keenam database
// bounded context, lalu membuat index di masing-masing (EnsureIndexes).
// primaryContext menentukan database yang dikembalikan oleh DB()/MongoDB() — yakni
// database milik service yang sedang diuji paket test tersebut.
func SetupTestSuite(primaryContext string) (*TestSuite, error) {
	ctx := context.Background()

	// Setup MongoDB
	mongoContainer, err := mongodb.Run(ctx,
		"mongo:7.0",
		testcontainers.WithWaitStrategy(
			wait.ForLog("Waiting for connections").WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start mongo container: %w", err)
	}

	mongoURI, err := mongoContainer.ConnectionString(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get mongo connection string: %w", err)
	}

	// Setup Redis
	redisContainer, err := redis.Run(ctx, "redis:7-alpine")
	if err != nil {
		return nil, fmt.Errorf("failed to start redis container: %w", err)
	}

	redisURL, err := redisContainer.ConnectionString(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get redis connection string: %w", err)
	}

	// Connect to MongoDB
	client, err := mongo.Connect(mongoopts.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping mongo: %w", err)
	}

	databases := make(map[string]*mongo.Database, len(contextDatabases))
	for _, name := range contextDatabases {
		databases[name] = client.Database(name)
	}

	ts := &TestSuite{
		MongoContainer: mongoContainer,
		RedisContainer: redisContainer,
		MongoClient:    client,
		Databases:      databases,
		PrimaryContext: primaryContext,
		RedisURL:       redisURL,
		MongoURI:       mongoURI,
		Ctx:            ctx,
	}

	// Buat index di seluruh database bounded context (mirip service/migrate).
	// Harness memakai satu container untuk keenam context, jadi memakai varian
	// single-client, bukan EnsureIndexes yang membuka enam koneksi.
	lg, _ := logger.NewLogger("test", nil)
	if err := database.EnsureIndexesOnClient(lg, client); err != nil {
		ts.Teardown()
		return nil, fmt.Errorf("failed to ensure indexes: %w", err)
	}

	// Seed minimal roles yang dipakai suite auth/user.
	if err := ts.SeedMinimalRoles(); err != nil {
		ts.Teardown()
		return nil, fmt.Errorf("failed to seed minimal roles: %w", err)
	}

	return ts, nil
}

// DB returns the primary-context database backing the integration tests.
func (ts *TestSuite) DB(ctx context.Context) *mongo.Database {
	return ts.Databases[ts.PrimaryContext]
}

// DBForContext returns the database for an explicit bounded context.
func (ts *TestSuite) DBForContext(ctx context.Context, context string) *mongo.Database {
	return ts.Databases[context]
}

// DBForCollection returns the database that owns the given collection.
func (ts *TestSuite) DBForCollection(coll string) *mongo.Database {
	if ctx, ok := collectionContext[coll]; ok {
		return ts.Databases[ctx]
	}
	return ts.Databases[ts.PrimaryContext]
}

// DBForService returns the database owned by the given service.
func (ts *TestSuite) DBForService(svc string) *mongo.Database {
	if ctx, ok := serviceContext[svc]; ok {
		return ts.Databases[ctx]
	}
	return ts.Databases[ts.PrimaryContext]
}

// Coll returns a handle to the given collection in whichever database owns it.
func (ts *TestSuite) Coll(coll string) *mongo.Collection {
	return ts.DBForCollection(coll).Collection(coll)
}

// MongoDB retains its historical name for compatibility with existing suites but
// now returns the primary-context MongoDB database backing the integration tests.
func (ts *TestSuite) MongoDB(ctx context.Context) (*mongo.Database, error) {
	return ts.DB(ctx), nil
}

// SeedMinimalRoles inserts default roles for integration tests (into ec_identity).
func (ts *TestSuite) SeedMinimalRoles() error {
	coll := ts.Databases[database.IdentityDB].Collection("roles")
	now := time.Now().UTC()

	for _, name := range []string{"ROLE_ADMIN", "ROLE_USER"} {
		count, _ := coll.CountDocuments(ts.Ctx, bson.M{"role_name": name})
		if count > 0 {
			continue
		}

		roleID, err := database.GetNextID(ts.Ctx, ts.Databases[database.IdentityDB], "roles")
		if err != nil {
			return err
		}

		_, err = coll.InsertOne(ts.Ctx, bson.M{
			"_id":        roleID,
			"role_name":  name,
			"created_at": now,
			"updated_at": now,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (ts *TestSuite) RedisClient() *goredis.Client {
	opts, _ := goredis.ParseURL(ts.RedisURL)
	return goredis.NewClient(opts)
}

func (ts *TestSuite) Teardown() {
	if ts.MongoClient != nil {
		_ = ts.MongoClient.Disconnect(context.Background())
	}
	if ts.MongoContainer != nil {
		if err := ts.MongoContainer.Terminate(ts.Ctx); err != nil {
			log.Printf("failed to terminate mongo container: %v", err)
		}
	}
	if ts.RedisContainer != nil {
		if err := ts.RedisContainer.Terminate(ts.Ctx); err != nil {
			log.Printf("failed to terminate redis container: %v", err)
		}
	}
}

type GRPCServerRunner interface {
	Serve(lis net.Listener) error
	Stop()
	GracefulStop()
}

func RunGRPCServer(server *grpc.Server) (string, error) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return "", err
	}
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Printf("grpc server error: %v", err)
		}
	}()
	return lis.Addr().String(), nil
}
