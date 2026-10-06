package handler

import (
	"github.com/MamangRust/microservice-ecommerce-grpc-merchant_policy/service"
	pbmerchant_policy "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant_policy"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	MerchantPolicyQuery   pbmerchant_policy.MerchantPolicyQueryServiceServer
	MerchantPolicyCommand pbmerchant_policy.MerchantPolicyCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		MerchantPolicyQuery:   NewMerchantPolicyQueryHandler(deps.Service.MerchantPoliciesQuery, deps.Logger),
		MerchantPolicyCommand: NewMerchantPolicyCommandHandler(deps.Service.MerchantPoliciesCommand, deps.Logger),
	}
}
