package service

import (
	"context"
	"time"

	"github.com/45ai/backend/internal/model"
)

// TransactionService defines the interface for transaction management
type TransactionService interface {
	// GetTransactionsByUserID retrieves transactions for a user with pagination
	GetTransactionsByUserID(ctx context.Context, userID int64, limit, offset int) ([]model.Transaction, error)
	
	// GetTransactionsByUserIDAndDateRange retrieves transactions for a user within a date range
	GetTransactionsByUserIDAndDateRange(ctx context.Context, userID int64, start, end time.Time) ([]model.Transaction, error)
	
	// CountTransactionsByUserID returns the total number of transactions for a user
	CountTransactionsByUserID(ctx context.Context, userID int64) (int, error)
	
	// GetTransactionByID retrieves a specific transaction by ID
	GetTransactionByID(ctx context.Context, id int64) (*model.Transaction, error)
}