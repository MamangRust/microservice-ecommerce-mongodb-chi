package handler

import (
	pbmerchant_policy "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant_policy"
)

type MerchantPolicyQueryHandler interface {
	pbmerchant_policy.MerchantPolicyQueryServiceServer
}

type MerchantPolicyCommandHandler interface {
	pbmerchant_policy.MerchantPolicyCommandServiceServer
}
