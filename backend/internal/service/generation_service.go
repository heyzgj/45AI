package service

import (
	"context"
	"io"
)

// GenerationService defines the interface for image generation business logic
type GenerationService interface {
	// GenerateImage processes an image generation request
	GenerateImage(ctx context.Context, userID int64, templateID int, imageData io.Reader) (*GenerationResult, error)
	
	// ValidateImage checks if an uploaded image is suitable for generation
	ValidateImage(ctx context.Context, imageData io.Reader) error
	
	// CheckContentSafety verifies image content is appropriate
	CheckContentSafety(ctx context.Context, imageData io.Reader) error
	
	// GetGenerationStatus retrieves the status of an ongoing generation
	GetGenerationStatus(ctx context.Context, jobID string) (*GenerationStatus, error)
	
	// GetGenerationResult retrieves the result of a completed generation
	GetGenerationResult(ctx context.Context, jobID string) (*GenerationResult, error)
}

// GenerationResult represents the result of an image generation
type GenerationResult struct {
	JobID    string `json:"job_id"`
	ImageURL string `json:"image_url"`
	Status   string `json:"status"`
}

// GenerationStatus represents the status of a generation request
type GenerationStatus struct {
	JobID    string `json:"job_id"`
	Status   string `json:"status"` // "pending", "processing", "completed", "failed"
	Progress int    `json:"progress"`
	ImageURL string `json:"image_url,omitempty"`
	Error    string `json:"error,omitempty"`
} 