// Package role adapts the role service query gRPC API into the shared domain
// model. Role assignment lives in the user_role service, so this adapter only
// exposes the read path.
package role

import (
	"context"
	"time"

	pbrole "github.com/MamangRust/microservice-ecommerce-grpc-pb/role"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-pkg/resilience"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/role_errors"
)

// QueryRepository is the contract consumers depend on for role reads.
type QueryRepository interface {
	FindById(ctx context.Context, roleID int) (*models.Role, error)
	FindByName(ctx context.Context, name string) (*models.Role, error)
}

// Repository implements QueryRepository on top of the generated role query
// client, guarding every call.
type Repository struct {
	client pbrole.RoleQueryServiceClient
	guard  *resilience.DependencyGuard
}

// New builds the role adapter from its query client.
func New(client pbrole.RoleQueryServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{client: client}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// NewAdapter builds the role adapter from its query and command clients. Only
// the query client is used; role assignment is served by the user_role adapter.
func NewAdapter(query pbrole.RoleQueryServiceClient, _ pbrole.RoleCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	return New(query, opts...)
}

// NewQueryAdapter wraps a generated role query client into a QueryRepository.
func NewQueryAdapter(client pbrole.RoleQueryServiceClient, opts ...adapter.GuardOption) QueryRepository {
	return New(client, opts...)
}

// SetGuard attaches the dependency guard used to wrap every outbound call.
func (r *Repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func (r *Repository) FindById(ctx context.Context, roleID int) (*models.Role, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbrole.ApiResponseRole, error) {
		return r.client.FindByIdRole(callCtx, &pbrole.FindByIdRoleRequest{RoleId: int32(roleID)})
	})
	if err != nil {
		return nil, role_errors.ErrRoleNotFound.WithInternal(err)
	}
	if res.GetData() == nil {
		return nil, role_errors.ErrRoleNotFound
	}
	return roleToModel(res.Data), nil
}

func (r *Repository) FindByName(ctx context.Context, name string) (*models.Role, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbrole.ApiResponseRole, error) {
		return r.client.FindByNameRole(callCtx, &pbrole.FindByNameRoleRequest{Name: name})
	})
	if err != nil {
		return nil, role_errors.ErrRoleNotFound.WithInternal(err)
	}
	if res.GetData() == nil {
		return nil, role_errors.ErrRoleNotFound
	}
	return roleToModel(res.Data), nil
}

func roleToModel(role *pbrole.RoleResponse) *models.Role {
	if role == nil {
		return nil
	}
	return &models.Role{
		RoleID:    role.Id,
		RoleName:  role.Name,
		CreatedAt: parseTimeString(role.CreatedAt),
		UpdatedAt: parseTimeString(role.UpdatedAt),
	}
}

func parseTimeString(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02 15:04:05.000", s)
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05Z", s)
		if err != nil {
			return nil
		}
	}
	return &t
}

var _ QueryRepository = (*Repository)(nil)
