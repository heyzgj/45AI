package model

import (
	"time"
)

// Generation represents a queued or completed image generation job
type Generation struct {
	ID          int64     `json:"id" db:"id"`
	JobID       string    `json:"job_id" db:"job_id"`
	UserID      int64     `json:"user_id" db:"user_id"`
	TemplateID  int       `json:"template_id" db:"template_id"`
	Status      string    `json:"status" db:"status"` // "pending", "processing", "completed", "failed"
	Progress    int       `json:"progress" db:"progress"` // 0-100
	ImageURL    string    `json:"image_url,omitempty" db:"image_url"`
	Error       string    `json:"error,omitempty" db:"error"`
	StartedAt   *time.Time `json:"started_at,omitempty" db:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// GenerationStatus represents the public status response
type GenerationStatus struct {
	JobID     string `json:"job_id"`
	Status    string `json:"status"`
	Progress  int    `json:"progress"`
	ImageURL  string `json:"image_url,omitempty"`
	Error     string `json:"error,omitempty"`
}

// GenerationResult represents the public result response
type GenerationResult struct {
	JobID    string `json:"job_id"`
	ImageURL string `json:"image_url"`
	Status   string `json:"status"`
} 