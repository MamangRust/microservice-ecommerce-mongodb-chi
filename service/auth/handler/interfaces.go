package handler

import (
	"context"

	pbauth "github.com/MamangRust/microservice-ecommerce-grpc-pb/auth"
)

type AuthHandleGrpc interface {
	pbauth.AuthServiceServer
	LoginUser(ctx context.Context, req *pbauth.LoginRequest) (*pbauth.ApiResponseLogin, error)
	RegisterUser(ctx context.Context, req *pbauth.RegisterRequest) (*pbauth.ApiResponseRegister, error)
}
