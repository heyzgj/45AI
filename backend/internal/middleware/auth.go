package middleware

import (
	"net/http"
	"strings"

	"github.com/45ai/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates a middleware function for JWT authentication
func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format, must be Bearer <token>"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		userID, err := authService.ValidateToken(c.Request.Context(), tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user ID in context for use in handlers
		c.Set("userID", userID)

		c.Next()
	}
}

// OptionalAuthMiddleware creates a middleware that doesn't require authentication
// but extracts user info if available
func OptionalAuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString := parts[1]
				userID, err := authService.ValidateToken(c.Request.Context(), tokenString)
				if err == nil {
					c.Set("userID", userID)
				}
			}
		}
		c.Next()
	}
} 