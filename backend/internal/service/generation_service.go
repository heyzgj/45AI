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
	GetGenerationStatus(ctx context.Context, requestID string) (*GenerationStatus, error)
}

// GenerationResult represents the result of an image generation
type GenerationResult struct {
	RequestID string   `json:"request_id"`
	Images    []string `json:"images"`
	Credits   int      `json:"credits_used"`
}

// GenerationStatus represents the status of a generation request
type GenerationStatus struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"` // "processing", "completed", "failed"
	Progress  int    `json:"progress,omitempty"`
	Error     string `json:"error,omitempty"`
} 