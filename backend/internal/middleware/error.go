package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/45ai/backend/internal/model"
)

// ErrorHandlerMiddleware handles errors and panics in the application
type ErrorHandlerMiddleware struct {
	logger *slog.Logger
}

// NewErrorHandlerMiddleware creates a new error handler middleware
func NewErrorHandlerMiddleware(logger *slog.Logger) *ErrorHandlerMiddleware {
	return &ErrorHandlerMiddleware{
		logger: logger,
	}
}

// RecoveryHandler handles panics and converts them to HTTP 500 responses
func (e *ErrorHandlerMiddleware) RecoveryHandler() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with stack trace
				e.logger.Error("Panic recovered",
					"error", err,
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
					"stack", string(debug.Stack()),
				)

							// Return a generic error response
			c.JSON(http.StatusInternalServerError, model.NewErrorResponse("INTERNAL_SERVER_ERROR", "An unexpected error occurred"))
				c.Abort()
			}
		}()
		c.Next()
	})
}

// ErrorHandler processes errors added to the Gin context
func (e *ErrorHandlerMiddleware) ErrorHandler() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Next()

		// Process any errors that occurred during request handling
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			// Log the error
			e.logger.Error("Request error",
				"error", err.Error(),
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
				"user_agent", c.Request.UserAgent(),
				"trace_id", getTraceID(c),
			)

			// If no response has been written yet, send an error response
			if !c.Writer.Written() {
				statusCode := getStatusCodeFromError(err.Error())
				
				c.JSON(statusCode, model.NewErrorResponse("REQUEST_FAILED", getPublicErrorMessage(err.Error())))
			}
		}
	})
}

// NotFoundHandler handles 404 errors
func (e *ErrorHandlerMiddleware) NotFoundHandler() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, model.NewErrorResponse("NOT_FOUND", fmt.Sprintf("Endpoint %s %s not found", c.Request.Method, c.Request.URL.Path)))
	})
}

// MethodNotAllowedHandler handles 405 errors
func (e *ErrorHandlerMiddleware) MethodNotAllowedHandler() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, model.NewErrorResponse("METHOD_NOT_ALLOWED", fmt.Sprintf("Method %s not allowed for endpoint %s", c.Request.Method, c.Request.URL.Path)))
	})
}

// getStatusCodeFromError determines the appropriate HTTP status code from error message
func getStatusCodeFromError(errorMsg string) int {
	// Add logic to determine status code based on error type
	switch {
	case contains(errorMsg, "validation"):
		return http.StatusBadRequest
	case contains(errorMsg, "unauthorized"):
		return http.StatusUnauthorized
	case contains(errorMsg, "forbidden"):
		return http.StatusForbidden
	case contains(errorMsg, "not found"):
		return http.StatusNotFound
	case contains(errorMsg, "conflict"):
		return http.StatusConflict
	case contains(errorMsg, "too many requests"):
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}

// getPublicErrorMessage returns a user-safe error message
func getPublicErrorMessage(errorMsg string) string {
	// Filter out sensitive information from error messages
	switch {
	case contains(errorMsg, "validation"):
		return "Invalid request data"
	case contains(errorMsg, "unauthorized"):
		return "Authentication required"
	case contains(errorMsg, "forbidden"):
		return "Access denied"
	case contains(errorMsg, "not found"):
		return "Resource not found"
	case contains(errorMsg, "duplicate"):
		return "Resource already exists"
	case contains(errorMsg, "database"):
		return "Data operation failed"
	case contains(errorMsg, "network"):
		return "Network error occurred"
	default:
		return "An error occurred while processing your request"
	}
}

// getTraceID extracts or generates a trace ID for request tracking
func getTraceID(c *gin.Context) string {
	// Try to get trace ID from headers (if using distributed tracing)
	if traceID := c.GetHeader("X-Trace-ID"); traceID != "" {
		return traceID
	}
	
	// Try to get from context (if set by earlier middleware)
	if traceID, exists := c.Get("trace_id"); exists {
		if id, ok := traceID.(string); ok {
			return id
		}
	}
	
	// Generate a simple trace ID (in production, use proper trace ID generation)
	return fmt.Sprintf("trace-%d", c.Request.Context().Value("request_id"))
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
		    (len(s) > len(substr) && 
		     (s[:len(substr)] == substr || 
		      s[len(s)-len(substr):] == substr ||
		      containsSubstring(s, substr))))
}

// containsSubstring performs case-insensitive substring search
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// AppError represents application-specific errors
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Status  int    `json:"-"`
}

// Error implements the error interface
func (ae *AppError) Error() string {
	return ae.Message
}

// NewAppError creates a new application error
func NewAppError(code, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

// Common application errors
var (
	ErrUnauthorized = NewAppError("UNAUTHORIZED", "Authentication required", http.StatusUnauthorized)
	ErrForbidden    = NewAppError("FORBIDDEN", "Access denied", http.StatusForbidden)
	ErrNotFound     = NewAppError("NOT_FOUND", "Resource not found", http.StatusNotFound)
	ErrConflict     = NewAppError("CONFLICT", "Resource already exists", http.StatusConflict)
	ErrBadRequest   = NewAppError("BAD_REQUEST", "Invalid request", http.StatusBadRequest)
	ErrInternal     = NewAppError("INTERNAL_ERROR", "Internal server error", http.StatusInternalServerError)
)

// HandleAppError handles application-specific errors
func (e *ErrorHandlerMiddleware) HandleAppError(c *gin.Context, err *AppError) {
	e.logger.Error("Application error",
		"code", err.Code,
		"message", err.Message,
		"details", err.Details,
		"path", c.Request.URL.Path,
		"method", c.Request.Method,
	)

	c.JSON(err.Status, model.NewErrorResponse(err.Code, err.Message))
} 