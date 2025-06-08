package model

import (
	"time"
)

// Template represents an AI style template
type Template struct {
	ID              int       `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	Description     string    `json:"description" db:"description"`
	PreviewImageURL string    `json:"preview_image_url" db:"preview_image_url"`
	CreditCost      int       `json:"credit_cost" db:"credit_cost"`
	IsActive        bool      `json:"is_active" db:"is_active"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// TemplateListResponse represents the response for template listing
type TemplateListResponse struct {
	Templates []Template `json:"templates"`
	Total     int        `json:"total"`
} 