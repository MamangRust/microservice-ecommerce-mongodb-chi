package handler

import (
	pbtransaction "github.com/MamangRust/microservice-ecommerce-grpc-pb/transaction"
)

type TransactionQueryHandler interface {
	pbtransaction.TransactionQueryServiceServer
}

type TransactionCommandHandler interface {
	pbtransaction.TransactionCommandServiceServer
}
