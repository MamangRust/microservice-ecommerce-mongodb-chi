package apps

import (
	"github.com/MamangRust/microservice-ecommerce-grpc-banner/cache"
	"github.com/MamangRust/microservice-ecommerce-grpc-banner/handler"
	"github.com/MamangRust/microservice-ecommerce-grpc-banner/repository"
	"github.com/MamangRust/microservice-ecommerce-grpc-banner/service"
	pbbanner "github.com/MamangRust/microservice-ecommerce-grpc-pb/banner"
	"github.com/MamangRust/microservice-ecommerce-pkg/server"
	"github.com/MamangRust/microservice-ecommerce-shared/observability"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	db := srv.MongoDB

	repos := repository.NewRepositories(db)
	obs, _ := observability.NewObservability("banner-server", srv.Logger)
	c := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Cache:         c,
		Logger:        srv.Logger,
		Repository:    repos,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{Service: svc, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbbanner.RegisterBannerQueryServiceServer(gs, h.BannerQuery)
		pbbanner.RegisterBannerCommandServiceServer(gs, h.BannerCommand)
	}

	return srv, nil
}
