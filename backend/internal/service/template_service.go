package service

import (
	"context"
	"github.com/45ai/backend/internal/model"
)

// TemplateService defines the interface for template business logic
type TemplateService interface {
	// GetAllTemplates retrieves all active templates
	GetAllTemplates(ctx context.Context) (*model.TemplateListResponse, error)
	
	// GetTemplateByID retrieves a specific template
	GetTemplateByID(ctx context.Context, id int) (*model.Template, error)
	
	// ValidateTemplateForUser checks if a user can use a template
	ValidateTemplateForUser(ctx context.Context, userID int64, templateID int) error
	
	// GetTemplateRequirements returns the requirements for using a template
	GetTemplateRequirements(ctx context.Context, templateID int) (credits int, err error)
} 