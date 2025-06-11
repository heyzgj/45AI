package handler

import (
	"net/http"

	"github.com/45ai/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type authHandlerImpl struct {
	authService service.AuthService
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &authHandlerImpl{
		authService: authService,
	}
}

func (h *authHandlerImpl) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.authService.LoginWithWechat(c.Request.Context(), req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User:  user,
	})
}

func (h *authHandlerImpl) Refresh(c *gin.Context) {
	// Implementation for a future task
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (h *authHandlerImpl) Logout(c *gin.Context) {
	// For JWT-based authentication, logout is handled client-side
	// The server just acknowledges the logout request
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

func (h *authHandlerImpl) GetProfile(c *gin.Context) {
	// Implementation for a future task
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

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