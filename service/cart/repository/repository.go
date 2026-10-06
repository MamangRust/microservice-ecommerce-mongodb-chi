package repository

import (
	pbproduct "github.com/MamangRust/microservice-ecommerce-grpc-pb/product"
	pbuser "github.com/MamangRust/microservice-ecommerce-grpc-pb/user"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	productadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/product"
	useradapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/user"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type GuardOptions struct {
	User    []adapter.GuardOption
	Product []adapter.GuardOption
}

type Repositories struct {
	CartQuery    CartQueryRepository
	CartCommand  CartCommandRepository
	UserQuery    useradapter.QueryRepository
	ProductQuery productadapter.QueryRepository
}

func NewRepositories(DB *mongo.Database,
	userQueryClient pbuser.UserQueryServiceClient,
	productQueryClient pbproduct.ProductQueryServiceClient,
	guards ...GuardOptions,
) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	return &Repositories{
		CartQuery:    NewCartQueryRepository(DB),
		CartCommand:  NewCartCommandRepository(DB),
		UserQuery:    useradapter.New(userQueryClient, nil, g.User...),
		ProductQuery: productadapter.New(productQueryClient, nil, g.Product...),
	}
}
