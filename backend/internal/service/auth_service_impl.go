package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/45ai/backend/internal/config"
	"github.com/45ai/backend/internal/model"
	"github.com/45ai/backend/internal/repository"
	"github.com/golang-jwt/jwt/v5"
)

type authServiceImpl struct {
	cfg        config.JWTConfig
	userRepo   repository.UserRepository
	wechatRepo repository.WechatRepository
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(cfg config.JWTConfig, userRepo repository.UserRepository, wechatRepo repository.WechatRepository) AuthService {
	return &authServiceImpl{
		cfg:        cfg,
		userRepo:   userRepo,
		wechatRepo: wechatRepo,
	}
}

func (s *authServiceImpl) LoginWithWechat(ctx context.Context, code string) (*model.User, string, error) {
	var user *model.User
	var err error
	
	// Check for development mock codes
	if code == "dev_mock_code" || strings.HasPrefix(code, "test_") {
		// Mock authentication for development
		mockOpenID := "mock_openid_" + code
		
		// Find or create mock user
		user, err = s.userRepo.GetByWechatOpenID(ctx, mockOpenID)
		if err != nil {
			if err == sql.ErrNoRows {
				// Create a new mock user
				user = &model.User{
					WechatOpenID: mockOpenID,
					Nickname:     "Mock User",
					AvatarURL:    "/static/default-avatar.png",
					Credits:      50, // Give new users 50 free credits
				}
				if createErr := s.userRepo.Create(ctx, user); createErr != nil {
					return nil, "", fmt.Errorf("failed to create mock user: %w", createErr)
				}
			} else {
				return nil, "", fmt.Errorf("failed to get mock user: %w", err)
			}
		}
	} else {
		// Use real WeChat API for production codes
		wechatResp, err := s.wechatRepo.Code2Session(code)
		if err != nil {
			return nil, "", fmt.Errorf("failed to exchange wechat code: %w", err)
		}

		// Find or create user
		user, err = s.userRepo.GetByWechatOpenID(ctx, wechatResp.OpenID)
		if err != nil {
			if err == sql.ErrNoRows {
				// If user not found, create a new one with default credits
				user = &model.User{
					WechatOpenID: wechatResp.OpenID,
					Nickname:     wechatResp.UnionID, // Use WeChat provided data
					AvatarURL:    "",                 // Will be updated from WeChat user info if available
					Credits:      50,                 // Give new users 50 free credits
				}
				if createErr := s.userRepo.Create(ctx, user); createErr != nil {
					return nil, "", fmt.Errorf("failed to create user: %w", createErr)
				}
			} else {
				return nil, "", fmt.Errorf("failed to get user: %w", err)
			}
		}
	}

	// Generate JWT token
	token, err := s.GenerateToken(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

func (s *authServiceImpl) RefreshToken(ctx context.Context, oldToken string) (string, error) {
	// Future implementation
	return "", fmt.Errorf("not implemented")
}

func (s *authServiceImpl) GetUserFromToken(ctx context.Context, token string) (*model.User, error) {
	// Future implementation
	return nil, fmt.Errorf("not implemented")
}

// GenerateToken generates a new JWT for a given user ID
func (s *authServiceImpl) GenerateToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(s.cfg.Expiry).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.Secret))
}

// ValidateToken validates a JWT and returns the user ID
func (s *authServiceImpl) ValidateToken(ctx context.Context, tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.Secret), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["sub"].(float64); ok {
			return int64(userID), nil
		}
		return 0, fmt.Errorf("invalid user ID in token")
	}

	return 0, fmt.Errorf("invalid token")
} 