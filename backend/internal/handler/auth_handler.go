package handler

import (
	"github.com/gin-gonic/gin"
)

// AuthHandler defines the interface for authentication HTTP handlers
type AuthHandler interface {
	// Login handles user login requests
	Login(c *gin.Context)
	
	// Refresh handles token refresh requests
	Refresh(c *gin.Context)
	
	// Logout handles user logout requests
	Logout(c *gin.Context)
	
	// GetProfile retrieves the current user's profile
	GetProfile(c *gin.Context)
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Code string `json:"code" binding:"required"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

// RefreshRequest represents the token refresh request
type RefreshRequest struct {
	Token string `json:"token" binding:"required"`
} 