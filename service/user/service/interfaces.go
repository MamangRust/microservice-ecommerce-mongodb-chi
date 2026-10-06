package service

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-user/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

// UserQueryService handles query operations related to user data.
type UserQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllUsers) ([]*repository.UserResult, *int, error)
	FindByID(ctx context.Context, id int) (*models.User, error)
	FindByEmailWithPassword(ctx context.Context, email string) (*models.User, error)
	FindByVerificationCode(ctx context.Context, code string) (*models.User, error)
	FindActive(ctx context.Context, req *requests.FindAllUsers) ([]*repository.UserResult, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*repository.UserResult, *int, error)
}

// UserCommandService handles command operations related to user management.
type UserCommandService interface {
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
