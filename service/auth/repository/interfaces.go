package repository

import (
	"context"

	useradapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/user"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByEmailAndVerify(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, user_id int) (*models.User, error)
	CreateUser(ctx context.Context, request *requests.RegisterRequest) (*models.User, error)
	UpdateUserIsVerified(ctx context.Context, user_id int, is_verified bool) (*models.User, error)
	UpdateUserPassword(ctx context.Context, user_id int, password string) (*models.User, error)
	FindByVerificationCode(ctx context.Context, verification_code string) (*models.User, error)
}

type ResetTokenRepository interface {
	FindByToken(ctx context.Context, code string) (*models.ResetToken, error)
	CreateResetToken(ctx context.Context, req *requests.CreateResetTokenRequest) (*models.ResetToken, error)
	CreateResetTokenInTx(ctx context.Context, tx *mongo.Database, req *requests.CreateResetTokenRequest) (*models.ResetToken, error)
	DeleteResetToken(ctx context.Context, user_id int) error
}

type RefreshTokenRepository interface {
	FindByToken(ctx context.Context, token string) (*models.RefreshToken, error)
	FindByUserId(ctx context.Context, user_id int) (*models.RefreshToken, error)
	CreateRefreshToken(ctx context.Context, req *requests.CreateRefreshToken) (*models.RefreshToken, error)
	UpdateRefreshToken(ctx context.Context, req *requests.UpdateRefreshToken) (*models.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteRefreshTokenByUserId(ctx context.Context, user_id int) error
}

// Repositories is declared in repository.go; this file only holds the narrow
// interface bounds the auth service depends on.

// Compile-time guard so the auth package cannot drift from the adapter
// surface it claims to depend on.
var (
	_ useradapter.QueryRepository   = (UserRepository)(nil)
	_ useradapter.CommandRepository = (UserRepository)(nil)
)
