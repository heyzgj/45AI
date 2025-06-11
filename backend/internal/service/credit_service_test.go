package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/45ai/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByWechatOpenID(ctx context.Context, openID string) (*model.User, error) {
	args := m.Called(ctx, openID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateCredits(ctx context.Context, userID int64, amount int) error {
	args := m.Called(ctx, userID, amount)
	return args.Error(0)
}

func (m *MockUserRepository) Exists(ctx context.Context, openID string) (bool, error) {
	args := m.Called(ctx, openID)
	return args.Bool(0), args.Error(1)
}

// MockTransactionRepository is a mock implementation of TransactionRepository
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(ctx context.Context, transaction *model.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetByID(ctx context.Context, id int64) (*model.Transaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]model.Transaction, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]model.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByUserIDAndDateRange(ctx context.Context, userID int64, start, end time.Time) ([]model.Transaction, error) {
	args := m.Called(ctx, userID, start, end)
	return args.Get(0).([]model.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) CountByUserID(ctx context.Context, userID int64) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func (m *MockTransactionRepository) SumCreditsByUserID(ctx context.Context, userID int64) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func TestCreditService_AddCredits(t *testing.T) {
	tests := []struct {
		name              string
		userID            int64
		amount            int
		description       string
		externalPaymentID *string
		mockUser          *model.User
		userRepoError     error
		transactionError  error
		updateCreditsError error
		expectedError     error
	}{
		{
			name:        "successful credit addition",
			userID:      1,
			amount:      100,
			description: "Purchase 100 credits",
			mockUser: &model.User{
				ID:      1,
				Credits: 50,
			},
			expectedError: nil,
		},
		{
			name:          "invalid amount",
			userID:        1,
			amount:        0,
			description:   "Invalid amount",
			expectedError: ErrInvalidAmount,
		},
		{
			name:          "negative amount",
			userID:        1,
			amount:        -50,
			description:   "Negative amount",
			expectedError: ErrInvalidAmount,
		},
		{
			name:          "user not found",
			userID:        999,
			amount:        100,
			description:   "Purchase for non-existent user",
			userRepoError: sql.ErrNoRows,
			expectedError: ErrUserNotFound,
		},
		{
			name:        "transaction creation fails",
			userID:      1,
			amount:      100,
			description: "Purchase 100 credits",
			mockUser: &model.User{
				ID:      1,
				Credits: 50,
			},
			transactionError: errors.New("transaction creation failed"),
			expectedError:    errors.New("transaction creation failed"),
		},
		{
			name:        "credit update fails",
			userID:      1,
			amount:      100,
			description: "Purchase 100 credits",
			mockUser: &model.User{
				ID:      1,
				Credits: 50,
			},
			updateCreditsError: errors.New("credit update failed"),
			expectedError:      errors.New("credit update failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			userRepo := new(MockUserRepository)
			transactionRepo := new(MockTransactionRepository)
			logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
			
			service := NewCreditService(userRepo, transactionRepo, logger)
			ctx := context.Background()

			// Setup expectations based on test case
			if tt.amount > 0 {
				if tt.userRepoError != nil {
					userRepo.On("GetByID", ctx, tt.userID).Return(nil, tt.userRepoError)
				} else if tt.mockUser != nil {
					userRepo.On("GetByID", ctx, tt.userID).Return(tt.mockUser, nil)
					
					if tt.transactionError != nil {
						transactionRepo.On("Create", ctx, mock.AnythingOfType("*model.Transaction")).Return(tt.transactionError)
					} else {
						transactionRepo.On("Create", ctx, mock.AnythingOfType("*model.Transaction")).Return(nil)
						
						if tt.updateCreditsError != nil {
							userRepo.On("UpdateCredits", ctx, tt.userID, tt.amount).Return(tt.updateCreditsError)
						} else {
							userRepo.On("UpdateCredits", ctx, tt.userID, tt.amount).Return(nil)
						}
					}
				}
			}

			// Execute test
			err := service.AddCredits(ctx, tt.userID, tt.amount, tt.description, tt.externalPaymentID)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				if tt.expectedError == ErrInvalidAmount || tt.expectedError == ErrUserNotFound {
					assert.Equal(t, tt.expectedError, err)
				} else {
					assert.Contains(t, err.Error(), tt.expectedError.Error())
				}
			} else {
				assert.NoError(t, err)
			}

			// Verify mock expectations
			userRepo.AssertExpectations(t)
			transactionRepo.AssertExpectations(t)
		})
	}
}

func TestCreditService_DeductCredits(t *testing.T) {
	tests := []struct {
		name              string
		userID            int64
		amount            int
		description       string
		templateID        *int
		mockUser          *model.User
		userRepoError     error
		transactionError  error
		updateCreditsError error
		expectedError     error
	}{
		{
			name:        "successful credit deduction",
			userID:      1,
			amount:      30,
			description: "Used template",
			templateID:  intPtr(5),
			mockUser: &model.User{
				ID:      1,
				Credits: 50,
			},
			expectedError: nil,
		},
		{
			name:          "invalid amount",
			userID:        1,
			amount:        0,
			description:   "Invalid amount",
			expectedError: ErrInvalidAmount,
		},
		{
			name:        "insufficient credits",
			userID:      1,
			amount:      100,
			description: "Used expensive template",
			mockUser: &model.User{
				ID:      1,
				Credits: 50,
			},
			expectedError: ErrInsufficientCredits,
		},
		{
			name:          "user not found",
			userID:        999,
			amount:        30,
			description:   "Deduct for non-existent user",
			userRepoError: sql.ErrNoRows,
			expectedError: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			userRepo := new(MockUserRepository)
			transactionRepo := new(MockTransactionRepository)
			logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
			
			service := NewCreditService(userRepo, transactionRepo, logger)
			ctx := context.Background()

			// Setup expectations
			if tt.amount > 0 {
				if tt.userRepoError != nil {
					userRepo.On("GetByID", ctx, tt.userID).Return(nil, tt.userRepoError)
				} else if tt.mockUser != nil {
					userRepo.On("GetByID", ctx, tt.userID).Return(tt.mockUser, nil)
					
					// If user has sufficient credits, expect transaction and update
					if tt.mockUser.Credits >= tt.amount && tt.expectedError != ErrInsufficientCredits {
						if tt.transactionError != nil {
							transactionRepo.On("Create", ctx, mock.AnythingOfType("*model.Transaction")).Return(tt.transactionError)
						} else {
							transactionRepo.On("Create", ctx, mock.AnythingOfType("*model.Transaction")).Return(nil)
							
							if tt.updateCreditsError != nil {
								userRepo.On("UpdateCredits", ctx, tt.userID, -tt.amount).Return(tt.updateCreditsError)
							} else {
								userRepo.On("UpdateCredits", ctx, tt.userID, -tt.amount).Return(nil)
							}
						}
					}
				}
			}

			// Execute test
			err := service.DeductCredits(ctx, tt.userID, tt.amount, tt.description, tt.templateID)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				if tt.expectedError == ErrInvalidAmount || tt.expectedError == ErrUserNotFound || tt.expectedError == ErrInsufficientCredits {
					assert.Equal(t, tt.expectedError, err)
				} else {
					assert.Contains(t, err.Error(), tt.expectedError.Error())
				}
			} else {
				assert.NoError(t, err)
			}

			// Verify mock expectations
			userRepo.AssertExpectations(t)
			transactionRepo.AssertExpectations(t)
		})
	}
}

func TestCreditService_GetBalance(t *testing.T) {
	tests := []struct {
		name          string
		userID        int64
		mockUser      *model.User
		userRepoError error
		expectedBalance int
		expectedError error
	}{
		{
			name:   "successful balance retrieval",
			userID: 1,
			mockUser: &model.User{
				ID:      1,
				Credits: 150,
			},
			expectedBalance: 150,
			expectedError:   nil,
		},
		{
			name:            "user not found",
			userID:          999,
			userRepoError:   sql.ErrNoRows,
			expectedBalance: 0,
			expectedError:   ErrUserNotFound,
		},
		{
			name:            "database error",
			userID:          1,
			userRepoError:   errors.New("database connection failed"),
			expectedBalance: 0,
			expectedError:   errors.New("database connection failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			userRepo := new(MockUserRepository)
			transactionRepo := new(MockTransactionRepository)
			logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
			
			service := NewCreditService(userRepo, transactionRepo, logger)
			ctx := context.Background()

			// Setup expectations
			if tt.userRepoError != nil {
				userRepo.On("GetByID", ctx, tt.userID).Return(nil, tt.userRepoError)
			} else {
				userRepo.On("GetByID", ctx, tt.userID).Return(tt.mockUser, nil)
			}

			// Execute test
			balance, err := service.GetBalance(ctx, tt.userID)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				if tt.expectedError == ErrUserNotFound {
					assert.Equal(t, tt.expectedError, err)
				} else {
					assert.Contains(t, err.Error(), tt.expectedError.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBalance, balance)
			}

			// Verify mock expectations
			userRepo.AssertExpectations(t)
		})
	}
}

func TestCreditService_ValidateBalance(t *testing.T) {
	tests := []struct {
		name            string
		userID          int64
		requiredAmount  int
		mockUser        *model.User
		userRepoError   error
		expectedError   error
	}{
		{
			name:           "sufficient balance",
			userID:         1,
			requiredAmount: 30,
			mockUser: &model.User{
				ID:      1,
				Credits: 50,
			},
			expectedError: nil,
		},
		{
			name:           "insufficient balance",
			userID:         1,
			requiredAmount: 100,
			mockUser: &model.User{
				ID:      1,
				Credits: 50,
			},
			expectedError: ErrInsufficientCredits,
		},
		{
			name:           "invalid required amount",
			userID:         1,
			requiredAmount: 0,
			expectedError:  ErrInvalidAmount,
		},
		{
			name:           "user not found",
			userID:         999,
			requiredAmount: 30,
			userRepoError:  sql.ErrNoRows,
			expectedError:  ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			userRepo := new(MockUserRepository)
			transactionRepo := new(MockTransactionRepository)
			logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
			
			service := NewCreditService(userRepo, transactionRepo, logger)
			ctx := context.Background()

			// Setup expectations
			if tt.requiredAmount > 0 {
				if tt.userRepoError != nil {
					userRepo.On("GetByID", ctx, tt.userID).Return(nil, tt.userRepoError)
				} else if tt.mockUser != nil {
					userRepo.On("GetByID", ctx, tt.userID).Return(tt.mockUser, nil)
				}
			}

			// Execute test
			err := service.ValidateBalance(ctx, tt.userID, tt.requiredAmount)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}

			// Verify mock expectations
			userRepo.AssertExpectations(t)
		})
	}
}

// Helper function to create int pointer
func intPtr(i int) *int {
	return &i
} 