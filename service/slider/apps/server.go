package apps

import (
	pbslider "github.com/MamangRust/microservice-ecommerce-grpc-pb/slider"
	"github.com/MamangRust/microservice-ecommerce-grpc-slider/cache"
	"github.com/MamangRust/microservice-ecommerce-grpc-slider/handler"
	"github.com/MamangRust/microservice-ecommerce-grpc-slider/repository"
	"github.com/MamangRust/microservice-ecommerce-grpc-slider/service"
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
	c := cache.NewMencache(srv.CacheStore)
	obs, _ := observability.NewObservability("slider-server", srv.Logger)

	svc := service.NewService(&service.Deps{
		Repositories:  repos,
		Mencache:      c,
		Logger:        srv.Logger,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{
		Service: svc,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbslider.RegisterSliderQueryServiceServer(gs, h.SliderQuery)
		pbslider.RegisterSliderCommandServiceServer(gs, h.SliderCommand)
	}

	return srv, nil
}
