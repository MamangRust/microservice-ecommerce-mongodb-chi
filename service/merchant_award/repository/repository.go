package repository

import (
	pbmerchant "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	merchantadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/merchant"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type GuardOptions struct {
	Merchant []adapter.GuardOption
}

type Repositories struct {
	MerchantAwardQuery   MerchantAwardQueryRepository
	MerchantAwardCommand MerchantAwardCommandRepository
	MerchantQuery        merchantadapter.QueryRepository
}

func NewRepositories(db *mongo.Database, merchantQueryClient pbmerchant.MerchantQueryServiceClient, guards ...GuardOptions) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	return &Repositories{
		MerchantAwardQuery:   NewMerchantAwardQueryRepository(db),
		MerchantAwardCommand: NewMerchantAwardCommandRepository(db),
		MerchantQuery:        merchantadapter.New(merchantQueryClient, g.Merchant...),
	}
}
