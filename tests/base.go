package tests

import (
	"context"
	"fmt"
	"reflect"
	"time"

	pbcategory "github.com/MamangRust/microservice-ecommerce-grpc-pb/category"
	pbmerchant "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant"
	pborder "github.com/MamangRust/microservice-ecommerce-grpc-pb/order"
	pborder_item "github.com/MamangRust/microservice-ecommerce-grpc-pb/order_item"
	pbproduct "github.com/MamangRust/microservice-ecommerce-grpc-pb/product"
	pbreview "github.com/MamangRust/microservice-ecommerce-grpc-pb/review"
	pbshipping_address "github.com/MamangRust/microservice-ecommerce-grpc-pb/shipping_address"
	pbuser "github.com/MamangRust/microservice-ecommerce-grpc-pb/user"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/observability"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/v2/mongo"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type BaseTestSuite struct {
	suite.Suite
	ts             *TestSuite
	PrimaryContext string
	Log            logger.LoggerInterface
	Obs            observability.TraceLoggerObservability
	Conns          map[string]*grpc.ClientConn
	Servers        []*grpc.Server
	Ctx            context.Context
	Cancel         context.CancelFunc

	// statsContainer is the ClickHouse testcontainer started by SetupStatsEnv
	// and terminated in TearDownSuite.
	statsContainer testcontainers.Container
}

// seedUserCounter and seedUniqueSuffix make seeded values unique across multiple
// seed calls within the same test process (suites share one database per package
// run), so fixed slugs/emails never collide between test methods.
var seedUserCounter int64

func uniqueSuffix() string {
	seedUserCounter++
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), seedUserCounter)
}

func (s *BaseTestSuite) SetupSuite() {
	s.Ctx, s.Cancel = context.WithCancel(context.Background())

	ts, err := SetupTestSuite(s.PrimaryContext)
	s.Require().NoError(err)
	s.ts = ts

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	s.Log, _ = logger.NewLogger("test", lp)

	if s.Log == nil || (reflect.ValueOf(s.Log).Kind() == reflect.Ptr && reflect.ValueOf(s.Log).IsNil()) {
		z, _ := zap.NewDevelopment()
		s.Log = &logger.Logger{Log: z}
	}

	s.Obs, err = observability.NewObservability("test", s.Log)
	s.Require().NoError(err)
	s.Require().NotNil(s.Obs)
	s.Conns = make(map[string]*grpc.ClientConn)
}

func (s *BaseTestSuite) TearDownSuite() {
	if s.statsContainer != nil {
		_ = s.statsContainer.Terminate(context.Background())
	}
	for _, conn := range s.Conns {
		conn.Close()
	}
	for _, server := range s.Servers {
		server.GracefulStop()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
	if s.Cancel != nil {
		s.Cancel()
	}
}

func (s *BaseTestSuite) RedisClient() *goredis.Client {
	return s.ts.RedisClient()
}

// MongoDB returns the primary-context database backing the integration tests.
func (s *BaseTestSuite) MongoDB() (*mongo.Database, error) {
	return s.ts.MongoDB(s.Ctx)
}

// MongoURI is the testcontainers connection string, for components that open
// their own Mongo connection instead of accepting a *mongo.Database (e.g. the
// stats-writer backfiller).
func (s *BaseTestSuite) MongoURI() string {
	return s.ts.MongoURI
}

// DB returns the primary-context database (convenience, no error).
func (s *BaseTestSuite) DB() *mongo.Database {
	return s.ts.DB(s.Ctx)
}

// DBForCollection returns the database that owns the given collection.
func (s *BaseTestSuite) DBForCollection(coll string) *mongo.Database {
	return s.ts.DBForCollection(coll)
}

// Coll returns a handle to the given collection in whichever database owns it.
func (s *BaseTestSuite) Coll(coll string) *mongo.Collection {
	return s.ts.Coll(coll)
}

// DBForService returns the database owned by the given service.
func (s *BaseTestSuite) DBForService(svc string) *mongo.Database {
	return s.ts.DBForService(svc)
}

// DBForContext returns the database for an explicit bounded context. Usable from
// any subpackage test — delegates to *TestSuite.DBForContext.
func (s *BaseTestSuite) DBForContext(ctx context.Context, context string) *mongo.Database {
	return s.ts.DBForContext(ctx, context)
}

func (s *BaseTestSuite) RegisterServer(server *grpc.Server) string {
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Servers = append(s.Servers, server)
	return addr
}

func (s *BaseTestSuite) GetConnection(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	s.Require().NoError(err)
	return conn
}

func (s *BaseTestSuite) SeedUser(ctx context.Context) int {
	// Each call seeds a unique email so suites with multiple test methods can
	// share one database without colliding on a fixed address.
	email := fmt.Sprintf("seed.user.%s@example.com", uniqueSuffix())
	res, err := pbuser.NewUserCommandServiceClient(s.Conns["user"]).Create(ctx, &pbuser.CreateUserRequest{
		Firstname:       "Seed",
		Lastname:        "User",
		Email:           email,
		Password:        "password123",
		ConfirmPassword: "password123",
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}

func (s *BaseTestSuite) SeedCategory(ctx context.Context) int {
	seedSuffix := uniqueSuffix()
	res, err := pbcategory.NewCategoryCommandServiceClient(s.Conns["category"]).Create(ctx, &pbcategory.CreateCategoryRequest{
		Name:          "Seed Category " + seedSuffix,
		Description:   "Seed Description",
		SlugCategory:  "seed-category-" + seedSuffix,
		ImageCategory: "seed.jpg",
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}

func (s *BaseTestSuite) SeedMerchant(ctx context.Context, userID int) int {
	res, err := pbmerchant.NewMerchantCommandServiceClient(s.Conns["merchant"]).Create(ctx, &pbmerchant.CreateMerchantRequest{
		UserId:       int32(userID),
		Name:         "Seed Merchant",
		Description:  "Seed Description",
		Address:      "Seed Address",
		ContactEmail: "merchant@example.com",
		ContactPhone: "08123456789",
		Status:       "active",
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}

func (s *BaseTestSuite) SeedProduct(ctx context.Context, merchantID int, categoryID int) int {
	seedSuffix := uniqueSuffix()
	res, err := pbproduct.NewProductCommandServiceClient(s.Conns["product"]).Create(ctx, &pbproduct.CreateProductRequest{
		MerchantId:   int32(merchantID),
		CategoryId:   int32(categoryID),
		Name:         "Seed Product " + seedSuffix,
		Description:  "Seed Description",
		Price:        10000,
		CountInStock: 100,
		Brand:        "Seed Brand",
		Weight:       1000,
		SlugProduct:  "seed-product-" + seedSuffix,
		ImageProduct: "seed.jpg",
		Barcode:      "123456789",
		Rating:       5,
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}

func (s *BaseTestSuite) SeedShippingAddress(ctx context.Context, orderID int) int {
	res, err := pbshipping_address.NewShippingCommandServiceClient(s.Conns["shipping-address"]).CreateShipping(ctx, &pbshipping_address.CreateShippingAddressRequest{
		OrderId:        int32(orderID),
		Alamat:         "Seed Address",
		Provinsi:       "Seed Province",
		Kota:           "Seed City",
		Negara:         "Seed Country",
		Courier:        "Seed Courier",
		ShippingMethod: "Seed Method",
		ShippingCost:   1000,
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}

func (s *BaseTestSuite) SeedOrder(ctx context.Context, userID int, merchID int, prodID int) int {
	res, err := pborder.NewOrderCommandServiceClient(s.Conns["order"]).Create(ctx, &pborder.CreateOrderRequest{
		UserId:     int32(userID),
		MerchantId: int32(merchID),
		TotalPrice: 10000,
		Items: []*pborder.CreateOrderItemRequest{
			{
				ProductId: int32(prodID),
				Quantity:  1,
				Price:     10000,
			},
		},
		Shipping: &pbshipping_address.CreateShippingAddressRequest{
			Alamat:         "Seed Address",
			Provinsi:       "Seed Province",
			Kota:           "Seed City",
			Negara:         "Seed Country",
			Courier:        "Seed Courier",
			ShippingMethod: "Seed Method",
			ShippingCost:   1000,
		},
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}

func (s *BaseTestSuite) SeedReview(ctx context.Context, userID int, productID int) int {
	res, err := pbreview.NewReviewCommandServiceClient(s.Conns["review"]).Create(ctx, &pbreview.CreateReviewRequest{
		UserId:    int32(userID),
		ProductId: int32(productID),
		Rating:    5,
		Comment:   "Seed Review",
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}

func (s *BaseTestSuite) SeedOrderItem(ctx context.Context, orderID int, productID int) int {
	res, err := pborder_item.NewOrderItemCommandServiceClient(s.Conns["order-item"]).CreateOrderItem(ctx, &pborder_item.CreateOrderItemRecordRequest{
		OrderId:   int32(orderID),
		ProductId: int32(productID),
		Quantity:  1,
		Price:     1000,
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}
