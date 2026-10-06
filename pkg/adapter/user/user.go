// Package user adapts the user service's query/command gRPC APIs into the
// shared domain model.
package user

import (
	"context"
	"time"

	pbuser "github.com/MamangRust/microservice-ecommerce-grpc-pb/user"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-pkg/resilience"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	sharedErrors "github.com/MamangRust/microservice-ecommerce-shared/errors"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/user_errors"
)

// QueryRepository is the contract consumers depend on for user reads.
type QueryRepository interface {
	FindByID(ctx context.Context, id int) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByEmailAndVerify(ctx context.Context, email string) (*models.User, error)
	FindByVerificationCode(ctx context.Context, code string) (*models.User, error)
}

// CommandRepository is the contract consumers depend on for user writes.
type CommandRepository interface {
	CreateUser(ctx context.Context, request *requests.RegisterRequest) (*models.User, error)
	UpdateUserIsVerified(ctx context.Context, userID int, isVerified bool) (*models.User, error)
	UpdateUserPassword(ctx context.Context, userID int, password string) (*models.User, error)
}

// Repository implements QueryRepository and CommandRepository on top of the
// generated user service clients, guarding every call.
type Repository struct {
	query   pbuser.UserQueryServiceClient
	command pbuser.UserCommandServiceClient
	guard   *resilience.DependencyGuard
}

// New builds a user adapter from its query and command clients.
func New(query pbuser.UserQueryServiceClient, command pbuser.UserCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{query: query, command: command}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// NewQueryAdapter builds a read-only user adapter.
func NewQueryAdapter(query pbuser.UserQueryServiceClient, opts ...adapter.GuardOption) *Repository {
	return New(query, nil, opts...)
}

// NewCommandAdapter builds a write-only user adapter.
func NewCommandAdapter(command pbuser.UserCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	return New(nil, command, opts...)
}

// SetGuard attaches the dependency guard used to wrap every outbound call.
func (r *Repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func (r *Repository) FindByID(ctx context.Context, id int) (*models.User, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbuser.ApiResponseUser, error) {
		return r.query.FindById(callCtx, &pbuser.FindByIdUserRequest{Id: int32(id)})
	})
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	if res.GetData() == nil {
		return nil, user_errors.ErrUserNotFound
	}
	return userToModel(res.Data), nil
}

// FindByEmail returns the user with the password hash so the auth service can
// verify credentials.
func (r *Repository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbuser.ApiResponseUserWithPassword, error) {
		return r.query.FindByEmail(callCtx, &pbuser.FindByEmailRequest{Email: email})
	})
	if err != nil {
		return nil, user_errors.ErrUserNotFound.WithInternal(err)
	}
	if res.GetData() == nil {
		return nil, user_errors.ErrUserNotFound
	}
	return userWithPasswordToModel(res.Data), nil
}

// FindByEmailAndVerify mirrors FindByEmail. The user service does not expose a
// verified-only lookup yet, so callers get the same user and must not assume
// verification has been enforced.
func (r *Repository) FindByEmailAndVerify(ctx context.Context, email string) (*models.User, error) {
	return r.FindByEmail(ctx, email)
}

func (r *Repository) FindByVerificationCode(ctx context.Context, code string) (*models.User, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbuser.ApiResponseUser, error) {
		return r.query.FindByVerificationCode(callCtx, &pbuser.FindByVerificationCodeRequest{VerificationCode: code})
	})
	if err != nil {
		return nil, user_errors.ErrUserNotFound.WithInternal(err)
	}
	if res.GetData() == nil {
		return nil, user_errors.ErrUserNotFound
	}
	return userToModel(res.Data), nil
}

// CreateUser registers a user through the user service. The password must
// already be hashed by the caller.
func (r *Repository) CreateUser(ctx context.Context, request *requests.RegisterRequest) (*models.User, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbuser.ApiResponseUser, error) {
		return r.command.Create(callCtx, &pbuser.CreateUserRequest{
			Firstname:       request.FirstName,
			Lastname:        request.LastName,
			Email:           request.Email,
			Password:        request.Password,
			ConfirmPassword: request.ConfirmPassword,
		})
	})
	if err != nil {
		return nil, user_errors.ErrCreateUser.WithInternal(err)
	}
	if res.GetData() == nil {
		return nil, user_errors.ErrCreateUser
	}
	return userToModel(res.Data), nil
}

func (r *Repository) UpdateUserIsVerified(ctx context.Context, userID int, isVerified bool) (*models.User, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbuser.ApiResponseUser, error) {
		return r.command.UpdateIsVerified(callCtx, &pbuser.UpdateUserIsVerifiedRequest{
			Id:         int32(userID),
			IsVerified: isVerified,
		})
	})
	if err != nil {
		return nil, user_errors.ErrUpdateUserVerificationCode.WithInternal(err)
	}
	if res.GetData() == nil {
		return nil, user_errors.ErrUpdateUserVerificationCode
	}
	return userToModel(res.Data), nil
}

func (r *Repository) UpdateUserPassword(ctx context.Context, userID int, password string) (*models.User, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbuser.ApiResponseUser, error) {
		return r.command.UpdatePassword(callCtx, &pbuser.UpdateUserPasswordRequest{
			Id:       int32(userID),
			Password: password,
		})
	})
	if err != nil {
		return nil, user_errors.ErrUpdateUserPassword.WithInternal(err)
	}
	if res.GetData() == nil {
		return nil, user_errors.ErrUpdateUserPassword
	}
	return userToModel(res.Data), nil
}

func userToModel(u *pbuser.UserResponse) *models.User {
	if u == nil {
		return nil
	}
	return &models.User{
		UserID:    u.Id,
		Firstname: u.Firstname,
		Lastname:  u.Lastname,
		Email:     u.Email,
		CreatedAt: parseTimeString(u.CreatedAt),
		UpdatedAt: parseTimeString(u.UpdatedAt),
	}
}

func userWithPasswordToModel(u *pbuser.UserResponseWithPassword) *models.User {
	if u == nil {
		return nil
	}
	return &models.User{
		UserID:    u.Id,
		Firstname: u.Firstname,
		Lastname:  u.Lastname,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: parseTimeString(u.CreatedAt),
		UpdatedAt: parseTimeString(u.UpdatedAt),
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

var (
	_ QueryRepository   = (*Repository)(nil)
	_ CommandRepository = (*Repository)(nil)
)
