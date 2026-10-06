package apps

import (
	"fmt"
	"time"

	pbcategory "github.com/MamangRust/microservice-ecommerce-grpc-pb/category"
	pbmerchant "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant"
	pbproduct "github.com/MamangRust/microservice-ecommerce-grpc-pb/product"

	"github.com/MamangRust/microservice-ecommerce-grpc-product/cache"
	"github.com/MamangRust/microservice-ecommerce-grpc-product/handler"
	"github.com/MamangRust/microservice-ecommerce-grpc-product/repository"
	"github.com/MamangRust/microservice-ecommerce-grpc-product/service"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
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

	categoryAddr := viper.GetString("GRPC_CATEGORY_ADDR")

	categoryConn, err := grpc.NewClient(categoryAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to category service: %w", err)
	}
	categoryQueryClient := pbcategory.NewCategoryQueryServiceClient(categoryConn)

	merchantAddr := viper.GetString("GRPC_MERCHANT_ADDR")

	merchantConn, err := grpc.NewClient(merchantAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to merchant service: %w", err)
	}
	merchantQueryClient := pbmerchant.NewMerchantQueryServiceClient(merchantConn)

	repos := repository.NewRepositories(srv.MongoDB,
		categoryQueryClient,
		merchantQueryClient,
		repository.GuardOptions{
			Category: []adapter.GuardOption{
				adapter.WithDependencyGuard(resilience.NewDependencyGuard("category", 5, 30, 100, 3*time.Second, srv.Logger)),
			},
			Merchant: []adapter.GuardOption{
				adapter.WithDependencyGuard(resilience.NewDependencyGuard("merchant", 5, 30, 100, 3*time.Second, srv.Logger)),
			},
		},
	)

	obs, _ := observability.NewObservability("product-server", srv.Logger)
	c := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Cache:         c,
		Logger:        srv.Logger,
		Repository:    repos,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{Service: svc, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbproduct.RegisterProductQueryServiceServer(gs, h.ProductQuery)
		pbproduct.RegisterProductCommandServiceServer(gs, h.ProductCommand)
	}

	return srv, nil
}
