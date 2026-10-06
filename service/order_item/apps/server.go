package apps

import (
	"github.com/MamangRust/microservice-ecommerce-grpc-order-item/cache"
	"github.com/MamangRust/microservice-ecommerce-grpc-order-item/handler"
	"github.com/MamangRust/microservice-ecommerce-grpc-order-item/repository"
	"github.com/MamangRust/microservice-ecommerce-grpc-order-item/service"
	pborder_item "github.com/MamangRust/microservice-ecommerce-grpc-pb/order_item"
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
	obs, _ := observability.NewObservability("order_item-server", srv.Logger)
	cache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Cache:         cache,
		Logger:        srv.Logger,
		Repository:    repos,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{Service: svc, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pborder_item.RegisterOrderItemQueryServiceServer(gs, h.OrderItemQuery)
		pborder_item.RegisterOrderItemCommandServiceServer(gs, h.OrderItemCommand)
	}

	return srv, nil
}
