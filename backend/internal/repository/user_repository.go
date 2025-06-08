package repository

import (
	"context"
	"github.com/45ai/backend/internal/model"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *model.User) error
	
	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id int64) (*model.User, error)
	
	// GetByWechatOpenID retrieves a user by WeChat OpenID
	GetByWechatOpenID(ctx context.Context, openID string) (*model.User, error)
	
	// Update updates user information
	Update(ctx context.Context, user *model.User) error
	
	// UpdateCredits updates user credit balance
	UpdateCredits(ctx context.Context, userID int64, amount int) error
	
	// Exists checks if a user exists by WeChat OpenID
	Exists(ctx context.Context, openID string) (bool, error)
} 