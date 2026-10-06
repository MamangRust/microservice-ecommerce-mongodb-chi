package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/role_errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type userRoleRepository struct {
	db *mongo.Database
}

func NewUserRoleRepository(db *mongo.Database) UserRoleRepository {
	return &userRoleRepository{db: db}
}

func (r *userRoleRepository) AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*models.UserRole, error) {
	id, err := roleGetNextID(r.db, "user_roles", ctx)
	if err != nil {
		return nil, role_errors.ErrAssignRole.WithInternal(err)
	}
	now := time.Now()
	userRole := &models.UserRole{
		UserRoleID: id,
		UserID:     int32(req.UserId),
		RoleID:     int32(req.RoleId),
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}
	_, err = r.db.Collection("user_roles").InsertOne(ctx, userRole)
	if err != nil {
		return nil, role_errors.ErrAssignRole.WithInternal(err)
	}
	return userRole, nil
}

func (r *userRoleRepository) RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error {
	result, err := r.db.Collection("user_roles").DeleteOne(ctx, bson.M{
		"user_id": req.UserId,
		"role_id": req.RoleId,
	})
	if err != nil {
		return role_errors.ErrRemoveRole.WithInternal(err)
	}
	if result.DeletedCount == 0 {
		return role_errors.ErrRemoveRole
	}
	return nil
}
