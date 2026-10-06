package handler

import (
	pbmerchant_business "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant_business"
)

type MerchantBusinessQueryHandler interface {
	pbmerchant_business.MerchantBusinessQueryServiceServer
}

type MerchantBusinessCommandHandler interface {
	pbmerchant_business.MerchantBusinessCommandServiceServer
}
