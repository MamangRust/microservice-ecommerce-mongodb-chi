package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MamangRust/microservice-ecommerce-grpc-apigateway/apierror"
	gateway_cache "github.com/MamangRust/microservice-ecommerce-grpc-apigateway/cache/gateway_cache"
	"github.com/MamangRust/microservice-ecommerce-grpc-apigateway/httpx"
	pbcategory "github.com/MamangRust/microservice-ecommerce-grpc-pb/category"
	pborder "github.com/MamangRust/microservice-ecommerce-grpc-pb/order"
	pbtransaction "github.com/MamangRust/microservice-ecommerce-grpc-pb/transaction"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/response"
	"github.com/MamangRust/microservice-ecommerce-shared/errors"
	statsmapper "github.com/MamangRust/microservice-ecommerce-shared/mapper/stats"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

// statsHandleApi exposes the stats-reader (ClickHouse) gRPC contracts as REST
// endpoints. It mirrors the full ecommerce stats-reader surface: 7 gRPC
// services / 36 rpc — order, category and transaction basic + ByMerchant +
// ById (category only), plus transaction payment-method stats. There is no
// cashier or order-by-id service in the ecommerce pipeline.
//
// All cache keys use the "apigw:" namespace so the gateway can never overwrite
// domain-service or stats-reader cache keys.
type statsHandleApi struct {
	order                 pborder.OrderStatsServiceClient
	orderByMerchant       pborder.OrderStatsByMerchantServiceClient
	category              pbcategory.CategoryStatsServiceClient
	categoryById          pbcategory.CategoryStatsByIdServiceClient
	categoryByMerchant    pbcategory.CategoryStatsByMerchantServiceClient
	transaction           pbtransaction.TransactionStatsServiceClient
	transactionByMerchant pbtransaction.TransactionStatsByMerchantServiceClient
	logger                logger.LoggerInterface
	mapping               statsmapper.StatsResponseMapper
	apiHandler            apierror.ApiHandler
	cache                 *gateway_cache.GatewayCache
}

// NewHandlerStats builds the 7 stats-reader clients from the shared
// StatsReader connection and registers every route.
func NewHandlerStats(
	router chi.Router,
	conn grpc.ClientConnInterface,
	logger logger.LoggerInterface,
	mapping statsmapper.StatsResponseMapper,
	apiHandler apierror.ApiHandler,
	cache *gateway_cache.GatewayCache,
) *statsHandleApi {
	statsHandler := &statsHandleApi{
		order:                 pborder.NewOrderStatsServiceClient(conn),
		orderByMerchant:       pborder.NewOrderStatsByMerchantServiceClient(conn),
		category:              pbcategory.NewCategoryStatsServiceClient(conn),
		categoryById:          pbcategory.NewCategoryStatsByIdServiceClient(conn),
		categoryByMerchant:    pbcategory.NewCategoryStatsByMerchantServiceClient(conn),
		transaction:           pbtransaction.NewTransactionStatsServiceClient(conn),
		transactionByMerchant: pbtransaction.NewTransactionStatsByMerchantServiceClient(conn),
		logger:                logger,
		mapping:               mapping,
		apiHandler:            apiHandler,
		cache:                 cache,
	}

	router.Route("/api/order-stats", func(orderStats chi.Router) {
		orderStats.Get("/monthly-total-revenue", apiHandler.Handle("get-order-stats-monthly-total-revenue", statsHandler.FindMonthlyTotalRevenue))
		orderStats.Get("/yearly-total-revenue", apiHandler.Handle("get-order-stats-yearly-total-revenue", statsHandler.FindYearlyTotalRevenue))
		orderStats.Get("/monthly-revenue", apiHandler.Handle("get-order-stats-monthly-revenue", statsHandler.FindMonthlyRevenue))
		orderStats.Get("/yearly-revenue", apiHandler.Handle("get-order-stats-yearly-revenue", statsHandler.FindYearlyRevenue))

		orderStats.Route("/merchant", func(m chi.Router) {
			m.Get("/monthly-total-revenue", apiHandler.Handle("get-order-stats-monthly-total-revenue-by-merchant", statsHandler.FindMonthlyTotalRevenueByMerchant))
			m.Get("/yearly-total-revenue", apiHandler.Handle("get-order-stats-yearly-total-revenue-by-merchant", statsHandler.FindYearlyTotalRevenueByMerchant))
			m.Get("/monthly-revenue", apiHandler.Handle("get-order-stats-monthly-revenue-by-merchant", statsHandler.FindMonthlyRevenueByMerchant))
			m.Get("/yearly-revenue", apiHandler.Handle("get-order-stats-yearly-revenue-by-merchant", statsHandler.FindYearlyRevenueByMerchant))
		})
	})

	router.Route("/api/category-stats", func(categoryStats chi.Router) {
		categoryStats.Get("/monthly-total-pricing", apiHandler.Handle("get-category-stats-monthly-total-pricing", statsHandler.FindCategoryMonthlyTotalPrice))
		categoryStats.Get("/yearly-total-pricing", apiHandler.Handle("get-category-stats-yearly-total-pricing", statsHandler.FindCategoryYearlyTotalPrice))
		categoryStats.Get("/monthly-pricing", apiHandler.Handle("get-category-stats-monthly-pricing", statsHandler.FindCategoryMonthPrice))
		categoryStats.Get("/yearly-pricing", apiHandler.Handle("get-category-stats-yearly-pricing", statsHandler.FindCategoryYearPrice))

		categoryStats.Route("/merchant", func(m chi.Router) {
			m.Get("/monthly-total-pricing", apiHandler.Handle("get-category-stats-monthly-total-pricing-by-merchant", statsHandler.FindCategoryMonthlyTotalPriceByMerchant))
			m.Get("/yearly-total-pricing", apiHandler.Handle("get-category-stats-yearly-total-pricing-by-merchant", statsHandler.FindCategoryYearlyTotalPriceByMerchant))
			m.Get("/monthly-pricing", apiHandler.Handle("get-category-stats-monthly-pricing-by-merchant", statsHandler.FindCategoryMonthPriceByMerchant))
			m.Get("/yearly-pricing", apiHandler.Handle("get-category-stats-yearly-pricing-by-merchant", statsHandler.FindCategoryYearPriceByMerchant))
		})

		categoryStats.Route("/mycategory", func(i chi.Router) {
			i.Get("/monthly-total-pricing", apiHandler.Handle("get-category-stats-monthly-total-pricing-by-id", statsHandler.FindCategoryMonthlyTotalPriceById))
			i.Get("/yearly-total-pricing", apiHandler.Handle("get-category-stats-yearly-total-pricing-by-id", statsHandler.FindCategoryYearlyTotalPriceById))
			i.Get("/monthly-pricing", apiHandler.Handle("get-category-stats-monthly-pricing-by-id", statsHandler.FindCategoryMonthPriceById))
			i.Get("/yearly-pricing", apiHandler.Handle("get-category-stats-yearly-pricing-by-id", statsHandler.FindCategoryYearPriceById))
		})
	})

	router.Route("/api/transaction-stats", func(transactionStats chi.Router) {
		transactionStats.Get("/monthly-success", apiHandler.Handle("get-transaction-stats-monthly-success", statsHandler.FindTransactionMonthlySuccess))
		transactionStats.Get("/yearly-success", apiHandler.Handle("get-transaction-stats-yearly-success", statsHandler.FindTransactionYearlySuccess))
		transactionStats.Get("/monthly-failed", apiHandler.Handle("get-transaction-stats-monthly-failed", statsHandler.FindTransactionMonthlyFailed))
		transactionStats.Get("/yearly-failed", apiHandler.Handle("get-transaction-stats-yearly-failed", statsHandler.FindTransactionYearlyFailed))
		transactionStats.Get("/monthly-method-success", apiHandler.Handle("get-transaction-stats-monthly-method-success", statsHandler.FindTransactionMonthMethodSuccess))
		transactionStats.Get("/yearly-method-success", apiHandler.Handle("get-transaction-stats-yearly-method-success", statsHandler.FindTransactionYearMethodSuccess))
		transactionStats.Get("/monthly-method-failed", apiHandler.Handle("get-transaction-stats-monthly-method-failed", statsHandler.FindTransactionMonthMethodFailed))
		transactionStats.Get("/yearly-method-failed", apiHandler.Handle("get-transaction-stats-yearly-method-failed", statsHandler.FindTransactionYearMethodFailed))

		transactionStats.Route("/merchant", func(m chi.Router) {
			m.Get("/monthly-success", apiHandler.Handle("get-transaction-stats-monthly-success-by-merchant", statsHandler.FindTransactionMonthlySuccessByMerchant))
			m.Get("/yearly-success", apiHandler.Handle("get-transaction-stats-yearly-success-by-merchant", statsHandler.FindTransactionYearlySuccessByMerchant))
			m.Get("/monthly-failed", apiHandler.Handle("get-transaction-stats-monthly-failed-by-merchant", statsHandler.FindTransactionMonthlyFailedByMerchant))
			m.Get("/yearly-failed", apiHandler.Handle("get-transaction-stats-yearly-failed-by-merchant", statsHandler.FindTransactionYearlyFailedByMerchant))
			m.Get("/monthly-method-success", apiHandler.Handle("get-transaction-stats-monthly-method-success-by-merchant", statsHandler.FindTransactionMonthMethodByMerchantSuccess))
			m.Get("/yearly-method-success", apiHandler.Handle("get-transaction-stats-yearly-method-success-by-merchant", statsHandler.FindTransactionYearMethodByMerchantSuccess))
			m.Get("/monthly-method-failed", apiHandler.Handle("get-transaction-stats-monthly-method-failed-by-merchant", statsHandler.FindTransactionMonthMethodByMerchantFailed))
			m.Get("/yearly-method-failed", apiHandler.Handle("get-transaction-stats-yearly-method-failed-by-merchant", statsHandler.FindTransactionYearMethodByMerchantFailed))
		})
	})

	return statsHandler
}

func (h *statsHandleApi) parseYearMonth(r *http.Request) (int, int, error) {
	year, err := parseQueryIntWithValidation(r, "year", 1, 9999)
	if err != nil {
		return 0, 0, errors.NewBadRequestError("invalid year")
	}
	month, err := parseQueryIntWithValidation(r, "month", 1, 12)
	if err != nil {
		return 0, 0, errors.NewBadRequestError("invalid month")
	}
	return year, month, nil
}

func (h *statsHandleApi) parseYear(r *http.Request) (int, error) {
	year, err := parseQueryIntWithValidation(r, "year", 1, 9999)
	if err != nil {
		return 0, errors.NewBadRequestError("invalid year")
	}
	return year, nil
}

// parseID validates a positive int32 id query param (merchant_id, category_id).
func (h *statsHandleApi) parseID(r *http.Request, key string) (int, error) {
	id, err := parseQueryIntWithValidation(r, key, 1, 2147483647)
	if err != nil {
		return 0, errors.NewBadRequestError("invalid " + key)
	}
	return id, nil
}

// ── Order stats ─────────────────────────────────────────────────────────────

// @Summary Monthly total revenue (order)
// @Tags Order Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Param month query int true "Month (1-12)"
// @Success 200 {object} response.ApiResponseOrderMonthlyTotalRevenue
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/order-stats/monthly-total-revenue [get]
// @Security BearerAuth
func (h *statsHandleApi) FindMonthlyTotalRevenue(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, month, err := h.parseYearMonth(r)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:order-stats:monthly-total-revenue:%d:%d", year, month)
	if cached, found := gateway_cache.Get[response.ApiResponseOrderMonthlyTotalRevenue](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.order.FindMonthlyTotalRevenue(ctx, &pborder.FindYearMonthTotalRevenue{Year: int32(year), Month: int32(month)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseOrderMonthlyTotalRevenue(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Yearly total revenue (order)
// @Tags Order Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Success 200 {object} response.ApiResponseOrderYearlyTotalRevenue
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/order-stats/yearly-total-revenue [get]
// @Security BearerAuth
func (h *statsHandleApi) FindYearlyTotalRevenue(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := h.parseYear(r)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:order-stats:yearly-total-revenue:%d", year)
	if cached, found := gateway_cache.Get[response.ApiResponseOrderYearlyTotalRevenue](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.order.FindYearlyTotalRevenue(ctx, &pborder.FindYearTotalRevenue{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseOrderYearlyTotalRevenue(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Monthly revenue (order)
// @Tags Order Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Success 200 {object} response.ApiResponseOrderMonthly
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/order-stats/monthly-revenue [get]
// @Security BearerAuth
func (h *statsHandleApi) FindMonthlyRevenue(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := h.parseYear(r)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:order-stats:monthly-revenue:%d", year)
	if cached, found := gateway_cache.Get[response.ApiResponseOrderMonthly](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.order.FindMonthlyRevenue(ctx, &pborder.FindYearOrder{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseOrderMonthly(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Yearly revenue (order)
// @Tags Order Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Success 200 {object} response.ApiResponseOrderYearly
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/order-stats/yearly-revenue [get]
// @Security BearerAuth
func (h *statsHandleApi) FindYearlyRevenue(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := h.parseYear(r)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:order-stats:yearly-revenue:%d", year)
	if cached, found := gateway_cache.Get[response.ApiResponseOrderYearly](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.order.FindYearlyRevenue(ctx, &pborder.FindYearOrder{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseOrderYearly(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Monthly total revenue by merchant (order)
// @Tags Order Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Param month query int true "Month (1-12)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseOrderMonthlyTotalRevenue
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/order-stats/merchant/monthly-total-revenue [get]
// @Security BearerAuth
func (h *statsHandleApi) FindMonthlyTotalRevenueByMerchant(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, month, err := h.parseYearMonth(r)
	if err != nil {
		return err
	}
	merchantID, err := h.parseID(r, "merchant_id")
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:order-stats:monthly-total-revenue:merchant:%d:%d:%d", merchantID, year, month)
	if cached, found := gateway_cache.Get[response.ApiResponseOrderMonthlyTotalRevenue](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.orderByMerchant.FindMonthlyTotalRevenueByMerchant(ctx, &pborder.FindYearMonthTotalRevenueByMerchant{Year: int32(year), Month: int32(month), MerchantId: int32(merchantID)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseOrderMonthlyTotalRevenue(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Yearly total revenue by merchant (order)
// @Tags Order Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseOrderYearlyTotalRevenue
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/order-stats/merchant/yearly-total-revenue [get]
// @Security BearerAuth
func (h *statsHandleApi) FindYearlyTotalRevenueByMerchant(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := h.parseYear(r)
	if err != nil {
		return err
	}
	merchantID, err := h.parseID(r, "merchant_id")
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:order-stats:yearly-total-revenue:merchant:%d:%d", merchantID, year)
	if cached, found := gateway_cache.Get[response.ApiResponseOrderYearlyTotalRevenue](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.orderByMerchant.FindYearlyTotalRevenueByMerchant(ctx, &pborder.FindYearTotalRevenueByMerchant{Year: int32(year), MerchantId: int32(merchantID)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseOrderYearlyTotalRevenue(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Monthly revenue by merchant (order)
// @Tags Order Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseOrderMonthly
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/order-stats/merchant/monthly-revenue [get]
// @Security BearerAuth
func (h *statsHandleApi) FindMonthlyRevenueByMerchant(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := h.parseYear(r)
	if err != nil {
		return err
	}
	merchantID, err := h.parseID(r, "merchant_id")
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:order-stats:monthly-revenue:merchant:%d:%d", merchantID, year)
	if cached, found := gateway_cache.Get[response.ApiResponseOrderMonthly](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.orderByMerchant.FindMonthlyRevenueByMerchant(ctx, &pborder.FindYearOrderByMerchant{MerchantId: int32(merchantID), Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseOrderMonthly(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Yearly revenue by merchant (order)
// @Tags Order Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseOrderYearly
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/order-stats/merchant/yearly-revenue [get]
// @Security BearerAuth
func (h *statsHandleApi) FindYearlyRevenueByMerchant(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := h.parseYear(r)
	if err != nil {
		return err
	}
	merchantID, err := h.parseID(r, "merchant_id")
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:order-stats:yearly-revenue:merchant:%d:%d", merchantID, year)
	if cached, found := gateway_cache.Get[response.ApiResponseOrderYearly](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.orderByMerchant.FindYearlyRevenueByMerchant(ctx, &pborder.FindYearOrderByMerchant{MerchantId: int32(merchantID), Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseOrderYearly(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// ── Category stats ──────────────────────────────────────────────────────────

// @Summary Monthly total pricing (category)
// @Tags Category Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Param month query int true "Month (1-12)"
// @Success 200 {object} response.ApiResponseCategoryMonthlyTotalPrice
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/category-stats/monthly-total-pricing [get]
// @Security BearerAuth
func (h *statsHandleApi) FindCategoryMonthlyTotalPrice(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, month, err := h.parseYearMonth(r)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:category-stats:monthly-total-pricing:%d:%d", year, month)
	if cached, found := gateway_cache.Get[response.ApiResponseCategoryMonthlyTotalPrice](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.category.FindMonthlyTotalPrices(ctx, &pbcategory.FindYearMonthTotalPrices{Year: int32(year), Month: int32(month)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCategoryMonthlyTotalPrice(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Yearly total pricing (category)
// @Tags Category Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Success 200 {object} response.ApiResponseCategoryYearlyTotalPrice
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/category-stats/yearly-total-pricing [get]
// @Security BearerAuth
func (h *statsHandleApi) FindCategoryYearlyTotalPrice(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := h.parseYear(r)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:category-stats:yearly-total-pricing:%d", year)
	if cached, found := gateway_cache.Get[response.ApiResponseCategoryYearlyTotalPrice](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.category.FindYearlyTotalPrices(ctx, &pbcategory.FindYearTotalPrices{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCategoryYearlyTotalPrice(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Monthly pricing (category)
// @Tags Category Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Success 200 {object} response.ApiResponseCategoryMonthPrice
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/category-stats/monthly-pricing [get]
// @Security BearerAuth
func (h *statsHandleApi) FindCategoryMonthPrice(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := h.parseYear(r)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:category-stats:monthly-pricing:%d", year)
	if cached, found := gateway_cache.Get[response.ApiResponseCategoryMonthPrice](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.category.FindMonthPrice(ctx, &pbcategory.FindYearCategory{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCategoryMonthPrice(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Yearly pricing (category)
// @Tags Category Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Success 200 {object} response.ApiResponseCategoryYearPrice
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/category-stats/yearly-pricing [get]
// @Security BearerAuth
func (h *statsHandleApi) FindCategoryYearPrice(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := h.parseYear(r)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:category-stats:yearly-pricing:%d", year)
	if cached, found := gateway_cache.Get[response.ApiResponseCategoryYearPrice](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.category.FindYearPrice(ctx, &pbcategory.FindYearCategory{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCategoryYearPrice(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Monthly total pricing by merchant (category)
// @Tags Category Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Param month query int true "Month (1-12)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseCategoryMonthlyTotalPrice
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/category-stats/merchant/monthly-total-pricing [get]
// @Security BearerAuth
func (h *statsHandleApi) FindCategoryMonthlyTotalPriceByMerchant(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, month, err := h.parseYearMonth(r)
	if err != nil {
		return err
	}
	merchantID, err := h.parseID(r, "merchant_id")
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:category-stats:monthly-total-pricing:merchant:%d:%d:%d", merchantID, year, month)
	if cached, found := gateway_cache.Get[response.ApiResponseCategoryMonthlyTotalPrice](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.categoryByMerchant.FindMonthlyTotalPricesByMerchant(ctx, &pbcategory.FindYearMonthTotalPriceByMerchant{Year: int32(year), Month: int32(month), MerchantId: int32(merchantID)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCategoryMonthlyTotalPrice(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Yearly total pricing by merchant (category)
// @Tags Category Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseCategoryYearlyTotalPrice
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/category-stats/merchant/yearly-total-pricing [get]
// @Security BearerAuth
func (h *statsHandleApi) FindCategoryYearlyTotalPriceByMerchant(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := h.parseYear(r)
	if err != nil {
		return err
	}
	merchantID, err := h.parseID(r, "merchant_id")
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:category-stats:yearly-total-pricing:merchant:%d:%d", merchantID, year)
	if cached, found := gateway_cache.Get[response.ApiResponseCategoryYearlyTotalPrice](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.categoryByMerchant.FindYearlyTotalPricesByMerchant(ctx, &pbcategory.FindYearTotalPriceByMerchant{Year: int32(year), MerchantId: int32(merchantID)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCategoryYearlyTotalPrice(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Monthly pricing by merchant (category)
// @Tags Category Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseCategoryMonthPrice
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/category-stats/merchant/monthly-pricing [get]
// @Security BearerAuth
func (h *statsHandleApi) FindCategoryMonthPriceByMerchant(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := h.parseYear(r)
	if err != nil {
		return err
	}
	merchantID, err := h.parseID(r, "merchant_id")
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:category-stats:monthly-pricing:merchant:%d:%d", merchantID, year)
	if cached, found := gateway_cache.Get[response.ApiResponseCategoryMonthPrice](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.categoryByMerchant.FindMonthPriceByMerchant(ctx, &pbcategory.FindYearCategoryByMerchant{MerchantId: int32(merchantID), Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCategoryMonthPrice(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Yearly pricing by merchant (category)
// @Tags Category Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseCategoryYearPrice
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/category-stats/merchant/yearly-pricing [get]
// @Security BearerAuth
func (h *statsHandleApi) FindCategoryYearPriceByMerchant(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := h.parseYear(r)
	if err != nil {
		return err
	}
	merchantID, err := h.parseID(r, "merchant_id")
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:category-stats:yearly-pricing:merchant:%d:%d", merchantID, year)
	if cached, found := gateway_cache.Get[response.ApiResponseCategoryYearPrice](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.categoryByMerchant.FindYearPriceByMerchant(ctx, &pbcategory.FindYearCategoryByMerchant{MerchantId: int32(merchantID), Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCategoryYearPrice(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Monthly total pricing by category (category)
// @Tags Category Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Param month query int true "Month (1-12)"
// @Param category_id query int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryMonthlyTotalPrice
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/category-stats/mycategory/monthly-total-pricing [get]
// @Security BearerAuth
func (h *statsHandleApi) FindCategoryMonthlyTotalPriceById(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, month, err := h.parseYearMonth(r)
	if err != nil {
		return err
	}
	categoryID, err := h.parseID(r, "category_id")
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:category-stats:monthly-total-pricing:id:%d:%d:%d", categoryID, year, month)
	if cached, found := gateway_cache.Get[response.ApiResponseCategoryMonthlyTotalPrice](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.categoryById.FindMonthlyTotalPricesById(ctx, &pbcategory.FindYearMonthTotalPriceById{Year: int32(year), Month: int32(month), CategoryId: int32(categoryID)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCategoryMonthlyTotalPrice(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Yearly total pricing by category (category)
// @Tags Category Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Param category_id query int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryYearlyTotalPrice
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/category-stats/mycategory/yearly-total-pricing [get]
// @Security BearerAuth
func (h *statsHandleApi) FindCategoryYearlyTotalPriceById(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := h.parseYear(r)
	if err != nil {
		return err
	}
	categoryID, err := h.parseID(r, "category_id")
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:category-stats:yearly-total-pricing:id:%d:%d", categoryID, year)
	if cached, found := gateway_cache.Get[response.ApiResponseCategoryYearlyTotalPrice](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.categoryById.FindYearlyTotalPricesById(ctx, &pbcategory.FindYearTotalPriceById{Year: int32(year), CategoryId: int32(categoryID)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCategoryYearlyTotalPrice(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Monthly pricing by category (category)
// @Tags Category Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Param category_id query int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryMonthPrice
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/category-stats/mycategory/monthly-pricing [get]
// @Security BearerAuth
func (h *statsHandleApi) FindCategoryMonthPriceById(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := h.parseYear(r)
	if err != nil {
		return err
	}
	categoryID, err := h.parseID(r, "category_id")
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:category-stats:monthly-pricing:id:%d:%d", categoryID, year)
	if cached, found := gateway_cache.Get[response.ApiResponseCategoryMonthPrice](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.categoryById.FindMonthPriceById(ctx, &pbcategory.FindYearCategoryById{CategoryId: int32(categoryID), Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCategoryMonthPrice(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Yearly pricing by category (category)
// @Tags Category Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Param category_id query int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryYearPrice
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/category-stats/mycategory/yearly-pricing [get]
// @Security BearerAuth
func (h *statsHandleApi) FindCategoryYearPriceById(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := h.parseYear(r)
	if err != nil {
		return err
	}
	categoryID, err := h.parseID(r, "category_id")
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:category-stats:yearly-pricing:id:%d:%d", categoryID, year)
	if cached, found := gateway_cache.Get[response.ApiResponseCategoryYearPrice](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.categoryById.FindYearPriceById(ctx, &pbcategory.FindYearCategoryById{CategoryId: int32(categoryID), Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCategoryYearPrice(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// ── Transaction status stats ────────────────────────────────────────────────

// @Summary Monthly success (transaction)
// @Tags Transaction Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Param month query int true "Month (1-12)"
// @Success 200 {object} response.ApiResponsesTransactionMonthSuccess
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/transaction-stats/monthly-success [get]
// @Security BearerAuth
func (h *statsHandleApi) FindTransactionMonthlySuccess(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, month, err := h.parseYearMonth(r)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:transaction-stats:monthly-success:%d:%d", year, month)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionMonthSuccess](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.transaction.GetMonthlyAmountSuccess(ctx, &pbtransaction.MonthAmountTransactionRequest{Year: int32(year), Month: int32(month)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionMonthSuccess(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Yearly success (transaction)
// @Tags Transaction Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Success 200 {object} response.ApiResponsesTransactionYearSuccess
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/transaction-stats/yearly-success [get]
// @Security BearerAuth
func (h *statsHandleApi) FindTransactionYearlySuccess(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := h.parseYear(r)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:transaction-stats:yearly-success:%d", year)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionYearSuccess](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.transaction.GetYearlyAmountSuccess(ctx, &pbtransaction.YearAmountTransactionRequest{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionYearSuccess(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Monthly failed (transaction)
// @Tags Transaction Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Param month query int true "Month (1-12)"
// @Success 200 {object} response.ApiResponsesTransactionMonthFailed
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/transaction-stats/monthly-failed [get]
// @Security BearerAuth
func (h *statsHandleApi) FindTransactionMonthlyFailed(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, month, err := h.parseYearMonth(r)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:transaction-stats:monthly-failed:%d:%d", year, month)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionMonthFailed](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.transaction.GetMonthlyAmountFailed(ctx, &pbtransaction.MonthAmountTransactionRequest{Year: int32(year), Month: int32(month)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionMonthFailed(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Yearly failed (transaction)
// @Tags Transaction Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Success 200 {object} response.ApiResponsesTransactionYearFailed
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/transaction-stats/yearly-failed [get]
// @Security BearerAuth
func (h *statsHandleApi) FindTransactionYearlyFailed(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := h.parseYear(r)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:transaction-stats:yearly-failed:%d", year)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionYearFailed](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.transaction.GetYearlyAmountFailed(ctx, &pbtransaction.YearAmountTransactionRequest{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionYearFailed(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Monthly success by merchant (transaction)
// @Tags Transaction Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Param month query int true "Month (1-12)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionMonthSuccess
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/transaction-stats/merchant/monthly-success [get]
// @Security BearerAuth
func (h *statsHandleApi) FindTransactionMonthlySuccessByMerchant(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, month, err := h.parseYearMonth(r)
	if err != nil {
		return err
	}
	merchantID, err := h.parseID(r, "merchant_id")
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:transaction-stats:monthly-success:merchant:%d:%d:%d", merchantID, year, month)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionMonthSuccess](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.transactionByMerchant.GetMonthlyAmountSuccessByMerchant(ctx, &pbtransaction.MonthAmountTransactionMerchantRequest{MerchantId: int32(merchantID), Year: int32(year), Month: int32(month)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionMonthSuccess(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Yearly success by merchant (transaction)
// @Tags Transaction Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionYearSuccess
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/transaction-stats/merchant/yearly-success [get]
// @Security BearerAuth
func (h *statsHandleApi) FindTransactionYearlySuccessByMerchant(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := h.parseYear(r)
	if err != nil {
		return err
	}
	merchantID, err := h.parseID(r, "merchant_id")
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:transaction-stats:yearly-success:merchant:%d:%d", merchantID, year)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionYearSuccess](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.transactionByMerchant.GetYearlyAmountSuccessByMerchant(ctx, &pbtransaction.YearAmountTransactionMerchantRequest{MerchantId: int32(merchantID), Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionYearSuccess(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Monthly failed by merchant (transaction)
// @Tags Transaction Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Param month query int true "Month (1-12)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionMonthFailed
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/transaction-stats/merchant/monthly-failed [get]
// @Security BearerAuth
func (h *statsHandleApi) FindTransactionMonthlyFailedByMerchant(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, month, err := h.parseYearMonth(r)
	if err != nil {
		return err
	}
	merchantID, err := h.parseID(r, "merchant_id")
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:transaction-stats:monthly-failed:merchant:%d:%d:%d", merchantID, year, month)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionMonthFailed](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.transactionByMerchant.GetMonthlyAmountFailedByMerchant(ctx, &pbtransaction.MonthAmountTransactionMerchantRequest{MerchantId: int32(merchantID), Year: int32(year), Month: int32(month)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionMonthFailed(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Yearly failed by merchant (transaction)
// @Tags Transaction Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionYearFailed
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/transaction-stats/merchant/yearly-failed [get]
// @Security BearerAuth
func (h *statsHandleApi) FindTransactionYearlyFailedByMerchant(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := h.parseYear(r)
	if err != nil {
		return err
	}
	merchantID, err := h.parseID(r, "merchant_id")
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:transaction-stats:yearly-failed:merchant:%d:%d", merchantID, year)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionYearFailed](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.transactionByMerchant.GetYearlyAmountFailedByMerchant(ctx, &pbtransaction.YearAmountTransactionMerchantRequest{MerchantId: int32(merchantID), Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionYearFailed(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// ── Transaction payment-method stats ────────────────────────────────────────

// @Summary Monthly success by payment method (transaction)
// @Tags Transaction Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Param month query int true "Month (1-12)"
// @Success 200 {object} response.ApiResponsesTransactionMonthMethod
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/transaction-stats/monthly-method-success [get]
// @Security BearerAuth
func (h *statsHandleApi) FindTransactionMonthMethodSuccess(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, month, err := h.parseYearMonth(r)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:transaction-method-stats:monthly-success:%d:%d", year, month)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionMonthMethod](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.transaction.GetMonthlyTransactionMethodSuccess(ctx, &pbtransaction.MonthMethodTransactionRequest{Year: int32(year), Month: int32(month)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionMonthMethod(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Yearly success by payment method (transaction)
// @Tags Transaction Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Success 200 {object} response.ApiResponsesTransactionYearMethod
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/transaction-stats/yearly-method-success [get]
// @Security BearerAuth
func (h *statsHandleApi) FindTransactionYearMethodSuccess(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := h.parseYear(r)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:transaction-method-stats:yearly-success:%d", year)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionYearMethod](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.transaction.GetYearlyTransactionMethodSuccess(ctx, &pbtransaction.YearMethodTransactionRequest{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionYearMethod(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Monthly failed by payment method (transaction)
// @Tags Transaction Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Param month query int true "Month (1-12)"
// @Success 200 {object} response.ApiResponsesTransactionMonthMethod
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/transaction-stats/monthly-method-failed [get]
// @Security BearerAuth
func (h *statsHandleApi) FindTransactionMonthMethodFailed(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, month, err := h.parseYearMonth(r)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:transaction-method-stats:monthly-failed:%d:%d", year, month)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionMonthMethod](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.transaction.GetMonthlyTransactionMethodFailed(ctx, &pbtransaction.MonthMethodTransactionRequest{Year: int32(year), Month: int32(month)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionMonthMethod(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Yearly failed by payment method (transaction)
// @Tags Transaction Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Success 200 {object} response.ApiResponsesTransactionYearMethod
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/transaction-stats/yearly-method-failed [get]
// @Security BearerAuth
func (h *statsHandleApi) FindTransactionYearMethodFailed(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := h.parseYear(r)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:transaction-method-stats:yearly-failed:%d", year)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionYearMethod](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.transaction.GetYearlyTransactionMethodFailed(ctx, &pbtransaction.YearMethodTransactionRequest{Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionYearMethod(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Monthly success by payment method and merchant (transaction)
// @Tags Transaction Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Param month query int true "Month (1-12)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionMonthMethod
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/transaction-stats/merchant/monthly-method-success [get]
// @Security BearerAuth
func (h *statsHandleApi) FindTransactionMonthMethodByMerchantSuccess(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, month, err := h.parseYearMonth(r)
	if err != nil {
		return err
	}
	merchantID, err := h.parseID(r, "merchant_id")
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:transaction-method-stats:monthly-success:merchant:%d:%d:%d", merchantID, year, month)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionMonthMethod](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.transactionByMerchant.GetMonthlyTransactionMethodByMerchantSuccess(ctx, &pbtransaction.MonthMethodTransactionMerchantRequest{MerchantId: int32(merchantID), Year: int32(year), Month: int32(month)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionMonthMethod(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Yearly success by payment method and merchant (transaction)
// @Tags Transaction Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionYearMethod
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/transaction-stats/merchant/yearly-method-success [get]
// @Security BearerAuth
func (h *statsHandleApi) FindTransactionYearMethodByMerchantSuccess(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := h.parseYear(r)
	if err != nil {
		return err
	}
	merchantID, err := h.parseID(r, "merchant_id")
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:transaction-method-stats:yearly-success:merchant:%d:%d", merchantID, year)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionYearMethod](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.transactionByMerchant.GetYearlyTransactionMethodByMerchantSuccess(ctx, &pbtransaction.YearMethodTransactionMerchantRequest{MerchantId: int32(merchantID), Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionYearMethod(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Monthly failed by payment method and merchant (transaction)
// @Tags Transaction Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Param month query int true "Month (1-12)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionMonthMethod
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/transaction-stats/merchant/monthly-method-failed [get]
// @Security BearerAuth
func (h *statsHandleApi) FindTransactionMonthMethodByMerchantFailed(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, month, err := h.parseYearMonth(r)
	if err != nil {
		return err
	}
	merchantID, err := h.parseID(r, "merchant_id")
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:transaction-method-stats:monthly-failed:merchant:%d:%d:%d", merchantID, year, month)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionMonthMethod](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.transactionByMerchant.GetMonthlyTransactionMethodByMerchantFailed(ctx, &pbtransaction.MonthMethodTransactionMerchantRequest{MerchantId: int32(merchantID), Year: int32(year), Month: int32(month)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionMonthMethod(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// @Summary Yearly failed by payment method and merchant (transaction)
// @Tags Transaction Stats
// @Accept json
// @Produce json
// @Param year query int true "Year (e.g. 2026)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionYearMethod
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/transaction-stats/merchant/yearly-method-failed [get]
// @Security BearerAuth
func (h *statsHandleApi) FindTransactionYearMethodByMerchantFailed(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := h.parseYear(r)
	if err != nil {
		return err
	}
	merchantID, err := h.parseID(r, "merchant_id")
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("apigw:transaction-method-stats:yearly-failed:merchant:%d:%d", merchantID, year)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionYearMethod](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.transactionByMerchant.GetYearlyTransactionMethodByMerchantFailed(ctx, &pbtransaction.YearMethodTransactionMerchantRequest{MerchantId: int32(merchantID), Year: int32(year)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionYearMethod(res)
	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}
