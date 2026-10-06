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

type roleseeder struct {
	db     *mongo.Database
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewRoleSeeder(db *mongo.Database, ctx context.Context, logger logger.LoggerInterface) *roleseeder {
	return &roleseeder{db: db, ctx: ctx, logger: logger}
}

func (r *roleseeder) Seed() error {
	count, err := r.db.Collection("roles").CountDocuments(r.ctx, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		r.logger.Debug("roles already seeded, skipping")
		return nil
	}

	now := time.Now()
	roleNames := []string{"ROLE_ADMIN", "ROLE_MERCHANT", "ROLE_USER", "ROLE_CASHIER"}

	for _, name := range roleNames {
		id, err := database.GetNextID(r.ctx, r.db, "roles")
		if err != nil {
			r.logger.Error("failed to get next role id", zap.Error(err))
			return err
		}
		role := models.Role{
			RoleID:    id,
			RoleName:  name,
			CreatedAt: &now,
			UpdatedAt: &now,
		}
		if _, err := r.db.Collection("roles").InsertOne(r.ctx, role); err != nil {
			r.logger.Error("failed to seed role", zap.String("name", name), zap.Error(err))
			return err
		}
	}

	r.logger.Info("roles seeded successfully")
	return nil
}
