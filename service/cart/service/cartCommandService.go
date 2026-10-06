package service

import (
	"context"

	mencache "github.com/MamangRust/microservice-ecommerce-grpc-cart/cache"
	"github.com/MamangRust/microservice-ecommerce-grpc-cart/repository"
	productadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/product"
	useradapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/user"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errorhandler"
	"github.com/MamangRust/microservice-ecommerce-shared/errors"
	"github.com/MamangRust/microservice-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type cartCommandService struct {
	cartCommandRepository  repository.CartCommandRepository
	productQueryRepository productadapter.QueryRepository
	userQueryRepository    useradapter.QueryRepository
	mencache               mencache.CartMencache
	observability          observability.TraceLoggerObservability
	logger                 logger.LoggerInterface
}

type CartCommandServiceDeps struct {
	CartCommandRepository  repository.CartCommandRepository
	ProductQueryRepository productadapter.QueryRepository
	UserQueryRepository    useradapter.QueryRepository
	Mencache               mencache.CartMencache
	Observability          observability.TraceLoggerObservability
	Logger                 logger.LoggerInterface
}

func NewCartCommandService(deps *CartCommandServiceDeps) *cartCommandService {
	return &cartCommandService{
		cartCommandRepository:  deps.CartCommandRepository,
		productQueryRepository: deps.ProductQueryRepository,
		userQueryRepository:    deps.UserQueryRepository,
		mencache:               deps.Mencache,
		logger:                 deps.Logger,
		observability:          deps.Observability,
	}
}

func (s *cartCommandService) Create(ctx context.Context, req *requests.CreateCartRequest) (*repository.CartCreateResult, error) {
	const method = "Create"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	product, err := s.productQueryRepository.FindByID(ctx, req.ProductID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*repository.CartCreateResult](
			s.logger,
			err,
			method,
			span,
			zap.Int("product_id", req.ProductID),
		)
	}

	if req.Quantity <= 0 {
		status = "error"
		return errorhandler.HandleError[*repository.CartCreateResult](s.logger, errors.ErrBadRequest.WithMessage("cart quantity must be greater than zero"), method, span)
	}
	if product.CountInStock < int32(req.Quantity) {
		status = "error"
		return errorhandler.HandleError[*repository.CartCreateResult](s.logger, errors.ErrBadRequest.WithMessage("Insufficient product stock"), method, span)
	}

	_, err = s.userQueryRepository.FindByID(ctx, req.UserID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*repository.CartCreateResult](
			s.logger,
			err,
			method,
			span,
			zap.Int("user_id", req.UserID),
		)
	}

	var imageProduct string
	if product.ImageProduct != nil {
		imageProduct = *product.ImageProduct
	}

	var weight int
	if product.Weight != nil {
		weight = int(*product.Weight)
	}

	cartRecord := &requests.CartCreateRecord{
		ProductID:    req.ProductID,
		UserID:       req.UserID,
		Name:         product.Name,
		Price:        int(product.Price),
		ImageProduct: imageProduct,
		Quantity:     req.Quantity,
		Weight:       weight,
	}

	res, err := s.cartCommandRepository.CreateCart(ctx, cartRecord)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*repository.CartCreateResult](
			s.logger,
			err,
			method,
			span,
			zap.Any("request", req),
		)
	}

	s.mencache.DeleteCartsCache(ctx, req.UserID)

	logSuccess("Successfully created cart", zap.Int("cartID", int(res.CartID)))
	return res, nil
}

func (s *cartCommandService) DeletePermanent(ctx context.Context, req *requests.DeleteCartRequest) (bool, error) {
	const method = "DeletePermanentCart"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("cartID", req.CartID))

	defer func() {
		end(status)
	}()

	success, err := s.cartCommandRepository.DeletePermanent(ctx, &requests.DeleteCartRequest{
		CartID: req.CartID,
		UserID: req.UserID,
	})
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			err,
			method,
			span,
			zap.Int("cart_id", req.CartID),
		)
	}

	s.mencache.DeleteCartsCache(ctx, req.UserID)

	logSuccess("Successfully deleted cart permanently", zap.Int("cartID", req.CartID))
	return success, nil
}

func (s *cartCommandService) DeleteAll(ctx context.Context, req *requests.DeleteAllCartRequest) (bool, error) {
	const method = "DeleteAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.cartCommandRepository.DeleteAllPermanently(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			err,
			method,
			span,
		)
	}

	s.mencache.DeleteCartsCache(ctx, req.UserID)

	logSuccess("Successfully deleted all carts permanently")
	return success, nil
}
