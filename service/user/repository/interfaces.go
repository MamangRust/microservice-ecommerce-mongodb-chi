package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

// UserResult is used for paginated list queries
type UserResult struct {
	UserID     int32
	Firstname  string
	Lastname   string
	Email      string
	Password   string
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
	DeletedAt  *time.Time
	TotalCount int64
}

// UserResultTrashed includes deleted_at in query results
type UserResultTrashed = UserResult

type UserQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllUsers) ([]*UserResult, error)
	FindActive(ctx context.Context, req *requests.FindAllUsers) ([]*UserResult, error)
	FindTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*UserResult, error)
	FindByID(ctx context.Context, user_id int) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByIDWithPassword(ctx context.Context, user_id int) (*models.User, error)
	FindByEmailWithPassword(ctx context.Context, email string) (*models.User, error)
	FindByVerificationCode(ctx context.Context, code string) (*models.User, error)
}

type UserCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateUserRequest) (*models.User, error)
	Update(ctx context.Context, request *requests.UpdateUserRequest) (*models.User, error)
	UpdateIsVerified(ctx context.Context, user_id int, is_verified bool) (*models.User, error)
	UpdatePassword(ctx context.Context, user_id int, password string) (*models.User, error)
	Trash(ctx context.Context, user_id int) (*models.User, error)
	Restore(ctx context.Context, user_id int) (*models.User, error)
	DeletePermanent(ctx context.Context, user_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
