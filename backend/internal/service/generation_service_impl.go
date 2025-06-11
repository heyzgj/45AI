package service

import (
	"context"
	"fmt"
	"io"
	"time"
	"github.com/45ai/backend/internal/model"
	"github.com/45ai/backend/internal/repository"
)

type generationServiceImpl struct {
	contentSafetyService ContentSafetyService
	userRepo             repository.UserRepository
	transactionRepo      repository.TransactionRepository
	templateRepo         repository.TemplateRepository
	comfyuiRepo          repository.ComfyUIRepository
	generationRepo       repository.GenerationRepository
}

func NewGenerationService(
	contentSafetyService ContentSafetyService,
	userRepo repository.UserRepository,
	transactionRepo repository.TransactionRepository,
	templateRepo repository.TemplateRepository,
	comfyuiRepo repository.ComfyUIRepository,
	generationRepo repository.GenerationRepository,
) GenerationService {
	return &generationServiceImpl{
		contentSafetyService: contentSafetyService,
		userRepo:             userRepo,
		transactionRepo:      transactionRepo,
		templateRepo:         templateRepo,
		comfyuiRepo:          comfyuiRepo,
		generationRepo:       generationRepo,
	}
}

func (s *generationServiceImpl) GenerateImage(ctx context.Context, userID int64, templateID int, imageData io.Reader) (*GenerationResult, error) {
	// 1. Validate the user's uploaded image
	if err := s.ValidateImage(ctx, imageData); err != nil {
		return nil, err
	}

	// 2. Check content safety
	if err := s.CheckContentSafety(ctx, imageData); err != nil {
		return nil, err
	}

	// 3. Check user credits
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	template, err := s.templateRepo.GetByID(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}
	if user.Credits < template.CreditCost {
		return nil, fmt.Errorf("insufficient credits")
	}

	// 4. Generate image
	imageURLs, err := s.comfyuiRepo.GenerateImage(ctx, templateID, imageData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate image: %w", err)
	}

	// 5. Deduct credits and create transaction
	if err := s.userRepo.UpdateCredits(ctx, userID, -template.CreditCost); err != nil {
		return nil, fmt.Errorf("failed to deduct credits: %w", err)
	}
	templateIDPtr := &template.ID
	transaction := &model.Transaction{
		UserID:           userID,
		Type:             "generation",
		Amount:           -template.CreditCost,
		Description:      fmt.Sprintf("Used '%s' template", template.Name),
		RelatedTemplateID: templateIDPtr,
	}
	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		// This is a critical error, as the user has been charged but the transaction was not recorded.
		// In a real application, this should be handled with more care (e.g., a retry mechanism or a refund).
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// For the synchronous version, we'll return the first image URL
	var imageURL string
	if len(imageURLs) > 0 {
		imageURL = imageURLs[0]
	}

	return &GenerationResult{
		JobID:    "sync-" + fmt.Sprintf("%d", time.Now().Unix()), // Temporary job ID for sync calls
		ImageURL: imageURL,
		Status:   "completed",
	}, nil
}

func (s *generationServiceImpl) ValidateImage(ctx context.Context, imageData io.Reader) error {
	// For now, we'll just do a basic check for nil.
	if imageData == nil {
		return fmt.Errorf("image data is required")
	}
	return nil
}

func (s *generationServiceImpl) CheckContentSafety(ctx context.Context, imageData io.Reader) error {
	safe, err := s.contentSafetyService.ValidateImage(ctx, imageData)
	if err != nil {
		return fmt.Errorf("content safety check failed: %w", err)
	}
	if !safe {
		return fmt.Errorf("image content is not safe")
	}
	return nil
}

func (s *generationServiceImpl) GetGenerationStatus(ctx context.Context, jobID string) (*GenerationStatus, error) {
	generation, err := s.generationRepo.GetByJobID(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("failed to get generation: %w", err)
	}

	status := &GenerationStatus{
		JobID:    generation.JobID,
		Status:   generation.Status,
		Progress: generation.Progress,
		ImageURL: generation.ImageURL,
		Error:    generation.Error,
	}

	return status, nil
}

func (s *generationServiceImpl) GetGenerationResult(ctx context.Context, jobID string) (*GenerationResult, error) {
	generation, err := s.generationRepo.GetByJobID(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("failed to get generation: %w", err)
	}

	if generation.Status != "completed" {
		return nil, fmt.Errorf("generation not completed yet, status: %s", generation.Status)
	}

	result := &GenerationResult{
		JobID:    generation.JobID,
		ImageURL: generation.ImageURL,
		Status:   generation.Status,
	}

	return result, nil
} 