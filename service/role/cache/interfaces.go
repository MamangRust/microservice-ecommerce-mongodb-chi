package cache

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-role/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type RoleQueryCache interface {
	SetCachedRoles(ctx context.Context, req *requests.FindAllRole, data []*repository.RoleResult, total *int)
	GetCachedRoles(ctx context.Context, req *requests.FindAllRole) ([]*repository.RoleResult, *int, bool)

	GetCachedRoleById(ctx context.Context, id int) (*models.Role, bool)
	SetCachedRoleById(ctx context.Context, id int, data *models.Role)

	GetCachedRoleByName(ctx context.Context, name string) (*models.Role, bool)
	SetCachedRoleByName(ctx context.Context, name string, data *models.Role)

	GetCachedRoleByUserId(ctx context.Context, userId int) ([]*models.Role, bool)
	SetCachedRoleByUserId(ctx context.Context, userId int, data []*models.Role)

	GetCachedRoleActive(ctx context.Context, req *requests.FindAllRole) ([]*repository.RoleResult, *int, bool)
	SetCachedRoleActive(ctx context.Context, req *requests.FindAllRole, data []*repository.RoleResult, total *int)

	GetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRole) ([]*repository.RoleResult, *int, bool)
	SetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRole, data []*repository.RoleResult, total *int)
}

type RoleCommandCache interface {
	DeleteCachedRole(ctx context.Context, id int)
}
