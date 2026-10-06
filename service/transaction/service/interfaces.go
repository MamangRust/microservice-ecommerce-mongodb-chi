package service

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-transaction/repository"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type TransactionQueryService interface {
	FindAll(
		ctx context.Context,
		req *requests.FindAllTransaction,
	) ([]*repository.TransactionResult, *int, error)

	FindActive(
		ctx context.Context,
		req *requests.FindAllTransaction,
	) ([]*repository.TransactionResult, *int, error)

	FindTrashed(
		ctx context.Context,
		req *requests.FindAllTransaction,
	) ([]*repository.TransactionResult, *int, error)

	FindByMerchant(
		ctx context.Context,
		req *requests.FindAllTransactionByMerchant,
	) ([]*repository.TransactionResult, *int, error)

	FindByID(
		ctx context.Context,
		transaction_id int,
	) (*repository.TransactionResult, error)

	FindByOrderID(
		ctx context.Context,
		order_id int,
	) (*repository.TransactionResult, error)
}

type TransactionCommandService interface {
	Create(
		ctx context.Context,
		request *requests.CreateTransactionRequest,
	) (*repository.TransactionResult, error)

	Update(
		ctx context.Context,
		request *requests.UpdateTransactionRequest,
	) (*repository.TransactionResult, error)

	Trash(
		ctx context.Context,
		transaction_id int,
	) (*repository.TransactionResult, error)

	Restore(
		ctx context.Context,
		transaction_id int,
	) (*repository.TransactionResult, error)

	DeletePermanent(
		ctx context.Context,
		transaction_id int,
	) (bool, error)

	DeleteByOrderIDPermanent(
		ctx context.Context,
		order_id int,
	) (bool, error)

	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
