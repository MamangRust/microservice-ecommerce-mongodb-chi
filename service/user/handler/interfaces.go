package handler

import (
	pbuser "github.com/MamangRust/microservice-ecommerce-grpc-pb/user"
)

type UserQueryHandler interface {
	pbuser.UserQueryServiceServer
}

type UserCommandHandler interface {
	pbuser.UserCommandServiceServer
}
