package service

import (
	mencache "github.com/MamangRust/microservice-ecommerce-auth/cache"
	"github.com/MamangRust/microservice-ecommerce-auth/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/auth"
	"github.com/MamangRust/microservice-ecommerce-pkg/hash"
	"github.com/MamangRust/microservice-ecommerce-pkg/kafka"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-pkg/outbox"
	"github.com/MamangRust/microservice-ecommerce-shared/observability"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Service struct {
	Login         LoginService
	Register      RegistrationService
	PasswordReset PasswordResetService
	Identify      IdentifyService
}

type Deps struct {
	Mencache      *mencache.Mencache
	Repositories  *repository.Repositories
	Token         auth.TokenManager
	Hash          hash.HashPassword
	Logger        logger.LoggerInterface
	Kafka         *kafka.Kafka
	MongoDB        *mongo.Database
	Outbox        *outbox.OutboxService
	Observability observability.TraceLoggerObservability
}

func NewService(deps *Deps) *Service {
	tokenService := NewTokenService(&tokenServiceDeps{
		Token:         deps.Token,
		RefreshToken:  deps.Repositories.RefreshToken,
		Logger:        deps.Logger,
		Observability: deps.Observability,
	})

	return &Service{
		Login:         newLogin(deps, tokenService, deps.Observability, deps.Mencache.LoginCache),
		Register:      newRegister(deps, deps.Observability, deps.Mencache.RegisterCache),
		PasswordReset: newPasswordReset(deps, deps.Observability, deps.Mencache.PasswordResetCache),
		Identify:      newIdentity(deps, tokenService, deps.Observability, deps.Mencache.IdentityCache),
	}
}

func newLogin(deps *Deps, tokenService *tokenService, observability observability.TraceLoggerObservability, cache mencache.LoginCache) LoginService {
	return NewLoginService(&LoginServiceDeps{
		Cache:          cache,
		Logger:         deps.Logger,
		Hash:           deps.Hash,
		UserRepository: deps.Repositories.User,
		RefreshToken:   deps.Repositories.RefreshToken,
		Token:          deps.Token,
		TokenService:   tokenService,
		Observability:  observability,
	})
}

func newRegister(deps *Deps, observability observability.TraceLoggerObservability, cache mencache.RegisterCache) RegistrationService {
	return NewRegisterService(&RegisterServiceDeps{
		Cache:         cache,
		User:          deps.Repositories.User,
		Role:          deps.Repositories.Role,
		UserRole:      deps.Repositories.UserRole,
		Hash:          deps.Hash,
		Kafka:         deps.Kafka,
		MongoDB:        deps.MongoDB,
		Outbox:        deps.Outbox,
		Logger:        deps.Logger,
		Observability: observability,
	})
}

func newPasswordReset(deps *Deps, observability observability.TraceLoggerObservability, cache mencache.PasswordResetCache) PasswordResetService {
	return NewPasswordResetService(&PasswordResetServiceDeps{
		Cache:         cache,
		Kafka:         deps.Kafka,
		MongoDB:        deps.MongoDB,
		Outbox:        deps.Outbox,
		Logger:        deps.Logger,
		User:          deps.Repositories.User,
		ResetToken:    deps.Repositories.ResetToken,
		Observability: observability,
	})
}

func newIdentity(deps *Deps, tokenService *tokenService, observability observability.TraceLoggerObservability, cache mencache.IdentityCache) IdentifyService {
	return NewIdentityService(&IdentityServiceDeps{
		Cache:         cache,
		Token:         deps.Token,
		RefreshToken:  deps.Repositories.RefreshToken,
		User:          deps.Repositories.User,
		Logger:        deps.Logger,
		TokenService:  tokenService,
		Observability: observability,
	})
}
