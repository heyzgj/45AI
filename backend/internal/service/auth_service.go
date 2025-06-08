package service

import (
	"context"
	"github.com/45ai/backend/internal/model"
)

// AuthService defines the interface for authentication business logic
type AuthService interface {
	// LoginWithWechat authenticates a user with WeChat code
	LoginWithWechat(ctx context.Context, code string) (*model.User, string, error)
	
	// ValidateToken validates a JWT token and returns the user ID
	ValidateToken(ctx context.Context, token string) (int64, error)
	
	// RefreshToken refreshes an existing JWT token
	RefreshToken(ctx context.Context, oldToken string) (string, error)
	
	// GetUserFromToken retrieves user information from a JWT token
	GetUserFromToken(ctx context.Context, token string) (*model.User, error)
}

// WechatLoginResponse represents the response from WeChat login
type WechatLoginResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid,omitempty"`
} 