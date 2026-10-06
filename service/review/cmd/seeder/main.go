package main

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/dotenv"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"go.uber.org/zap"

	"github.com/MamangRust/microservice-ecommerce-grpc-review/seeder"
)

func main() {
	log, err := logger.NewLogger("seeder", nil)
	if err != nil {
		log.Fatal("Failed to initialize logger", zap.Error(err))
	}

	if err := dotenv.Viper(); err != nil {
		log.Fatal("Failed to load .env file", zap.Error(err))
	}

	ctx := context.Background()

	client, err := database.NewMongoClientForContext(log, database.ExperienceDB)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	db := database.NewMongoDatabase(client, database.ExperienceDB)

	s := seeder.NewReviewSeeder(db, ctx, log)
	if err := s.Seed(); err != nil {
		log.Fatal("Failed to seed reviews", zap.Error(err))
	}

	log.Info("reviews seeded successfully")
}
