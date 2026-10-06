package repository

import (
	"context"

	useradapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/user"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

// userRepository is a thin façade over useradapter.Repository. It exists to
// keep the auth-side UserRepository interface stable while the user gRPC
// client moves out of auth's own package. The dial site constructs the adapter
// with the generated pb clients; auth only sees the narrow interface bounds.
type userRepository struct {
	adapter *useradapter.Repository
}

func NewUserRepository(adapter *useradapter.Repository) *userRepository {
	return &userRepository{adapter: adapter}
}

func (r *userRepository) FindByID(ctx context.Context, user_id int) (*models.User, error) {
	return r.adapter.FindByID(ctx, user_id)
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	return r.adapter.FindByEmail(ctx, email)
}

func (r *userRepository) FindByEmailAndVerify(ctx context.Context, email string) (*models.User, error) {
	return r.adapter.FindByEmailAndVerify(ctx, email)
}

func (r *userRepository) FindByVerificationCode(ctx context.Context, code string) (*models.User, error) {
	return r.adapter.FindByVerificationCode(ctx, code)
}

func (r *userRepository) CreateUser(ctx context.Context, req *requests.RegisterRequest) (*models.User, error) {
	return r.adapter.CreateUser(ctx, req)
}

func (r *userRepository) UpdateUserIsVerified(ctx context.Context, user_id int, is_verified bool) (*models.User, error) {
	return r.adapter.UpdateUserIsVerified(ctx, user_id, is_verified)
}

func (r *userRepository) UpdateUserPassword(ctx context.Context, user_id int, password string) (*models.User, error) {
	return r.adapter.UpdateUserPassword(ctx, user_id, password)
}
