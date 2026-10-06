package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/MamangRust/microservice-ecommerce-grpc-transaction/cache"
	"github.com/MamangRust/microservice-ecommerce-grpc-transaction/repository"
	merchantadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/merchant"
	orderadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/order"
	orderitemadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/order_item"
	shippingadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/shipping_address"
	useradapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/user"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-pkg/email"
	"github.com/MamangRust/microservice-ecommerce-pkg/event"
	"github.com/MamangRust/microservice-ecommerce-pkg/kafka"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/events"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errorhandler"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/transaction_errors"
	"github.com/google/uuid"

	"github.com/MamangRust/microservice-ecommerce-shared/observability"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type transactionCommandService struct {
	observability      observability.TraceLoggerObservability
	kafka              *kafka.Kafka
	outbox             repository.OutboxRepository
	db                 *mongo.Database
	cache              cache.TransactionCommandCache
	transactionQuery   repository.TransactionQueryRepository
	transactionCommand repository.TransactionCommandRepository
	userQuery          useradapter.QueryRepository
	merchantQuery      merchantadapter.QueryRepository
	orderQuery         orderadapter.QueryRepository
	orderItem          orderitemadapter.QueryRepository
	shippingAddress    shippingadapter.QueryRepository
	logger             logger.LoggerInterface
}

type TransactionCommandServiceDeps struct {
	Observability      observability.TraceLoggerObservability
	Kafka              *kafka.Kafka
	Outbox             repository.OutboxRepository
	DB                 *mongo.Database
	Cache              cache.TransactionCommandCache
	TransactionQuery   repository.TransactionQueryRepository
	TransactionCommand repository.TransactionCommandRepository
	UserQuery          useradapter.QueryRepository
	MerchantQuery      merchantadapter.QueryRepository
	OrderQuery         orderadapter.QueryRepository
	OrderItem          orderitemadapter.QueryRepository
	ShippingAddress    shippingadapter.QueryRepository
	Logger             logger.LoggerInterface
}

func NewTransactionCommandService(deps *TransactionCommandServiceDeps) TransactionCommandService {
	return &transactionCommandService{
		observability:      deps.Observability,
		kafka:              deps.Kafka,
		outbox:             deps.Outbox,
		db:                 deps.DB,
		cache:              deps.Cache,
		transactionQuery:   deps.TransactionQuery,
		transactionCommand: deps.TransactionCommand,
		userQuery:          deps.UserQuery,
		merchantQuery:      deps.MerchantQuery,
		orderQuery:         deps.OrderQuery,
		orderItem:          deps.OrderItem,
		shippingAddress:    deps.ShippingAddress,
		logger:             deps.Logger,
	}
}

func (s *transactionCommandService) Create(ctx context.Context, req *requests.CreateTransactionRequest) (*repository.TransactionResult, error) {
	const method = "Create"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user_id", req.UserID),
		attribute.Int("merchant_id", req.MerchantID),
		attribute.Int("order_id", req.OrderID))

	defer func() {
		end(status)
	}()

	user, err := s.userQuery.FindByID(ctx, req.UserID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*repository.TransactionResult](s.logger, err, method, span)
	}

	_, err = s.merchantQuery.FindByID(ctx, req.MerchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*repository.TransactionResult](s.logger, err, method, span)
	}

	_, err = s.orderQuery.FindByID(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*repository.TransactionResult](s.logger, err, method, span)
	}

	orderItems, err := s.orderItem.FindOrderItemByOrder(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*repository.TransactionResult](s.logger, err, method, span)
	}

	shipping, err := s.shippingAddress.FindByOrder(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*repository.TransactionResult](s.logger, err, method, span)
	}

	var merchandiseTotal int
	for _, item := range orderItems {
		if item.Quantity <= 0 || item.Price < 0 {
			status = "error"
			return errorhandler.HandleError[*repository.TransactionResult](s.logger, transaction_errors.ErrFailedOrderItemEmpty, method, span)
		}
		merchandiseTotal += int(item.Price) * int(item.Quantity)
	}

	totalAmount := merchandiseTotal + int(shipping.ShippingCost)
	ppn := totalAmount * 11 / 100
	totalAmountWithTax := totalAmount + ppn

	span.SetAttributes(attribute.Int("calculated_amount", totalAmountWithTax))

	if req.PaymentStatus != nil && *req.PaymentStatus != "" {
		if !IsValidPaymentStatus(*req.PaymentStatus) {
			status = "error"
			return errorhandler.HandleError[*repository.TransactionResult](s.logger, transaction_errors.ErrFailedPaymentStatusInvalid, method, span)
		}
	}

	var paymentStatus string
	if req.Amount >= totalAmountWithTax {
		paymentStatus = PaymentStatusSuccess
	} else {
		status = "error"
		return errorhandler.HandleError[*repository.TransactionResult](s.logger, transaction_errors.ErrFailedPaymentInsufficientBalance, method, span)
	}

	if req.PaymentStatus != nil && *req.PaymentStatus != "" && !CanTransitionPaymentStatus(*req.PaymentStatus, paymentStatus) {
		status = "error"
		return errorhandler.HandleError[*repository.TransactionResult](s.logger, transaction_errors.ErrFailedPaymentStatusCannotBeModified, method, span)
	}

	req.Amount = totalAmountWithTax
	req.PaymentStatus = &paymentStatus

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   "Transaction Successful",
		"Message": fmt.Sprintf("Your transaction of %d has been processed successfully.", req.Amount),
		"Button":  "View History",
		"Link":    "https://sanedge.example.com/transaction/history",
	})

	payloadBytes, err := event.MarshalEmail("transaction.created", user.Email, "Transaction Successful - SanEdge", htmlBody)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*repository.TransactionResult](s.logger, err, method, span)
	}

	var transaction *repository.TransactionResult
	transaction, err = s.transactionCommand.Create(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*repository.TransactionResult](s.logger, err, method, span)
	}

	if s.outbox != nil {
		merchantPayload, marshalErr := json.Marshal(map[string]any{
			"merchantId":    transaction.MerchantID,
			"transactionId": transaction.TransactionID,
			"amount":        transaction.Amount,
			"status":        transaction.PaymentStatus,
			"timestamp":     time.Now().UnixMilli(),
		})
		if marshalErr != nil {
			status = "error"
			return errorhandler.HandleError[*repository.TransactionResult](s.logger, marshalErr, method, span)
		}

		statsEvent, statsErr := json.Marshal(events.StatsEnvelope{
			EventID: uuid.NewString(),
			Payload: mustJSON(events.TransactionEvent{
				TransactionID: transaction.TransactionID,
				OrderID:       transaction.OrderID,
				MerchantID:    transaction.MerchantID,
				PaymentMethod: transaction.PaymentMethod,
				Amount:        transaction.Amount,
				Status:        transaction.PaymentStatus,
				EventTime:     time.Now().UTC().Format(time.RFC3339),
			}),
		})
		if statsErr != nil {
			status = "error"
			return errorhandler.HandleError[*repository.TransactionResult](s.logger, statsErr, method, span)
		}

		enqueue := func(topic, key string, payload []byte) error {
			_, enqueueErr := s.outbox.Create(ctx, topic, key, payload)
			return enqueueErr
		}

		for _, event := range []struct {
			topic   string
			key     string
			payload []byte
		}{
			{topic: "email-service-topic-transaction-create", key: strconv.Itoa(int(transaction.TransactionID)), payload: payloadBytes},
			{topic: "merchant-service-topic-transaction-event", key: strconv.Itoa(int(transaction.MerchantID)), payload: merchantPayload},
			{topic: "stats.ecommerce.transaction.event", key: strconv.Itoa(int(transaction.TransactionID)), payload: statsEvent},
		} {
			if enqueueErr := enqueue(event.topic, event.key, event.payload); enqueueErr != nil {
				s.logger.Error("failed to enqueue outbox event", zap.Error(enqueueErr), zap.String("topic", event.topic), zap.Int32("transaction_id", transaction.TransactionID))
			}
		}
	}

	s.cache.InvalidateTransactionCache(ctx)

	logSuccess("Successfully created transaction", zap.Int32("transaction_id", transaction.TransactionID))

	return transaction, nil
}

func (s *transactionCommandService) Update(ctx context.Context, req *requests.UpdateTransactionRequest) (*repository.TransactionResult, error) {
	const method = "Update"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transaction_id", *req.TransactionID))

	defer func() {
		end(status)
	}()

	existingTx, err := s.transactionQuery.FindByID(ctx, *req.TransactionID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*repository.TransactionResult](s.logger, err, method, span)
	}

	if req.PaymentStatus != nil && *req.PaymentStatus != "" {
		if !IsValidPaymentStatus(*req.PaymentStatus) {
			status = "error"
			return errorhandler.HandleError[*repository.TransactionResult](s.logger, transaction_errors.ErrFailedPaymentStatusInvalid, method, span)
		}
		if *req.PaymentStatus != existingTx.PaymentStatus && !CanTransitionPaymentStatus(existingTx.PaymentStatus, *req.PaymentStatus) {
			status = "error"
			return errorhandler.HandleError[*repository.TransactionResult](s.logger, transaction_errors.ErrFailedPaymentStatusCannotBeModified, method, span)
		}
	}

	if req.MerchantID == 0 {
		req.MerchantID = int(existingTx.MerchantID)
	}
	_, err = s.merchantQuery.FindByID(ctx, req.MerchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*repository.TransactionResult](s.logger, err, method, span)
	}

	if req.OrderID == 0 {
		req.OrderID = int(existingTx.OrderID)
	}
	_, err = s.orderQuery.FindByID(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*repository.TransactionResult](s.logger, err, method, span)
	}

	orderItems, err := s.orderItem.FindOrderItemByOrder(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*repository.TransactionResult](s.logger, err, method, span)
	}

	shipping, err := s.shippingAddress.FindByOrder(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*repository.TransactionResult](s.logger, err, method, span)
	}

	var merchandiseTotal int
	for _, item := range orderItems {
		if item.Quantity <= 0 || item.Price < 0 {
			status = "error"
			return errorhandler.HandleError[*repository.TransactionResult](s.logger, transaction_errors.ErrFailedOrderItemEmpty, method, span)
		}
		merchandiseTotal += int(item.Price) * int(item.Quantity)
	}

	totalAmount := merchandiseTotal + int(shipping.ShippingCost)

	if req.Amount == 0 {
		req.Amount = int(existingTx.Amount)
	}

	if req.PaymentMethod == "" {
		req.PaymentMethod = existingTx.PaymentMethod
	}

	ppn := totalAmount * 11 / 100
	totalAmountWithTax := totalAmount + ppn

	paymentStatus := PaymentStatusSuccess
	if req.Amount < totalAmountWithTax {
		status = "error"
		return errorhandler.HandleError[*repository.TransactionResult](s.logger, transaction_errors.ErrFailedPaymentInsufficientBalance, method, span)
	}
	if req.PaymentStatus != nil && *req.PaymentStatus != "" && CanTransitionPaymentStatus(existingTx.PaymentStatus, *req.PaymentStatus) && *req.PaymentStatus != PaymentStatusSuccess {
		paymentStatus = *req.PaymentStatus
	}
	if !CanTransitionPaymentStatus(existingTx.PaymentStatus, paymentStatus) {
		status = "error"
		return errorhandler.HandleError[*repository.TransactionResult](s.logger, transaction_errors.ErrFailedPaymentStatusCannotBeModified, method, span)
	}

	req.Amount = totalAmountWithTax
	req.PaymentStatus = &paymentStatus

	transaction, err := s.transactionCommand.Update(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*repository.TransactionResult](s.logger, err, method, span)
	}

	s.cache.DeleteTransactionCache(ctx, *req.TransactionID)

	// Stats pipeline: republish so stats-writer replaces the row via
	// ReplacingMergeTree(event_version) with the latest payment status.
	// Best-effort: a stats publish failure must not fail the update.
	if s.outbox != nil {
		statsEvent, statsErr := json.Marshal(events.StatsEnvelope{
			EventID: uuid.NewString(),
			Payload: mustJSON(events.TransactionEvent{
				TransactionID: transaction.TransactionID,
				OrderID:       transaction.OrderID,
				MerchantID:    transaction.MerchantID,
				PaymentMethod: transaction.PaymentMethod,
				Amount:        transaction.Amount,
				Status:        transaction.PaymentStatus,
				EventTime:     time.Now().UTC().Format(time.RFC3339),
			}),
		})
		if statsErr != nil {
			s.logger.Error("failed to marshal stats event", zap.Error(statsErr), zap.Int32("transaction_id", transaction.TransactionID))
		} else if _, enqueueErr := s.outbox.Create(ctx, "stats.ecommerce.transaction.event", strconv.Itoa(int(transaction.TransactionID)), statsEvent); enqueueErr != nil {
			s.logger.Error("failed to enqueue stats event", zap.Error(enqueueErr), zap.Int32("transaction_id", transaction.TransactionID))
		}
	}

	logSuccess("Successfully updated transaction", zap.Int32("transaction_id", transaction.TransactionID))

	return transaction, nil
}

func (s *transactionCommandService) Trash(ctx context.Context, transactionID int) (*repository.TransactionResult, error) {
	const method = "Trash"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transaction_id", transactionID))

	defer func() {
		end(status)
	}()

	res, err := s.transactionCommand.Trash(ctx, transactionID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*repository.TransactionResult](s.logger, err, method, span)
	}

	s.cache.DeleteTransactionCache(ctx, transactionID)

	logSuccess("Successfully trashed transaction", zap.Int("transaction_id", transactionID))

	return res, nil
}

func (s *transactionCommandService) Restore(ctx context.Context, transactionID int) (*repository.TransactionResult, error) {
	const method = "Restore"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transaction_id", transactionID))

	defer func() {
		end(status)
	}()

	res, err := s.transactionCommand.Restore(ctx, transactionID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*repository.TransactionResult](s.logger, err, method, span)
	}

	s.cache.DeleteTransactionCache(ctx, transactionID)

	logSuccess("Successfully restored transaction", zap.Int("transaction_id", transactionID))

	return res, nil
}

func (s *transactionCommandService) DeletePermanent(ctx context.Context, transactionID int) (bool, error) {
	const method = "DeletePermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transaction_id", transactionID))

	defer func() {
		end(status)
	}()

	success, err := s.transactionCommand.DeletePermanent(ctx, transactionID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](s.logger, err, method, span)
	}

	s.cache.DeleteTransactionCache(ctx, transactionID)

	logSuccess("Successfully permanently deleted transaction", zap.Int("transaction_id", transactionID))

	return success, nil
}

func (s *transactionCommandService) DeleteByOrderIDPermanent(ctx context.Context, orderID int) (bool, error) {
	const method = "DeleteByOrderIDPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order_id", orderID))

	defer func() {
		end(status)
	}()

	success, err := s.transactionCommand.DeleteByOrderIDPermanent(ctx, orderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](s.logger, err, method, span)
	}

	s.cache.InvalidateTransactionCache(ctx)

	logSuccess("Successfully permanently deleted transactions by order", zap.Int("order_id", orderID))

	return success, nil
}

func (s *transactionCommandService) RestoreAll(ctx context.Context) (bool, error) {
	const method = "RestoreAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.transactionCommand.RestoreAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](s.logger, err, method, span)
	}

	s.cache.InvalidateTransactionCache(ctx)
	logSuccess("Successfully restored all transactions")

	return success, nil
}

func (s *transactionCommandService) DeleteAll(ctx context.Context) (bool, error) {
	const method = "DeleteAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.transactionCommand.DeleteAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](s.logger, err, method, span)
	}

	s.cache.InvalidateTransactionCache(ctx)
	logSuccess("Successfully permanently deleted all transactions")

	return success, nil
}

// placeholder to avoid unused import
var _ = models.Transaction{}

func mustJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}
