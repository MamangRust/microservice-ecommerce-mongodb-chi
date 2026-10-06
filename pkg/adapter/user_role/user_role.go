// Package user_role adapts the UserRole service gRPC API into the shared domain
// model. It owns the only place that talks to pb/user_role's UserRoleService.
package user_role

import (
	"context"

	pbrole "github.com/MamangRust/microservice-ecommerce-grpc-pb/role"
	pbuserrole "github.com/MamangRust/microservice-ecommerce-grpc-pb/user_role"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-pkg/resilience"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	userrole_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/user_role_errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

// QueryRepository is the read path consumers use to resolve a user's roles.
type QueryRepository interface {
	FindByUserId(ctx context.Context, userID int) ([]*models.Role, error)
}

// CommandRepository is the write path consumers use to (un)assign roles.
type CommandRepository interface {
	AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*models.UserRole, error)
	RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error
}

// Repository implements QueryRepository and CommandRepository on top of the
// generated user-role service client, guarding every call.
type Repository struct {
	client pbuserrole.UserRoleServiceClient
	guard  *resilience.DependencyGuard
}

// New builds a user-role adapter from the generated client.
func New(client pbuserrole.UserRoleServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{client: client}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// NewAdapter builds a user-role adapter from the generated client.
func NewAdapter(client pbuserrole.UserRoleServiceClient, opts ...adapter.GuardOption) *Repository {
	return New(client, opts...)
}

// NewCommandAdapter returns the adapter restricted to the CommandRepository
// surface for consumers that only (un)assign roles.
func NewCommandAdapter(client pbuserrole.UserRoleServiceClient, opts ...adapter.GuardOption) CommandRepository {
	return New(client, opts...)
}

// SetGuard attaches the dependency guard used to wrap every outbound call.
func (r *Repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

// FindByUserId implements QueryRepository.
func (r *Repository) FindByUserId(ctx context.Context, userID int) ([]*models.Role, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbrole.ApiResponsesRole, error) {
		return r.client.FindByUserId(callCtx, &pbuserrole.FindByIdUserRoleRequest{UserId: int32(userID)})
	})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}

	roles := make([]*models.Role, 0, len(res.Data))
	for _, item := range res.Data {
		if item == nil {
			continue
		}
		roles = append(roles, &models.Role{
			RoleID:   item.Id,
			RoleName: item.Name,
		})
	}
	return roles, nil
}

// AssignRoleToUser implements CommandRepository.
func (r *Repository) AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*models.UserRole, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbuserrole.ApiResponseUserRole, error) {
		return r.client.AssignRoleToUser(callCtx, &pbuserrole.AssignRoleToUserRequest{
			UserId: int32(req.UserId),
			RoleId: int32(req.RoleId),
		})
	})
	if err != nil {
		return nil, userrole_errors.ErrAssignRoleToUser.WithInternal(err)
	}
	if res.GetData() == nil {
		return &models.UserRole{UserID: int32(req.UserId), RoleID: int32(req.RoleId)}, nil
	}
	return &models.UserRole{
		UserRoleID: res.Data.UserRoleId,
		UserID:     res.Data.UserId,
		RoleID:     res.Data.RoleId,
	}, nil
}

// RemoveRoleFromUser implements CommandRepository.
func (r *Repository) RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error {
	_, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*emptypb.Empty, error) {
		return r.client.RemoveRoleFromUser(callCtx, &pbuserrole.RemoveRoleFromUserRequest{
			UserId: int32(req.UserId),
			RoleId: int32(req.RoleId),
		})
	})
	if err != nil {
		return userrole_errors.ErrRemoveRole.WithInternal(err)
	}
	return nil
}

var (
	_ QueryRepository   = (*Repository)(nil)
	_ CommandRepository = (*Repository)(nil)
)
