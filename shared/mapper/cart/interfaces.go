package cartapimapper

import (
	pbcart "github.com/MamangRust/microservice-ecommerce-grpc-pb/cart"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/response"
)

type CartBaseResponseMapper interface {
	ToResponseCart(pbResponse *pbcart.CartResponse) *response.CartResponse
	ToResponseCarts(pbResponse []*pbcart.CartResponse) []*response.CartResponse
	ToApiResponseCart(pbResponse *pbcart.ApiResponseCart) *response.ApiResponseCart
}

type CartQueryResponseMapper interface {
	CartBaseResponseMapper
	ToApiResponseCartPagination(pbResponse *pbcart.ApiResponsePaginationCart) *response.ApiResponseCartPagination
}

type CartCommandResponseMapper interface {
	CartBaseResponseMapper
	ToApiResponseCartDelete(pbResponse *pbcart.ApiResponseCartDelete) *response.ApiResponseCartDelete
	ToApiResponseCartAll(pbResponse *pbcart.ApiResponseCartAll) *response.ApiResponseCartAll
}
