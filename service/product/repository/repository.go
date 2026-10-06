package repository

import (
	pbcategory "github.com/MamangRust/microservice-ecommerce-grpc-pb/category"
	pbmerchant "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	categoryadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/category"
	merchantadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/merchant"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type GuardOptions struct {
	Category []adapter.GuardOption
	Merchant []adapter.GuardOption
}

type Repositories struct {
	ProductQuery   ProductQueryRepository
	ProductCommand ProductCommandRepository
	CategoryQuery  categoryadapter.QueryRepository
	MerchantQuery  merchantadapter.QueryRepository
}

func NewRepositories(db *mongo.Database,
	categoryQueryClient pbcategory.CategoryQueryServiceClient,
	merchantQueryClient pbmerchant.MerchantQueryServiceClient,
	guards ...GuardOptions,
) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	categoryQuery := categoryadapter.New(categoryQueryClient, g.Category...)

	return &Repositories{
		ProductQuery:   NewProductQueryRepository(db, categoryQuery),
		ProductCommand: NewProductCommandRepository(db),
		CategoryQuery:  categoryQuery,
		MerchantQuery:  merchantadapter.New(merchantQueryClient, g.Merchant...),
	}
}
