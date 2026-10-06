package orderapimapper

import (
	pborder "github.com/MamangRust/microservice-ecommerce-grpc-pb/order"
	pborderstats "github.com/MamangRust/microservice-ecommerce-grpc-pb/order"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/response"
)

type OrderBaseResponseMapper interface {
	ToResponseOrder(order *pborder.OrderResponse) *response.OrderResponse
	ToResponsesOrder(orders []*pborder.OrderResponse) []*response.OrderResponse
	ToApiResponseOrder(pbResponse *pborder.ApiResponseOrder) *response.ApiResponseOrder
}

type OrderQueryResponseMapper interface {
	OrderBaseResponseMapper
	ToApiResponsesOrder(pbResponse *pborder.ApiResponsesOrder) *response.ApiResponsesOrder
	ToApiResponsePaginationOrder(pbResponse *pborder.ApiResponsePaginationOrder) *response.ApiResponsePaginationOrder
	ToApiResponsePaginationOrderDeleteAt(pbResponse *pborder.ApiResponsePaginationOrderDeleteAt) *response.ApiResponsePaginationOrderDeleteAt
}

type OrderCommandResponseMapper interface {
	OrderBaseResponseMapper
	ToResponseOrderDeleteAt(order *pborder.OrderResponseDeleteAt) *response.OrderResponseDeleteAt
	ToResponsesOrderDeleteAt(orders []*pborder.OrderResponseDeleteAt) []*response.OrderResponseDeleteAt
	ToApiResponseOrderDeleteAt(pbResponse *pborder.ApiResponseOrderDeleteAt) *response.ApiResponseOrderDeleteAt
	ToApiResponseOrderDelete(pbResponse *pborder.ApiResponseOrderDelete) *response.ApiResponseOrderDelete
	ToApiResponseOrderAll(pbResponse *pborder.ApiResponseOrderAll) *response.ApiResponseOrderAll
}

type OrderStatsResponseMapper interface {
	ToOrderMonthlyPrice(category *pborderstats.OrderMonthlyResponse) *response.OrderMonthlyResponse
	ToOrderMonthlyPrices(c []*pborderstats.OrderMonthlyResponse) []*response.OrderMonthlyResponse
	ToOrderYearlyPrice(category *pborderstats.OrderYearlyResponse) *response.OrderYearlyResponse
	ToOrderYearlyPrices(c []*pborderstats.OrderYearlyResponse) []*response.OrderYearlyResponse
	ToResponseOrderMonthlyTotalRevenue(c *pborderstats.OrderMonthlyTotalRevenueResponse) *response.OrderMonthlyTotalRevenueResponse
	ToResponseOrderMonthlyTotalRevenues(c []*pborderstats.OrderMonthlyTotalRevenueResponse) []*response.OrderMonthlyTotalRevenueResponse
	ToResponseOrderYearlyTotalRevenue(c *pborderstats.OrderYearlyTotalRevenueResponse) *response.OrderYearlyTotalRevenueResponse
	ToResponseOrderYearlyTotalRevenues(c []*pborderstats.OrderYearlyTotalRevenueResponse) []*response.OrderYearlyTotalRevenueResponse

	ToApiResponseMonthlyOrder(pbResponse *pborderstats.ApiResponseOrderMonthly) *response.ApiResponseOrderMonthly
	ToApiResponseYearlyOrder(pbResponse *pborderstats.ApiResponseOrderYearly) *response.ApiResponseOrderYearly
	ToApiResponseMonthlyTotalRevenue(pbResponse *pborderstats.ApiResponseOrderMonthlyTotalRevenue) *response.ApiResponseOrderMonthlyTotalRevenue
	ToApiResponseYearlyTotalRevenue(pbResponse *pborderstats.ApiResponseOrderYearlyTotalRevenue) *response.ApiResponseOrderYearlyTotalRevenue
}
