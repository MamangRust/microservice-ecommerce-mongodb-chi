package apps

import (
	"github.com/MamangRust/microservice-ecommerce-grpc-category/cache"
	"github.com/MamangRust/microservice-ecommerce-grpc-category/handler"
	"github.com/MamangRust/microservice-ecommerce-grpc-category/repository"
	"github.com/MamangRust/microservice-ecommerce-grpc-category/service"
	pbcategory "github.com/MamangRust/microservice-ecommerce-grpc-pb/category"
	"github.com/MamangRust/microservice-ecommerce-pkg/server"
	"github.com/MamangRust/microservice-ecommerce-shared/observability"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	repos := repository.NewRepositories(srv.MongoDB)
	obs, _ := observability.NewObservability("category-server", srv.Logger)
	c := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Cache:         c,
		Logger:        srv.Logger,
		Repositories:  repos,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{Service: svc, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbcategory.RegisterCategoryQueryServiceServer(gs, h.CategoryQuery)
		pbcategory.RegisterCategoryCommandServiceServer(gs, h.CategoryCommand)
	}

	return srv, nil
}
