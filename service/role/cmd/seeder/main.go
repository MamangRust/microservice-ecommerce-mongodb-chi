package main

import (
	"context"
	"fmt"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/dotenv"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"

	"github.com/MamangRust/microservice-ecommerce-grpc-role/seeder"
)

// open connects to the server owning the given bounded context.
func open(log logger.LoggerInterface, dbName string) (*mongo.Database, error) {
	client, err := database.NewMongoClientForContext(log, dbName)
	if err != nil {
		return nil, fmt.Errorf("connect to %s: %w", dbName, err)
	}
	return database.NewMongoDatabase(client, dbName), nil
}

func main() {
	log, err := logger.NewLogger("seeder", nil)
	if err != nil {
		log.Fatal("Failed to initialize logger", zap.Error(err))
	}

	if err := dotenv.Viper(); err != nil {
		log.Fatal("Failed to load .env file", zap.Error(err))
	}

	ctx := context.Background()

	// roles dan users berada di bounded context yang sama (ec_identity), jadi
	// cukup satu koneksi — bukan satu per domain seperti dulu.
	identityDB, err := open(log, database.IdentityDB)
	if err != nil {
		log.Fatal("Failed to connect to identity database", zap.Error(err))
	}

	role := seeder.NewRoleSeeder(identityDB, ctx, log)
	if err := role.Seed(); err != nil {
		log.Fatal("Failed to seed roles", zap.Error(err))
	}

	userRole := seeder.NewUserRoleSeeder(identityDB, identityDB, ctx, log)
	if err := userRole.Seed(); err != nil {
		log.Fatal("Failed to seed user_roles", zap.Error(err))
	}

	log.Info("roles and user_roles seeded successfully")
}
