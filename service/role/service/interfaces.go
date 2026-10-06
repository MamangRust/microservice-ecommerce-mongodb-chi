package service

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-role/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type RoleQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllRole) ([]*repository.RoleResult, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllRole) ([]*repository.RoleResult, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllRole) ([]*repository.RoleResult, *int, error)
	FindByID(ctx context.Context, role_id int) (*models.Role, error)
	FindByName(ctx context.Context, name string) (*models.Role, error)
	FindByUserId(ctx context.Context, id int) ([]*models.Role, error)
}

type RoleCommandService interface {
	Create(ctx context.Context, request *requests.CreateRoleRequest) (*models.Role, error)
	Update(ctx context.Context, request *requests.UpdateRoleRequest) (*models.Role, error)
	Trash(ctx context.Context, role_id int) (*models.Role, error)
	Restore(ctx context.Context, role_id int) (*models.Role, error)
	DeletePermanent(ctx context.Context, role_id int) (bool, error)

	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)

	AssignRoleToUser(ctx context.Context, request *requests.CreateUserRoleRequest) (*models.UserRole, error)
	RemoveRoleFromUser(ctx context.Context, request *requests.RemoveUserRoleRequest) error
}
