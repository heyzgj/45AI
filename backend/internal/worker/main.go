package main

import (
	"context"
	"log"
	"time"
	"bytes"

	"github.com/45ai/backend/internal/config"
	"github.com/45ai/backend/internal/repository"
	"github.com/45ai/backend/internal/service"
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

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB)
	transactionRepo := repository.NewTransactionRepository(db.DB)
	templateRepo := repository.NewTemplateRepository(db.DB)
	comfyuiRepo := repository.NewMockComfyUIRepository()

	// Initialize services
	contentSafetyService := service.NewMockContentSafetyService()
	generationService := service.NewGenerationService(contentSafetyService, userRepo, transactionRepo, templateRepo, comfyuiRepo)
	queueService := service.NewInMemoryQueueService()

	log.Println("Worker starting...")

	for {
		job, err := queueService.GetJob(context.Background())
		if err != nil {
			log.Printf("Failed to get job from queue: %v", err)
			continue
		}

		if job != nil {
			log.Printf("Processing job for user %d", job.UserID)
			_, err := generationService.GenerateImage(context.Background(), job.UserID, job.TemplateID, bytes.NewReader(job.ImageData))
			if err != nil {
				log.Printf("Failed to process job: %v", err)
			}
		}

		// Wait for a second before checking for new jobs
		time.Sleep(1 * time.Second)
	}
} 