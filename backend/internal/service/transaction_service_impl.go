package service

import (
	"context"
	"time"

	"github.com/45ai/backend/internal/model"
	"github.com/45ai/backend/internal/repository"
)

type transactionServiceImpl struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionServiceImpl{repo: repo}
}

func (s *transactionServiceImpl) GetTransactionsByUserID(ctx context.Context, userID int64, limit, offset int) ([]model.Transaction, error) {
	return s.repo.GetByUserID(ctx, userID, limit, offset)
}

func (s *transactionServiceImpl) GetTransactionsByUserIDAndDateRange(ctx context.Context, userID int64, start, end time.Time) ([]model.Transaction, error) {
	return s.repo.GetByUserIDAndDateRange(ctx, userID, start, end)
}

func (s *transactionServiceImpl) CountTransactionsByUserID(ctx context.Context, userID int64) (int, error) {
	return s.repo.CountByUserID(ctx, userID)
}

func (s *transactionServiceImpl) GetTransactionByID(ctx context.Context, id int64) (*model.Transaction, error) {
	return s.repo.GetByID(ctx, id)
} 