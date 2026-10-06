package repository

import (
	pbmerchant "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant"
	pborder "github.com/MamangRust/microservice-ecommerce-grpc-pb/order"
	pborder_item "github.com/MamangRust/microservice-ecommerce-grpc-pb/order_item"
	pbshipping_address "github.com/MamangRust/microservice-ecommerce-grpc-pb/shipping_address"
	pbuser "github.com/MamangRust/microservice-ecommerce-grpc-pb/user"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	merchantadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/merchant"
	orderadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/order"
	orderitemadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/order_item"
	shippingaddressadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/shipping_address"
	useradapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/user"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repositories struct {
	TransactionCommand TransactionCommandRepository
	TransactionQuery   TransactionQueryRepository
	OrderItem          orderitemadapter.QueryRepository
	OrderQuery         orderadapter.QueryRepository
	MerchantQuery      merchantadapter.QueryRepository
	ShippingAddress    shippingaddressadapter.QueryRepository
	UserQuery          useradapter.QueryRepository
	Outbox             OutboxRepository
}

// GuardOptions carries the resilience guard options for each outbound
// dependency.
type GuardOptions struct {
	User      []adapter.GuardOption
	Merchant  []adapter.GuardOption
	Order     []adapter.GuardOption
	OrderItem []adapter.GuardOption
	Shipping  []adapter.GuardOption
}

type Deps struct {
	DB *mongo.Database

	UserQueryClient      pbuser.UserQueryServiceClient
	MerchantQueryClient  pbmerchant.MerchantQueryServiceClient
	OrderQueryClient     pborder.OrderQueryServiceClient
	OrderItemQueryClient pborder_item.OrderItemQueryServiceClient
	ShippingQueryClient  pbshipping_address.ShippingQueryServiceClient

	Guards GuardOptions
}

func NewRepositories(deps *Deps) *Repositories {
	g := deps.Guards

	return &Repositories{
		TransactionCommand: NewTransactionCommandRepository(deps.DB),
		TransactionQuery:   NewTransactionQueryRepository(deps.DB),
		OrderItem:          orderitemadapter.New(deps.OrderItemQueryClient, nil, g.OrderItem...),
		OrderQuery:         orderadapter.New(deps.OrderQueryClient, g.Order...),
		MerchantQuery:      merchantadapter.New(deps.MerchantQueryClient, g.Merchant...),
		ShippingAddress:    shippingaddressadapter.New(deps.ShippingQueryClient, nil, g.Shipping...),
		UserQuery:          useradapter.New(deps.UserQueryClient, nil, g.User...),
		Outbox:             NewOutboxRepository(deps.DB),
	}
}
