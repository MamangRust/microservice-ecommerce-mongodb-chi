package category_test

import (
	"context"
	cat_handler "github.com/MamangRust/microservice-ecommerce-grpc-category/handler"
	cat_repo "github.com/MamangRust/microservice-ecommerce-grpc-category/repository"
	cat_service "github.com/MamangRust/microservice-ecommerce-grpc-category/service"
	pbcategory "github.com/MamangRust/microservice-ecommerce-grpc-pb/category"
	"testing"

	cat_cache "github.com/MamangRust/microservice-ecommerce-grpc-category/cache"
	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	"github.com/MamangRust/microservice-ecommerce-shared/observability"
	tests "github.com/MamangRust/microservice-ecommerce-test"
	"google.golang.org/protobuf/types/known/emptypb"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

type CategoryGapiTestSuite struct {
	tests.BaseTestSuite
	queryClient   pbcategory.CategoryQueryServiceClient
	commandClient pbcategory.CategoryCommandServiceClient
}

func (s *CategoryGapiTestSuite) SetupSuite() {
	s.PrimaryContext = database.CatalogDB
	s.BaseTestSuite.SetupSuite()

	mongoDB, err := s.MongoDB()
	s.Require().NoError(err)

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)

	// Category dependencies
	mencache := cat_cache.NewMencache(cacheStore)
	repos := cat_repo.NewRepositories(mongoDB)
	svc := cat_service.NewService(&cat_service.Deps{
		Cache:         mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})

	// Handler
	handler := cat_handler.NewHandler(&cat_handler.Deps{
		Service: svc,
		Logger:  s.Log,
	})

	// GRPC Server
	server := grpc.NewServer()
	pbcategory.RegisterCategoryQueryServiceServer(server, handler.CategoryQuery)
	pbcategory.RegisterCategoryCommandServiceServer(server, handler.CategoryCommand)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.queryClient = pbcategory.NewCategoryQueryServiceClient(conn)
	s.commandClient = pbcategory.NewCategoryCommandServiceClient(conn)
}

func (s *CategoryGapiTestSuite) TestCategoryGapiLifecycle() {
	ctx := context.Background()

	// 1. Create
	createRes, err := s.commandClient.Create(ctx, &pbcategory.CreateCategoryRequest{
		Name:          "GAPI Category",
		Description:   "Testing via GRPC",
		SlugCategory:  "gapi-cat",
		ImageCategory: "gapi.jpg",
	})
	s.NoError(err)
	s.NotNil(createRes)
	catID := createRes.Data.Id

	// 2. FindById
	getRes, err := s.queryClient.FindById(ctx, &pbcategory.FindByIdCategoryRequest{Id: catID})
	s.NoError(err)
	s.Equal("GAPI Category", getRes.Data.Name)

	// 3. FindAll
	allRes, err := s.queryClient.FindAll(ctx, &pbcategory.FindAllCategoryRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(allRes.Data)

	// 4. FindByActive
	activeRes, err := s.queryClient.FindByActive(ctx, &pbcategory.FindAllCategoryRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(activeRes.Data)

	// 5. Update
	updateRes, err := s.commandClient.Update(ctx, &pbcategory.UpdateCategoryRequest{
		CategoryId:    catID,
		Name:          "GAPI Category Updated",
		Description:   "Updated via GRPC",
		SlugCategory:  "gapi-cat-updated",
		ImageCategory: "gapi-updated.jpg",
	})
	s.NoError(err)
	s.Equal("GAPI Category Updated", updateRes.Data.Name)

	// 6. Trash
	_, err = s.commandClient.TrashedCategory(ctx, &pbcategory.FindByIdCategoryRequest{Id: catID})
	s.NoError(err)

	// 7. FindByTrashed
	trashedRes, err := s.queryClient.FindByTrashed(ctx, &pbcategory.FindAllCategoryRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 8. Restore
	_, err = s.commandClient.RestoreCategory(ctx, &pbcategory.FindByIdCategoryRequest{Id: catID})
	s.NoError(err)

	// 9. DeletePermanent
	_, _ = s.commandClient.TrashedCategory(ctx, &pbcategory.FindByIdCategoryRequest{Id: catID})
	_, err = s.commandClient.DeleteCategoryPermanent(ctx, &pbcategory.FindByIdCategoryRequest{Id: catID})
	s.NoError(err)

	// 10. RestoreAll
	_, err = s.commandClient.RestoreAllCategory(ctx, &emptypb.Empty{})
	s.NoError(err)

	// 11. DeleteAll
	_, err = s.commandClient.DeleteAllCategoryPermanent(ctx, &emptypb.Empty{})
	s.NoError(err)
}

func TestCategoryGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CategoryGapiTestSuite))
}
