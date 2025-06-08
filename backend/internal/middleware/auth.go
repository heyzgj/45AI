package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates a middleware function for JWT authentication
func AuthMiddleware(authService interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Check Bearer prefix
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format",
			})
			c.Abort()
			return
		}

		token := parts[1]

		// TODO: Validate token using authService
		// For now, we'll just set a placeholder user ID
		// userID, err := authService.ValidateToken(c.Request.Context(), token)
		
		// Set user ID in context for use in handlers
		c.Set("userID", int64(1)) // Placeholder
		c.Set("token", token)

		c.Next()
	}
}

// OptionalAuthMiddleware creates a middleware that doesn't require authentication
// but extracts user info if available
func OptionalAuthMiddleware(authService interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				token := parts[1]
				// TODO: Validate token and set user ID if valid
				c.Set("userID", int64(1)) // Placeholder
				c.Set("token", token)
			}
		}
		c.Next()
	}
} 