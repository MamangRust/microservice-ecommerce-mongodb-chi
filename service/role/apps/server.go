package apps

import (
	pbrole "github.com/MamangRust/microservice-ecommerce-grpc-pb/role"
	pbuserrole "github.com/MamangRust/microservice-ecommerce-grpc-pb/user_role"
	"github.com/MamangRust/microservice-ecommerce-grpc-role/cache"
	"github.com/MamangRust/microservice-ecommerce-grpc-role/handler"
	"github.com/MamangRust/microservice-ecommerce-grpc-role/repository"
	"github.com/MamangRust/microservice-ecommerce-grpc-role/service"
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
	obs, _ := observability.NewObservability("role-server", srv.Logger)
	cache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Cache:         cache,
		Logger:        srv.Logger,
		Repository:    repos,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{
		Service: svc,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbrole.RegisterRoleQueryServiceServer(gs, h.RoleQuery)
		pbrole.RegisterRoleCommandServiceServer(gs, h.RoleCommand)
		pbuserrole.RegisterUserRoleServiceServer(gs, h.UserRole)
	}

	return srv, nil
}
