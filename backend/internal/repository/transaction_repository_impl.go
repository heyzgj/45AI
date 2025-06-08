package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/45ai/backend/internal/model"
)

type transactionRepositoryImpl struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepositoryImpl{db: db}
}

func (r *transactionRepositoryImpl) Create(ctx context.Context, transaction *model.Transaction) error {
	query := "INSERT INTO transactions (user_id, type, amount, description, external_payment_id, related_template_id) VALUES (?, ?, ?, ?, ?, ?)"
	result, err := r.db.ExecContext(ctx, query, transaction.UserID, transaction.Type, transaction.Amount, transaction.Description, transaction.ExternalPaymentID, transaction.RelatedTemplateID)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	transaction.ID = id
	return nil
}

func (r *transactionRepositoryImpl) GetByID(ctx context.Context, id int64) (*model.Transaction, error) {
	// To be implemented in a future task
	return nil, nil
}

func (r *transactionRepositoryImpl) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]model.Transaction, error) {
	query := "SELECT id, user_id, type, amount, description, external_payment_id, related_template_id, created_at FROM transactions WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?"
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var t model.Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.Type, &t.Amount, &t.Description, &t.ExternalPaymentID, &t.RelatedTemplateID, &t.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (r *transactionRepositoryImpl) GetByUserIDAndDateRange(ctx context.Context, userID int64, start, end time.Time) ([]model.Transaction, error) {
	// To be implemented in a future task
	return nil, nil
}

func (r *transactionRepositoryImpl) CountByUserID(ctx context.Context, userID int64) (int, error) {
	// To be implemented in a future task
	return 0, nil
}

func (r *transactionRepositoryImpl) SumCreditsByUserID(ctx context.Context, userID int64) (int, error) {
	// To be implemented in a future task
	return 0, nil
} 