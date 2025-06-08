package model

import (
	"time"
)

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypePurchase   TransactionType = "purchase"
	TransactionTypeGeneration TransactionType = "generation"
)

// Transaction represents a credit transaction
type Transaction struct {
	ID                int64           `json:"id" db:"id"`
	UserID            int64           `json:"user_id" db:"user_id"`
	Type              TransactionType `json:"type" db:"type"`
	Amount            int             `json:"amount" db:"amount"`
	Description       string          `json:"description" db:"description"`
	ExternalPaymentID *string         `json:"external_payment_id,omitempty" db:"external_payment_id"`
	RelatedTemplateID *int            `json:"related_template_id,omitempty" db:"related_template_id"`
	CreatedAt         time.Time       `json:"created_at" db:"created_at"`
}

// TransactionCreateRequest represents the request to create a transaction
type TransactionCreateRequest struct {
	UserID            int64           `json:"user_id" binding:"required"`
	Type              TransactionType `json:"type" binding:"required"`
	Amount            int             `json:"amount" binding:"required"`
	Description       string          `json:"description" binding:"required"`
	ExternalPaymentID *string         `json:"external_payment_id,omitempty"`
	RelatedTemplateID *int            `json:"related_template_id,omitempty"`
} 