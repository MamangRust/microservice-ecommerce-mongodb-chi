package cart_test

import (
	"context"
	"testing"

	pbproduct "github.com/MamangRust/microservice-ecommerce-grpc-pb/product"
	pbuser "github.com/MamangRust/microservice-ecommerce-grpc-pb/user"

	"github.com/stretchr/testify/suite"

	cart_cache "github.com/MamangRust/microservice-ecommerce-grpc-cart/cache"
	"github.com/MamangRust/microservice-ecommerce-grpc-cart/repository"
	"github.com/MamangRust/microservice-ecommerce-grpc-cart/service"
	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/observability"
	tests "github.com/MamangRust/microservice-ecommerce-test"
)

type CartServiceTestSuite struct {
	tests.BaseTestSuite
	svc *service.Service
}

func (s *CartServiceTestSuite) SetupSuite() {
	s.PrimaryContext = database.ExperienceDB
	s.BaseTestSuite.SetupSuite()

	mongoDB, err := s.MongoDB()
	s.Require().NoError(err)

	// Setup dependencies
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)

	// Cart dependencies
	mencache := cart_cache.NewMencache(cacheStore)
	repos := repository.NewRepositories(
		mongoDB,
		pbuser.NewUserQueryServiceClient(s.Conns["user"]),
		pbproduct.NewProductQueryServiceClient(s.Conns["product"]),
	)

	s.svc = service.NewService(&service.Deps{
		Cache:         mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
}

func (s *CartServiceTestSuite) TestCartServiceLifecycle() {
	ctx := context.Background()

	// 1. Seed dependencies
	userID := s.SeedUser(ctx)
	categoryID := s.SeedCategory(ctx)
	merchantID := s.SeedMerchant(ctx, userID)
	prodID := s.SeedProduct(ctx, merchantID, categoryID)

	// 2. Add to Cart via Service
	req := &requests.CreateCartRequest{
		UserID:    userID,
		ProductID: prodID,
		Quantity:  3,
	}

	created, err := s.svc.CartCommand.Create(ctx, req)
	s.NoError(err)
	s.NotNil(created)

	// 3. Get Cart Items via Service
	items, _, err := s.svc.CartQuery.FindAll(ctx, &requests.FindAllCarts{
		UserID:   userID,
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotEmpty(items)
}

func TestCartServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CartServiceTestSuite))
}
