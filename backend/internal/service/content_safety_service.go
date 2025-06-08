package service

import (
	"context"
	"io"
)

type ContentSafetyService interface {
	ValidateImage(ctx context.Context, image io.Reader) (bool, error)
}

type mockContentSafetyService struct{}

func NewMockContentSafetyService() ContentSafetyService {
	return &mockContentSafetyService{}
}

func (s *mockContentSafetyService) ValidateImage(ctx context.Context, image io.Reader) (bool, error) {
	// In a real implementation, this would call a third-party API.
	// For now, we'll just assume all images are safe.
	return true, nil
} 