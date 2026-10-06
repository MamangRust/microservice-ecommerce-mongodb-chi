package seeder

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-pkg/hash"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

type userSeeder struct {
	db     *mongo.Database
	hash   hash.HashPassword
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewUserSeeder(db *mongo.Database, hash hash.HashPassword, ctx context.Context, logger logger.LoggerInterface) *userSeeder {
	return &userSeeder{db: db, hash: hash, ctx: ctx, logger: logger}
}

func (r *userSeeder) Seed() error {
	count, err := r.db.Collection("users").CountDocuments(r.ctx, bson.M{"deleted_at": bson.M{"$eq": nil}})
	if err != nil {
		return err
	}
	if count > 0 {
		r.logger.Debug("users already seeded, skipping")
		return nil
	}

	now := time.Now()
	for i := 1; i <= 10; i++ {
		email := fmt.Sprintf("user_%s@example.com", uuid.NewString())
		rawPassword := fmt.Sprintf("password%d", i)

		hashedPassword, err := r.hash.HashPassword(rawPassword)
		if err != nil {
			r.logger.Error("failed to hash password", zap.Int("user", i), zap.Error(err))
			return fmt.Errorf("failed to hash password for user %d: %w", i, err)
		}

		id, err := database.GetNextID(r.ctx, r.db, "users")
		if err != nil {
			r.logger.Error("failed to get next user id", zap.Int("user", i), zap.Error(err))
			return err
		}

		user := models.User{
			UserID:    id,
			Firstname: fmt.Sprintf("User%d", i),
			Lastname:  fmt.Sprintf("Last%d", i),
			Email:     email,
			Password:  hashedPassword,
			CreatedAt: &now,
			UpdatedAt: &now,
		}

		if _, err := r.db.Collection("users").InsertOne(r.ctx, user); err != nil {
			r.logger.Error("failed to seed user", zap.Int("user", i), zap.Error(err))
			return fmt.Errorf("failed to seed user %d: %w", i, err)
		}
	}

	r.logger.Info("users seeded successfully")
	return nil
}
