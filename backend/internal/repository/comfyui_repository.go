package repository

import (
	"context"
	"io"
)

type ComfyUIRepository interface {
	GenerateImage(ctx context.Context, templateID int, imageData io.Reader) ([]string, error)
}

type mockComfyUIRepository struct{}

func NewMockComfyUIRepository() ComfyUIRepository {
	return &mockComfyUIRepository{}
}

func (r *mockComfyUIRepository) GenerateImage(ctx context.Context, templateID int, imageData io.Reader) ([]string, error) {
	// In a real implementation, this would call the ComfyUI API.
	// For now, we'll just return some mock image URLs.
	return []string{
		"https://example.com/image1.png",
		"https://example.com/image2.png",
		"https://example.com/image3.png",
	}, nil
} 