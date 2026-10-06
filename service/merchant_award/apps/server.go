package apps

import (
	"fmt"
	"time"

	pbmerchant "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant"
	pbmerchant_award "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant_award"

	"github.com/MamangRust/microservice-ecommerce-grpc-merchant_award/cache"
	"github.com/MamangRust/microservice-ecommerce-grpc-merchant_award/handler"
	"github.com/MamangRust/microservice-ecommerce-grpc-merchant_award/repository"
	"github.com/MamangRust/microservice-ecommerce-grpc-merchant_award/service"
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

	// MongoDB

	merchantAddr := viper.GetString("GRPC_MERCHANT_ADDR")

	merchantConn, err := grpc.NewClient(
		merchantAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to merchant service: %w", err)
	}

	merchantQueryClient := pbmerchant.NewMerchantQueryServiceClient(merchantConn)

	repos := repository.NewRepositories(srv.MongoDB, merchantQueryClient,
		repository.GuardOptions{
			Merchant: []adapter.GuardOption{
				adapter.WithDependencyGuard(resilience.NewDependencyGuard("merchant", 5, 30, 100, 3*time.Second, srv.Logger)),
			},
		},
	)

	observability, _ := observability.NewObservability("merchant_award-server", srv.Logger)

	cache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Cache:         cache,
		Logger:        srv.Logger,
		Repository:    repos,
		Observability: observability,
	})

	h := handler.NewHandler(&handler.Deps{Service: svc, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbmerchant_award.RegisterMerchantAwardQueryServiceServer(gs, h.MerchantAwardQuery)
		pbmerchant_award.RegisterMerchantAwardCommandServiceServer(gs, h.MerchantAwardCommand)
	}

	return srv, nil
}
