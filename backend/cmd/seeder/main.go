package main

import (
	"context"
	"log"

	"github.com/45ai/backend/internal/config"
	"github.com/45ai/backend/internal/model"
	"github.com/45ai/backend/internal/repository"
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

	// Initialize repository
	templateRepo := repository.NewTemplateRepository(db.DB)

	// Seed data
	seedTemplates(templateRepo)
}

func seedTemplates(repo repository.TemplateRepository) {
	ctx := context.Background()

	templates := []model.Template{
		{Name: "Cyberpunk", Description: "Futuristic cityscapes and neon lights.", PreviewImageURL: "/images/cyberpunk.png", CreditCost: 10, IsActive: true},
		{Name: "Van Gogh", Description: "Classic impressionist style.", PreviewImageURL: "/images/vangogh.png", CreditCost: 15, IsActive: true},
		{Name: "Ghibli", Description: "Hayao Miyazaki inspired anime style.", PreviewImageURL: "/images/ghibli.png", CreditCost: 20, IsActive: true},
		{Name: "3D Render", Description: "Pixar-like 3D characters.", PreviewImageURL: "/images/3d.png", CreditCost: 25, IsActive: true},
		{Name: "Watercolor", Description: "Soft and vibrant watercolor painting.", PreviewImageURL: "/images/watercolor.png", CreditCost: 10, IsActive: true},
	}

	for _, t := range templates {
		// Idempotency check
		_, err := repo.GetByName(ctx, t.Name)
		if err == nil {
			log.Printf("Template '%s' already exists, skipping.", t.Name)
			continue
		}

		if err := repo.Create(ctx, &t); err != nil {
			log.Printf("Failed to seed template %s: %v", t.Name, err)
		} else {
			log.Printf("Seeded template: %s", t.Name)
		}
	}
}
