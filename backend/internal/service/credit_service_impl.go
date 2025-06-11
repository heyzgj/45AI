package service

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/45ai/backend/internal/model"
	"github.com/45ai/backend/internal/repository"
)

// creditServiceImpl implements the CreditService interface
type creditServiceImpl struct {
	userRepo        repository.UserRepository
	transactionRepo repository.TransactionRepository
	logger          *slog.Logger
}

// NewCreditService creates a new credit service instance
func NewCreditService(userRepo repository.UserRepository, transactionRepo repository.TransactionRepository, logger *slog.Logger) CreditService {
	return &creditServiceImpl{
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
		logger:          logger,
	}
}

// AddCredits adds credits to a user's balance and logs the transaction
func (s *creditServiceImpl) AddCredits(ctx context.Context, userID int64, amount int, description string, externalPaymentID *string) error {
	// Validate input
	if amount <= 0 {
		return ErrInvalidAmount
	}

	s.logger.InfoContext(ctx, "Adding credits to user", 
		"user_id", userID, 
		"amount", amount, 
		"description", description)

	// Check if user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrUserNotFound
		}
		s.logger.ErrorContext(ctx, "Failed to get user", "error", err, "user_id", userID)
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Create transaction record
	transaction := &model.Transaction{
		UserID:            userID,
		Type:              model.TransactionTypePurchase,
		Amount:            amount,
		Description:       description,
		ExternalPaymentID: externalPaymentID,
		CreatedAt:         time.Now(),
	}

	// Log the transaction first
	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create transaction", "error", err, "user_id", userID)
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	// Update user credits
	if err := s.userRepo.UpdateCredits(ctx, userID, amount); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update user credits", "error", err, "user_id", userID)
		return fmt.Errorf("failed to update user credits: %w", err)
	}

	s.logger.InfoContext(ctx, "Successfully added credits", 
		"user_id", userID, 
		"amount", amount, 
		"new_balance", user.Credits+amount)

	return nil
}

// DeductCredits deducts credits from a user's balance and logs the transaction
func (s *creditServiceImpl) DeductCredits(ctx context.Context, userID int64, amount int, description string, templateID *int) error {
	// Validate input
	if amount <= 0 {
		return ErrInvalidAmount
	}

	s.logger.InfoContext(ctx, "Deducting credits from user", 
		"user_id", userID, 
		"amount", amount, 
		"description", description)

	// Check if user has sufficient credits
	if err := s.ValidateBalance(ctx, userID, amount); err != nil {
		return err
	}

	// Create transaction record with negative amount
	transaction := &model.Transaction{
		UserID:            userID,
		Type:              model.TransactionTypeGeneration,
		Amount:            -amount,
		Description:       description,
		RelatedTemplateID: templateID,
		CreatedAt:         time.Now(),
	}

	// Log the transaction first
	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create transaction", "error", err, "user_id", userID)
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	// Update user credits (deduct)
	if err := s.userRepo.UpdateCredits(ctx, userID, -amount); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update user credits", "error", err, "user_id", userID)
		return fmt.Errorf("failed to update user credits: %w", err)
	}

	s.logger.InfoContext(ctx, "Successfully deducted credits", 
		"user_id", userID, 
		"amount", amount)

	return nil
}

// GetBalance returns the current credit balance for a user
func (s *creditServiceImpl) GetBalance(ctx context.Context, userID int64) (int, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrUserNotFound
		}
		s.logger.ErrorContext(ctx, "Failed to get user", "error", err, "user_id", userID)
		return 0, fmt.Errorf("failed to get user: %w", err)
	}

	return user.Credits, nil
}

// ValidateBalance checks if a user has sufficient credits for a transaction
func (s *creditServiceImpl) ValidateBalance(ctx context.Context, userID int64, requiredAmount int) error {
	if requiredAmount <= 0 {
		return ErrInvalidAmount
	}

	currentBalance, err := s.GetBalance(ctx, userID)
	if err != nil {
		return err
	}

	if currentBalance < requiredAmount {
		s.logger.WarnContext(ctx, "Insufficient credits", 
			"user_id", userID, 
			"current_balance", currentBalance, 
			"required_amount", requiredAmount)
		return ErrInsufficientCredits
	}

	return nil
} 