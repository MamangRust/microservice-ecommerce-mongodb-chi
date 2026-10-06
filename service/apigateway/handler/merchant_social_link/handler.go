package merchantsociallinkhandler

import (
	"github.com/MamangRust/microservice-ecommerce-grpc-apigateway/apierror"
	pbmerchant_social_link "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant_social_link"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	apimapper "github.com/MamangRust/microservice-ecommerce-shared/mapper/merchant_social_link"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

type DepsMerchantSocialLink struct {
	Client     *grpc.ClientConn
	Router     chi.Router
	Logger     logger.LoggerInterface
	ApiHandler apierror.ApiHandler
}

func RegisterMerchantSocialLinkHandler(deps *DepsMerchantSocialLink) {
	mapper := apimapper.NewMerchantSocialLinkResponseMapper()

	NewMerchantSocialLinkCommandHandleApi(&merchantSocialLinkCommandHandleDeps{
		client:     pbmerchant_social_link.NewMerchantSocialCommandServiceClient(deps.Client),
		router:     deps.Router,
		logger:     deps.Logger,
		mapper:     mapper.CommandMapper(),
		apiHandler: deps.ApiHandler,
	})
}
