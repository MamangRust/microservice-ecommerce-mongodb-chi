package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/dotenv"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"go.uber.org/zap"
)

func main() {
	flags := flag.NewFlagSet("migrate", flag.ExitOnError)
	flags.Usage = func() {
		fmt.Println("Usage: migrate up")
		fmt.Println("\nCreates MongoDB indexes and counters for the ecommerce database.")
	}

	if err := flags.Parse([]string{}); err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}

	err := dotenv.Viper()
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	l, err := logger.NewLogger("migrate", nil)
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}

	// EnsureIndexes membuka satu koneksi per bounded context ("1 database = 1
	// MongoDB"), jadi tidak ada client bersama di sini.
	if err := database.EnsureIndexes(l); err != nil {
		l.Fatal("Failed to create indexes", zap.Error(err))
	}

	l.Info("Migration completed successfully")
}
