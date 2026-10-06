package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/MamangRust/microservice-ecommerce-grpc-merchant/cache"
	"github.com/MamangRust/microservice-ecommerce-grpc-merchant/repository"
	useradapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/user"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-pkg/email"
	"github.com/MamangRust/microservice-ecommerce-pkg/event"
	"github.com/MamangRust/microservice-ecommerce-pkg/kafka"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-pkg/outbox"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errorhandler"
	"github.com/MamangRust/microservice-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type merchantCommandService struct {
	kafka              *kafka.Kafka
	observability      observability.TraceLoggerObservability
	cache              cache.MerchantCommandCache
	merchantRepository repository.MerchantCommandRepository
	merchantQuery      repository.MerchantQueryRepository
	userRepository     useradapter.QueryRepository
	mongoDB             *mongo.Database
	outbox             *outbox.OutboxService
	logger             logger.LoggerInterface
}

type MerchantCommandServiceDeps struct {
	Kafka              *kafka.Kafka
	Observability      observability.TraceLoggerObservability
	Cache              cache.MerchantCommandCache
	MerchantRepository repository.MerchantCommandRepository
	MerchantQuery      repository.MerchantQueryRepository
	UserRepository     useradapter.QueryRepository
	MongoDB             *mongo.Database
	Outbox             *outbox.OutboxService
	Logger             logger.LoggerInterface
}

func NewMerchantCommandService(deps *MerchantCommandServiceDeps) MerchantCommandService {
	return &merchantCommandService{
		kafka:              deps.Kafka,
		observability:      deps.Observability,
		cache:              deps.Cache,
		merchantRepository: deps.MerchantRepository,
		merchantQuery:      deps.MerchantQuery,
		userRepository:     deps.UserRepository,
		mongoDB:             deps.MongoDB,
		outbox:             deps.Outbox,
		logger:             deps.Logger,
	}
}

func (s *merchantCommandService) Create(ctx context.Context, request *requests.CreateMerchantRequest) (*models.Merchant, error) {
	const method = "Create"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("user.id", request.UserID))

	defer func() {
		end(status)
	}()

	user, err := s.userRepository.FindByID(ctx, request.UserID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*models.Merchant](
			s.logger,
			err,
			method,
			span,
			zap.Int("user.id", request.UserID),
		)
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   "Welcome to SanEdge Merchant Portal",
		"Message": "Your merchant account has been created successfully.",
		"Button":  "Upload Documents",
		"Link":    fmt.Sprintf("https://sanedge.example.com/merchant/%d/documents", user.UserID),
	})

	payloadBytes, err := event.MarshalEmail("merchant.created", user.Email, "Initial Verification - SanEdge", htmlBody)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*models.Merchant](
			s.logger,
			err,
			method,
			span,
			zap.Int("user.id", request.UserID),
		)
	}

	var res *models.Merchant
	{
		res, err = s.merchantRepository.Create(ctx, request)
		if err != nil {
			status = "error"
			return errorhandler.HandleError[*models.Merchant](
				s.logger,
				err,
				method,
				span,
				zap.Int("user.id", request.UserID),
			)
		}
		if s.outbox != nil {
			if enqErr := s.outbox.EnqueueInTx(ctx, NewOutboxQuerier(s.mongoDB), "email-service-topic-merchant-create", strconv.Itoa(int(res.MerchantID)), payloadBytes); enqErr != nil {
				s.logger.Error("failed to enqueue merchant create email", zap.Error(enqErr))
			}
		} else if s.kafka != nil {
			if sendErr := s.kafka.SendMessage("email-service-topic-merchant-create", strconv.Itoa(int(res.MerchantID)), payloadBytes); sendErr != nil {
				s.logger.Error("Failed to send email to Kafka", zap.Error(sendErr))
			}
		}
	}

	logSuccess("Successfully created merchant", zap.Int("merchant.id", int(res.MerchantID)))

	return res, nil
}

func (s *merchantCommandService) Update(ctx context.Context, request *requests.UpdateMerchantRequest) (*models.Merchant, error) {
	const method = "Update"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchant.id", *request.MerchantID))

	defer func() {
		end(status)
	}()

	res, err := s.merchantRepository.Update(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*models.Merchant](
			s.logger,
			err,
			method,
			span,
			zap.Int("merchant.id", *request.MerchantID),
		)
	}

	s.cache.DeleteMerchantCache(ctx, *request.MerchantID)

	logSuccess("Successfully updated merchant", zap.Int("merchant.id", *request.MerchantID))

	return res, nil
}

func (s *merchantCommandService) UpdateMerchantStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*models.Merchant, error) {
	const method = "UpdateMerchantStatus"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchant.id", *request.MerchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantQuery.FindByID(ctx, *request.MerchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*models.Merchant](
			s.logger,
			err,
			method,
			span,
			zap.Int("merchant.id", *request.MerchantID),
		)
	}

	user, err := s.userRepository.FindByID(ctx, int(merchant.UserID))
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*models.Merchant](
			s.logger,
			err,
			method,
			span,
			zap.Int("user.id", int(merchant.UserID)),
		)
	}

	statusReq := request.Status
	subject := ""
	message := ""
	buttonLabel := "Go to Portal"
	link := fmt.Sprintf("https://sanedge.example.com/merchant/%d/dashboard", *request.MerchantID)

	switch statusReq {
	case "active":
		subject = "Your Merchant Account is Now Active"
		message = "Congratulations! Your merchant account is now active."
	case "inactive":
		subject = "Merchant Account Set to Inactive"
		message = "Your merchant account status has been set to inactive."
	case "rejected":
		subject = "Merchant Account Rejected"
		message = "We're sorry, your merchant account has been rejected."
	}

	var res *models.Merchant
	{
		res, err = s.merchantRepository.UpdateStatus(ctx, request)
		if err != nil {
			status = "error"
			return errorhandler.HandleError[*models.Merchant](
				s.logger,
				err,
				method,
				span,
				zap.Int("merchant.id", *request.MerchantID),
			)
		}
		if subject != "" {
			htmlBody := email.GenerateEmailHTML(map[string]string{
				"Title":   subject,
				"Message": message,
				"Button":  buttonLabel,
				"Link":    link,
			})
			payloadBytes, marshalErr := event.MarshalEmail("merchant.status_updated", user.Email, subject, htmlBody)
			if marshalErr != nil {
				s.logger.Error("failed to marshal merchant status email", zap.Error(marshalErr), zap.Int32("merchant_id", res.MerchantID))
			} else if s.outbox != nil {
				if enqErr := s.outbox.EnqueueInTx(ctx, NewOutboxQuerier(s.mongoDB), "email-service-topic-merchant-update-status", strconv.Itoa(int(res.MerchantID)), payloadBytes); enqErr != nil {
					s.logger.Error("failed to enqueue merchant status email", zap.Error(enqErr))
				}
			} else if s.kafka != nil {
				if sendErr := s.kafka.SendMessage("email-service-topic-merchant-update-status", strconv.Itoa(int(res.MerchantID)), payloadBytes); sendErr != nil {
					s.logger.Error("failed to publish merchant status email", zap.Error(sendErr), zap.Int32("merchant_id", res.MerchantID))
				}
			}
		}
	}

	if s.kafka != nil && subject != "" {
		if statusEvent, marshalErr := json.Marshal(map[string]any{
			"merchantId": res.MerchantID,
			"status":     request.Status,
			"timestamp":  time.Now().UnixMilli(),
		}); marshalErr != nil {
			s.logger.Error("failed to marshal merchant status event", zap.Error(marshalErr), zap.Int32("merchant_id", res.MerchantID))
		} else if sendErr := s.kafka.SendMessage("transaction-service-topic-merchant-status-event", strconv.Itoa(int(res.MerchantID)), statusEvent); sendErr != nil {
			s.logger.Error("failed to publish merchant status event", zap.Error(sendErr), zap.Int32("merchant_id", res.MerchantID))
		}
	}

	s.cache.DeleteMerchantCache(ctx, *request.MerchantID)

	logSuccess("Successfully updated merchant status", zap.Int("merchant.id", *request.MerchantID))

	return res, nil
}

func (s *merchantCommandService) Trash(ctx context.Context, merchantID int) (*models.Merchant, error) {
	const method = "Trash"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchant.id", merchantID))

	defer func() {
		end(status)
	}()

	res, err := s.merchantRepository.Trash(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*models.Merchant](
			s.logger,
			err,
			method,
			span,
			zap.Int("merchant.id", merchantID),
		)
	}

	s.cache.DeleteMerchantCache(ctx, merchantID)

	logSuccess("Successfully trashed merchant", zap.Int("merchant.id", merchantID))

	return res, nil
}

func (s *merchantCommandService) Restore(ctx context.Context, merchantID int) (*models.Merchant, error) {
	const method = "Restore"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchant.id", merchantID))

	defer func() {
		end(status)
	}()

	res, err := s.merchantRepository.Restore(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*models.Merchant](
			s.logger,
			err,
			method,
			span,
			zap.Int("merchant.id", merchantID),
		)
	}

	s.cache.DeleteMerchantCache(ctx, merchantID)

	logSuccess("Successfully restored merchant", zap.Int("merchant.id", merchantID))

	return res, nil
}

func (s *merchantCommandService) DeletePermanent(ctx context.Context, merchantID int) (bool, error) {
	const method = "DeletePermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchant.id", merchantID))

	defer func() {
		end(status)
	}()

	res, err := s.merchantRepository.DeletePermanent(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			err,
			method,
			span,
			zap.Int("merchant.id", merchantID),
		)
	}

	s.cache.DeleteMerchantCache(ctx, merchantID)

	logSuccess("Successfully permanently deleted merchant", zap.Int("merchant.id", merchantID))

	return res, nil
}

func (s *merchantCommandService) RestoreAll(ctx context.Context) (bool, error) {
	const method = "RestoreAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.merchantRepository.RestoreAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			err,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all merchants")

	return res, nil
}

func (s *merchantCommandService) DeleteAll(ctx context.Context) (bool, error) {
	const method = "DeleteAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.merchantRepository.DeleteAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			err,
			method,
			span,
		)
	}

	logSuccess("Successfully permanently deleted all merchants")

	return res, nil
}
