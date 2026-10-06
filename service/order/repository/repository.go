package repository

import (
	pbmerchant "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant"
	pborder_item "github.com/MamangRust/microservice-ecommerce-grpc-pb/order_item"
	pbproduct "github.com/MamangRust/microservice-ecommerce-grpc-pb/product"
	pbshipping_address "github.com/MamangRust/microservice-ecommerce-grpc-pb/shipping_address"
	pbtransaction "github.com/MamangRust/microservice-ecommerce-grpc-pb/transaction"
	pbuser "github.com/MamangRust/microservice-ecommerce-grpc-pb/user"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	merchantadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/merchant"
	orderitemadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/order_item"
	productadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/product"
	shippingaddressadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/shipping_address"
	transactionadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/transaction"
	useradapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/user"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type GuardOptions struct {
	User        []adapter.GuardOption
	Product     []adapter.GuardOption
	Merchant    []adapter.GuardOption
	OrderItem   []adapter.GuardOption
	Shipping    []adapter.GuardOption
	Transaction []adapter.GuardOption
}

type Repositories struct {
	OrderQuery         OrderQueryRepository
	OrderCommand       OrderCommandRepository
	UserQuery          useradapter.QueryRepository
	ProductQuery       productadapter.QueryRepository
	ProductCommand     productadapter.CommandRepository
	OrderItemQuery     orderitemadapter.QueryRepository
	OrderItemCommand   orderitemadapter.CommandRepository
	MerchantQuery      merchantadapter.QueryRepository
	ShippingAddress    shippingaddressadapter.CommandRepository
	TransactionCommand transactionadapter.CommandRepository
	ShippingQuery      shippingaddressadapter.QueryRepository
	StockReservation   StockReservationRepository
	Outbox             OutboxRepository
}

type Deps struct {
	DB        *mongo.Database
	CatalogDB *mongo.Database

	UserQueryClient          pbuser.UserQueryServiceClient
	ProductQueryClient       pbproduct.ProductQueryServiceClient
	ProductCommandClient     pbproduct.ProductCommandServiceClient
	MerchantQueryClient      pbmerchant.MerchantQueryServiceClient
	OrderItemQueryClient     pborder_item.OrderItemQueryServiceClient
	OrderItemCommandClient   pborder_item.OrderItemCommandServiceClient
	ShippingCommandClient    pbshipping_address.ShippingCommandServiceClient
	ShippingQueryClient      pbshipping_address.ShippingQueryServiceClient
	TransactionCommandClient pbtransaction.TransactionCommandServiceClient

	Guards GuardOptions
}

func NewRepositories(deps *Deps) *Repositories {
	g := deps.Guards

	productAdapter := productadapter.New(deps.ProductQueryClient, deps.ProductCommandClient, g.Product...)
	orderItemAdapter := orderitemadapter.New(deps.OrderItemQueryClient, deps.OrderItemCommandClient, g.OrderItem...)
	shippingAdapter := shippingaddressadapter.New(deps.ShippingQueryClient, deps.ShippingCommandClient, g.Shipping...)

	return &Repositories{
		OrderQuery:         NewOrderQueryRepository(deps.DB),
		OrderCommand:       NewOrderCommandRepository(deps.DB),
		UserQuery:          useradapter.New(deps.UserQueryClient, nil, g.User...),
		ProductQuery:       productAdapter,
		ProductCommand:     productAdapter,
		OrderItemQuery:     orderItemAdapter,
		OrderItemCommand:   orderItemAdapter,
		MerchantQuery:      merchantadapter.New(deps.MerchantQueryClient, g.Merchant...),
		ShippingAddress:    shippingAdapter,
		TransactionCommand: transactionadapter.NewCommandAdapter(deps.TransactionCommandClient, g.Transaction...),
		ShippingQuery:      shippingAdapter,
		StockReservation:   NewStockReservationRepository(deps.DB, deps.CatalogDB),
		Outbox:             NewOutboxRepository(deps.DB),
	}
}
