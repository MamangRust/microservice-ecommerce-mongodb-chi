package handler

import (
	pbmerchant_detail "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant_detail"
	pbmerchant_social_link "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant_social_link"
)

type MerchantDetailQueryHandler interface {
	pbmerchant_detail.MerchantDetailQueryServiceServer
}

type MerchantDetailCommandHandler interface {
	pbmerchant_detail.MerchantDetailCommandServiceServer
}

type MerchantSocialLinkQueryHandler interface {
	pbmerchant_social_link.MerchantSocialCommandServiceServer
}

type MerchantSocialLinkCommandHandler interface {
	pbmerchant_social_link.MerchantSocialCommandServiceServer
}
