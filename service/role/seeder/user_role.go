package seeder

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

type userRoleSeeder struct {
	userDB *mongo.Database
	roleDB *mongo.Database
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewUserRoleSeeder(userDB, roleDB *mongo.Database, ctx context.Context, logger logger.LoggerInterface) *userRoleSeeder {
	return &userRoleSeeder{userDB: userDB, roleDB: roleDB, ctx: ctx, logger: logger}
}

func (r *userRoleSeeder) Seed() error {
	userCount, err := r.userDB.Collection("users").CountDocuments(r.ctx, bson.M{"deleted_at": bson.M{"$eq": nil}})
	if err != nil {
		return err
	}
	roleCount, err := r.roleDB.Collection("roles").CountDocuments(r.ctx, bson.M{})
	if err != nil {
		return err
	}

	if userCount == 0 || roleCount == 0 {
		r.logger.Debug("no users or roles found, skipping user_role seeding")
		return nil
	}

	// Get user IDs
	userCursor, err := r.userDB.Collection("users").Find(r.ctx, bson.M{"deleted_at": bson.M{"$eq": nil}})
	if err != nil {
		return err
	}
	defer userCursor.Close(r.ctx)
	var users []bson.M
	if err := userCursor.All(r.ctx, &users); err != nil {
		return err
	}
	if len(users) > 20 {
		users = users[:20]
	}

	// Get role IDs
	roleCursor, err := r.roleDB.Collection("roles").Find(r.ctx, bson.M{})
	if err != nil {
		return err
	}
	defer roleCursor.Close(r.ctx)
	var roles []bson.M
	if err := roleCursor.All(r.ctx, &roles); err != nil {
		return err
	}
	if len(roles) > 4 {
		roles = roles[:4]
	}

	now := time.Now()
	for i, u := range users {
		uid := getInt32FromBsonM(u, "_id")
		rid := getInt32FromBsonM(roles[i%len(roles)], "_id")

		id, err := database.GetNextID(r.ctx, r.roleDB, "user_roles")
		if err != nil {
			r.logger.Error("failed to get next user_role id", zap.Error(err))
			return err
		}

		ur := models.UserRole{
			UserRoleID: id,
			UserID:     uid,
			RoleID:     rid,
			CreatedAt:  &now,
			UpdatedAt:  &now,
		}

		_, err = r.roleDB.Collection("user_roles").InsertOne(r.ctx, ur)
		if err != nil {
			r.logger.Error("failed to assign role", zap.Error(err))
			return err
		}
	}

	r.logger.Info("user roles seeded successfully")
	return nil
}

func getInt32FromBsonM(m bson.M, key string) int32 {
	v, ok := m[key]
	if !ok {
		return 0
	}
	switch t := v.(type) {
	case int32:
		return t
	case int64:
		return int32(t)
	case float64:
		return int32(t)
	default:
		return 0
	}
}
