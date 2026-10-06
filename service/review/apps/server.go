package apps

import (
	"fmt"
	"time"

	pbproduct "github.com/MamangRust/microservice-ecommerce-grpc-pb/product"
	pbreview "github.com/MamangRust/microservice-ecommerce-grpc-pb/review"
	pbuser "github.com/MamangRust/microservice-ecommerce-grpc-pb/user"

	"github.com/MamangRust/microservice-ecommerce-grpc-review/cache"
	"github.com/MamangRust/microservice-ecommerce-grpc-review/handler"
	"github.com/MamangRust/microservice-ecommerce-grpc-review/repository"
	"github.com/MamangRust/microservice-ecommerce-grpc-review/service"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/resilience"
	"github.com/MamangRust/microservice-ecommerce-pkg/server"
	"github.com/MamangRust/microservice-ecommerce-shared/observability"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	userAddr := viper.GetString("GRPC_USER_ADDR")
	userConn, err := grpc.NewClient(userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}
	userQueryClient := pbuser.NewUserQueryServiceClient(userConn)

	productAddr := viper.GetString("GRPC_PRODUCT_ADDR")
	productConn, err := grpc.NewClient(productAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}
	productQueryClient := pbproduct.NewProductQueryServiceClient(productConn)

	// ec_catalog punya server MongoDB sendiri, jadi ini membuka koneksi kedua.
	catalogDB, err := srv.ContextDatabase(database.CatalogDB)
	if err != nil {
		return nil, err
	}

	repos := repository.NewRepositories(srv.MongoDB,
		catalogDB,
		userQueryClient,
		productQueryClient,
		repository.GuardOptions{
			User: []adapter.GuardOption{
				adapter.WithDependencyGuard(resilience.NewDependencyGuard("user", 5, 30, 100, 3*time.Second, srv.Logger)),
			},
			Product: []adapter.GuardOption{
				adapter.WithDependencyGuard(resilience.NewDependencyGuard("product", 5, 30, 100, 3*time.Second, srv.Logger)),
			},
		},
	)

	obs, _ := observability.NewObservability("review-server", srv.Logger)
	c := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Observability: obs,
		Cache:         c,
		Repositories:  repos,
		Logger:        srv.Logger,
	})

	h := handler.NewHandler(&handler.Deps{Service: svc, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbreview.RegisterReviewQueryServiceServer(gs, h.ReviewQuery)
		pbreview.RegisterReviewCommandServiceServer(gs, h.ReviewCommand)
	}

	return srv, nil
}
