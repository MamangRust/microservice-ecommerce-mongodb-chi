package cache

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-transaction/repository"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type TransactionQueryCache interface {
	GetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransaction) ([]*repository.TransactionResult, *int, bool)
	SetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransaction, data []*repository.TransactionResult, total *int)

	GetCachedTransactionByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant) ([]*repository.TransactionResult, *int, bool)
	SetCachedTransactionByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant, data []*repository.TransactionResult, total *int)

	GetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransaction) ([]*repository.TransactionResult, *int, bool)
	SetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransaction, data []*repository.TransactionResult, total *int)

	GetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransaction) ([]*repository.TransactionResult, *int, bool)
	SetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransaction, data []*repository.TransactionResult, total *int)

	GetCachedTransactionCache(ctx context.Context, id int) (*repository.TransactionResult, bool)
	SetCachedTransactionCache(ctx context.Context, data *repository.TransactionResult)

	GetCachedTransactionByOrderId(ctx context.Context, orderID int) (*repository.TransactionResult, bool)
	SetCachedTransactionByOrderId(ctx context.Context, orderID int, data *repository.TransactionResult)
}

type TransactionCommandCache interface {
	DeleteTransactionCache(ctx context.Context, transactionID int)
	InvalidateTransactionCache(ctx context.Context)
}
