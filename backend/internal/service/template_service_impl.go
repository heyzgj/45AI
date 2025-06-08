package service

import (
	"context"
	"github.com/45ai/backend/internal/model"
	"github.com/45ai/backend/internal/repository"
)

type templateServiceImpl struct {
	repo repository.TemplateRepository
}

func NewTemplateService(repo repository.TemplateRepository) TemplateService {
	return &templateServiceImpl{repo: repo}
}

func (s *templateServiceImpl) GetAllTemplates(ctx context.Context) (*model.TemplateListResponse, error) {
	templates, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return &model.TemplateListResponse{
		Templates: templates,
		Total:     len(templates),
	}, nil
}

func (s *templateServiceImpl) GetTemplateByID(ctx context.Context, id int) (*model.Template, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *templateServiceImpl) ValidateTemplateForUser(ctx context.Context, userID int64, templateID int) error {
	// To be implemented in a future task
	return nil
}

func (s *templateServiceImpl) GetTemplateRequirements(ctx context.Context, templateID int) (credits int, err error) {
	// To be implemented in a future task
	return 0, nil
} 