package handler

import (
	pbmerchant_award "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant_award"
)

type MerchantAwardQueryHandler interface {
	pbmerchant_award.MerchantAwardQueryServiceServer
}

type MerchantAwardCommandHandler interface {
	pbmerchant_award.MerchantAwardCommandServiceServer
}
