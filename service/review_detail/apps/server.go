package apps

import (
	pbreview_detail "github.com/MamangRust/microservice-ecommerce-grpc-pb/review_detail"
	"github.com/MamangRust/microservice-ecommerce-grpc-review-detail/cache"
	"github.com/MamangRust/microservice-ecommerce-grpc-review-detail/handler"
	"github.com/MamangRust/microservice-ecommerce-grpc-review-detail/repository"
	"github.com/MamangRust/microservice-ecommerce-grpc-review-detail/service"
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

	obs, _ := observability.NewObservability("review_detail-server", srv.Logger)
	c := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Observability: obs,
		Cache:         c,
		Repositories:  repos,
		Logger:        srv.Logger,
	})

	h := handler.NewHandler(&handler.Deps{Service: svc, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbreview_detail.RegisterReviewDetailQueryServiceServer(gs, h.ReviewDetailQuery)
		pbreview_detail.RegisterReviewDetailCommandServiceServer(gs, h.ReviewDetailCommand)
	}

	return srv, nil
}
