package authapimapper

import (
	pbauth "github.com/MamangRust/microservice-ecommerce-grpc-pb/auth"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/response"
)

type AuthBaseResponseMapper interface {
}

type AuthQueryResponseMapper interface {
	AuthBaseResponseMapper
	ToResponseGetMe(res *pbauth.ApiResponseGetMe) *response.ApiResponseGetMe
}

type AuthCommandResponseMapper interface {
	AuthBaseResponseMapper
	ToResponseVerifyCode(res *pbauth.ApiResponseVerifyCode) *response.ApiResponseVerifyCode
	ToResponseForgotPassword(res *pbauth.ApiResponseForgotPassword) *response.ApiResponseForgotPassword
	ToResponseResetPassword(res *pbauth.ApiResponseResetPassword) *response.ApiResponseResetPassword
	ToResponseLogin(res *pbauth.ApiResponseLogin) *response.ApiResponseLogin
	ToResponseRegister(res *pbauth.ApiResponseRegister) *response.ApiResponseRegister
	ToResponseRefreshToken(res *pbauth.ApiResponseRefreshToken) *response.ApiResponseRefreshToken
}
