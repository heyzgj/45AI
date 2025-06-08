package main

import (
	"log"

	"github.com/45ai/backend/internal/config"
	"github.com/45ai/backend/pkg/database"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Connect to database
	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Run migrations
	migrationRunner := database.NewMigrationRunner(db.DB, "./migrations")
	if err := migrationRunner.Migrate(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	log.Println("Migrations completed successfully")
} 