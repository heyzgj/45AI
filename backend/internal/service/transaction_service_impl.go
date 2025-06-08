package service

import (
	"context"

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