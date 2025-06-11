package service

import (
	"context"
	"errors"
)

// CreditService defines the interface for credit management
type CreditService interface {
	// AddCredits adds credits to a user's balance and logs the transaction
	AddCredits(ctx context.Context, userID int64, amount int, description string, externalPaymentID *string) error
	
	// DeductCredits deducts credits from a user's balance and logs the transaction
	// Returns an error if the user doesn't have sufficient credits
	DeductCredits(ctx context.Context, userID int64, amount int, description string, templateID *int) error
	
	// GetBalance returns the current credit balance for a user
	GetBalance(ctx context.Context, userID int64) (int, error)
	
	// ValidateBalance checks if a user has sufficient credits for a transaction
	ValidateBalance(ctx context.Context, userID int64, requiredAmount int) error
}

// Common credit service errors
var (
	ErrInsufficientCredits = errors.New("insufficient credits")
	ErrInvalidAmount       = errors.New("invalid amount: must be positive")
	ErrUserNotFound        = errors.New("user not found")
) 