package service

import (
	"context"
	"strconv"
	"time"

	mencache "github.com/MamangRust/microservice-ecommerce-auth/cache"
	"github.com/MamangRust/microservice-ecommerce-auth/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/email"
	"github.com/MamangRust/microservice-ecommerce-pkg/event"
	"github.com/MamangRust/microservice-ecommerce-pkg/hash"
	"github.com/MamangRust/microservice-ecommerce-pkg/kafka"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-pkg/outbox"
	randomstring "github.com/MamangRust/microservice-ecommerce-pkg/randomstring"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/microservice-ecommerce-shared/errorhandler"
	roleadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/role"
	userroleadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/user_role"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	user_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/user_errors"
	"github.com/MamangRust/microservice-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RegisterServiceDeps struct {
	Cache         mencache.RegisterCache
	User          repository.UserRepository
	Role          roleadapter.QueryRepository
	UserRole      userroleadapter.CommandRepository
	Hash          hash.HashPassword
	Kafka         *kafka.Kafka
	MongoDB        *mongo.Database
	Outbox        *outbox.OutboxService
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

type registerService struct {
	mencache      mencache.RegisterCache
	user          repository.UserRepository
	role          roleadapter.QueryRepository
	userRole      userroleadapter.CommandRepository
	hash          hash.HashPassword
	kafka         *kafka.Kafka
	mongoDB        *mongo.Database
	outbox        *outbox.OutboxService
	logger        logger.LoggerInterface
	observability observability.TraceLoggerObservability
}

func NewRegisterService(params *RegisterServiceDeps) *registerService {
	return &registerService{
		mencache:      params.Cache,
		user:          params.User,
		role:          params.Role,
		userRole:      params.UserRole,
		hash:          params.Hash,
		kafka:         params.Kafka,
		mongoDB:        params.MongoDB,
		outbox:        params.Outbox,
		logger:        params.Logger,
		observability: params.Observability,
	}
}

func (s *registerService) Register(ctx context.Context, request *requests.RegisterRequest) (*models.User, error) {
	const method = "Register"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("email", request.Email))

	defer func() {
		end(status)
	}()

	existingUser, err := s.user.FindByEmail(ctx, request.Email)
	if err == nil && existingUser != nil {
		status = "error"
		return sharederrorhandler.HandleError[*models.User](
			s.logger,
			user_errors.ErrUserEmailAlready,
			method,
			span,
			zap.String("email", request.Email),
		)
	}

	const defaultRoleName = "ROLE_ADMIN"
	role, err := s.role.FindByName(ctx, defaultRoleName)
	if err != nil || role == nil {
		status = "error"
		return sharederrorhandler.HandleError[*models.User](s.logger, err, method, span, zap.String("role_name", defaultRoleName))
	}

	random, err := randomstring.GenerateRandomString(10)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*models.User](s.logger, err, method, span)
	}
	request.VerifiedCode = random
	request.IsVerified = false

	newUser, err := s.user.CreateUser(ctx, request)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*models.User](s.logger, err, method, span)
	}

	_, err = s.userRole.AssignRoleToUser(ctx, &requests.CreateUserRoleRequest{
		UserId: int(newUser.UserID),
		RoleId: int(role.RoleID),
	})
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*models.User](s.logger, err, method, span, zap.Int("user.id", int(newUser.UserID)))
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   "Welcome to SanEdge",
		"Message": "Your account has been successfully created.",
		"Button":  "Verify Now",
		"Link":    "https://sanedge.example.com/login?verify_code=" + request.VerifiedCode,
	})

	payloadBytes, err := event.MarshalEmail("auth.register", request.Email, "Welcome to SanEdge", htmlBody)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*models.User](s.logger, err, method, span)
	}

	if s.outbox != nil {
		if enqueueErr := s.outbox.Enqueue(ctx, "email-service-topic-auth-register", strconv.Itoa(int(newUser.UserID)), payloadBytes); enqueueErr != nil {
			s.logger.Error("failed to enqueue registration email to outbox", zap.Error(enqueueErr), zap.String("email", request.Email))
		}
	} else if s.kafka != nil {
		go func() {
			if sendErr := s.kafka.SendMessage("email-service-topic-auth-register", strconv.Itoa(int(newUser.UserID)), payloadBytes); sendErr != nil {
				s.logger.Error("failed to send registration email via kafka", zap.Error(sendErr), zap.String("email", request.Email))
			}
		}()
	}

	s.mencache.SetVerificationCodeCache(ctx, request.Email, random, 15*time.Minute)

	logSuccess("User registered successfully",
		zap.String("email", request.Email),
		zap.String("first_name", request.FirstName),
		zap.String("last_name", request.LastName),
	)

	return newUser, nil
}
