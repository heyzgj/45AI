package repository

import (
	"context"
	"time"
	"github.com/45ai/backend/internal/model"
)

// TransactionRepository defines the interface for transaction data access
type TransactionRepository interface {
	// Create creates a new transaction
	Create(ctx context.Context, transaction *model.Transaction) error
	
	// GetByID retrieves a transaction by ID
	GetByID(ctx context.Context, id int64) (*model.Transaction, error)
	
	// GetByUserID retrieves all transactions for a user
	GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]model.Transaction, error)
	
	// GetByUserIDAndDateRange retrieves transactions for a user within a date range
	GetByUserIDAndDateRange(ctx context.Context, userID int64, start, end time.Time) ([]model.Transaction, error)
	
	// CountByUserID returns the total number of transactions for a user
	CountByUserID(ctx context.Context, userID int64) (int, error)
	
	// SumCreditsByUserID calculates the total credits for a user
	SumCreditsByUserID(ctx context.Context, userID int64) (int, error)
} 