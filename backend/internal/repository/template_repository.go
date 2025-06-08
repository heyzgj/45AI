package repository

import (
	"context"
	"github.com/45ai/backend/internal/model"
)

// TemplateRepository defines the interface for template data access
type TemplateRepository interface {
	// GetAll retrieves all active templates
	GetAll(ctx context.Context) ([]model.Template, error)
	
	// GetByID retrieves a template by ID
	GetByID(ctx context.Context, id int) (*model.Template, error)
	
	// GetByName retrieves a template by name
	GetByName(ctx context.Context, name string) (*model.Template, error)
	
	// Create creates a new template
	Create(ctx context.Context, template *model.Template) error
	
	// Update updates template information
	Update(ctx context.Context, template *model.Template) error
	
	// SetActive activates or deactivates a template
	SetActive(ctx context.Context, id int, isActive bool) error
	
	// Count returns the total number of active templates
	Count(ctx context.Context) (int, error)
} 