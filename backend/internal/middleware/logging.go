package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware handles structured logging for requests and responses
type LoggingMiddleware struct {
	logger *slog.Logger
}

// NewLoggingMiddleware creates a new logging middleware instance
func NewLoggingMiddleware(logger *slog.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

// RequestResponseLogger logs HTTP requests and responses with structured data
func (l *LoggingMiddleware) RequestResponseLogger() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		startTime := time.Now()
		
		// Generate request ID for tracing
		requestID := generateRequestID()
		c.Set("request_id", requestID)
		
		// Read request body for logging (but restore it for handlers)
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}
		
		// Create response writer wrapper to capture response data
		responseWriter := &responseWriterWrapper{
			ResponseWriter: c.Writer,
			body:          &bytes.Buffer{},
		}
		c.Writer = responseWriter
		
		// Log request
		l.logRequest(c, requestID, requestBody)
		
		// Process request
		c.Next()
		
		// Calculate duration
		duration := time.Since(startTime)
		
		// Log response
		l.logResponse(c, requestID, responseWriter, duration)
	})
}

// logRequest logs incoming HTTP requests
func (l *LoggingMiddleware) logRequest(c *gin.Context, requestID string, body []byte) {
	// Filter sensitive data from request body
	filteredBody := filterSensitiveData(string(body))
	
	l.logger.Info("HTTP Request",
		"request_id", requestID,
		"method", c.Request.Method,
		"path", c.Request.URL.Path,
		"query", c.Request.URL.RawQuery,
		"user_agent", c.Request.UserAgent(),
		"remote_addr", c.ClientIP(),
		"content_type", c.Request.Header.Get("Content-Type"),
		"content_length", c.Request.ContentLength,
		"request_body", truncateString(filteredBody, 1000),
		"headers", filterSensitiveHeaders(c.Request.Header),
	)
}

// logResponse logs HTTP responses
func (l *LoggingMiddleware) logResponse(c *gin.Context, requestID string, rw *responseWriterWrapper, duration time.Duration) {
	statusCode := c.Writer.Status()
	responseSize := rw.body.Len()
	
	// Filter sensitive data from response body
	responseBody := filterSensitiveData(rw.body.String())
	
	// Determine log level based on status code
	logLevel := getLogLevelForStatus(statusCode)
	
	logAttrs := []slog.Attr{
		slog.String("request_id", requestID),
		slog.String("method", c.Request.Method),
		slog.String("path", c.Request.URL.Path),
		slog.Int("status_code", statusCode),
		slog.Duration("duration", duration),
		slog.Int("response_size", responseSize),
		slog.String("response_body", truncateString(responseBody, 1000)),
	}
	
	// Add error information if request failed
	if len(c.Errors) > 0 {
		logAttrs = append(logAttrs, slog.String("errors", c.Errors.String()))
	}
	
	// Add user context if available
	if userID, exists := c.Get("user_id"); exists {
		logAttrs = append(logAttrs, slog.Any("user_id", userID))
	}
	
	switch logLevel {
	case slog.LevelError:
		l.logger.LogAttrs(c.Request.Context(), slog.LevelError, "HTTP Response (Error)", logAttrs...)
	case slog.LevelWarn:
		l.logger.LogAttrs(c.Request.Context(), slog.LevelWarn, "HTTP Response (Warning)", logAttrs...)
	default:
		l.logger.LogAttrs(c.Request.Context(), slog.LevelInfo, "HTTP Response", logAttrs...)
	}
}

// HealthCheckLogger provides lighter logging for health check endpoints
func (l *LoggingMiddleware) HealthCheckLogger() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		startTime := time.Now()
		
		c.Next()
		
		duration := time.Since(startTime)
		
		// Only log if there's an error or if it takes too long
		if c.Writer.Status() >= 400 || duration > time.Second {
			l.logger.Warn("Health check issue",
				"path", c.Request.URL.Path,
				"status", c.Writer.Status(),
				"duration", duration,
			)
		}
	})
}

// responseWriterWrapper wraps gin.ResponseWriter to capture response body
type responseWriterWrapper struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (rw *responseWriterWrapper) Write(data []byte) (int, error) {
	// Write to both the original writer and our buffer
	rw.body.Write(data)
	return rw.ResponseWriter.Write(data)
}

// generateRequestID creates a unique request ID
func generateRequestID() string {
	// In production, use a proper UUID library
	return fmt.Sprintf("req_%d_%d", time.Now().UnixNano(), time.Now().Nanosecond()%1000000)
}

// filterSensitiveData removes or masks sensitive information from strings
func filterSensitiveData(data string) string {
	if data == "" {
		return ""
	}
	
	// List of sensitive fields to filter
	sensitiveFields := []string{
		"password", "secret", "token", "key", "auth", "credential",
		"openid", "session", "cookie", "authorization",
	}
	
	filtered := data
	for _, field := range sensitiveFields {
		// Simple pattern matching - in production, use regex for better accuracy
		if strings.Contains(strings.ToLower(filtered), field) {
			// Mask the value after the field
			filtered = maskSensitiveValue(filtered, field)
		}
	}
	
	return filtered
}

// maskSensitiveValue masks sensitive values in JSON-like strings
func maskSensitiveValue(data, field string) string {
	// Simple masking - replace value with asterisks
	// In production, use proper JSON parsing and regex
	fieldPattern := `"` + field + `"`
	if strings.Contains(strings.ToLower(data), strings.ToLower(fieldPattern)) {
		// This is a simplified approach - in production, use proper JSON masking
		return strings.ReplaceAll(data, field, field+"_masked")
	}
	return data
}

// filterSensitiveHeaders removes sensitive information from HTTP headers
func filterSensitiveHeaders(headers http.Header) map[string][]string {
	filtered := make(map[string][]string)
	
	sensitiveHeaders := map[string]bool{
		"authorization": true,
		"cookie":        true,
		"x-api-key":     true,
		"x-auth-token":  true,
	}
	
	for key, values := range headers {
		lowerKey := strings.ToLower(key)
		if sensitiveHeaders[lowerKey] {
			filtered[key] = []string{"[MASKED]"}
		} else {
			filtered[key] = values
		}
	}
	
	return filtered
}

// getLogLevelForStatus determines appropriate log level based on HTTP status code
func getLogLevelForStatus(statusCode int) slog.Level {
	switch {
	case statusCode >= 500:
		return slog.LevelError
	case statusCode >= 400:
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}

// truncateString truncates a string to a maximum length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// SkipPaths returns a middleware that skips logging for specific paths
func (l *LoggingMiddleware) SkipPaths(paths ...string) gin.HandlerFunc {
	skipMap := make(map[string]bool)
	for _, path := range paths {
		skipMap[path] = true
	}
	
	return gin.HandlerFunc(func(c *gin.Context) {
		if skipMap[c.Request.URL.Path] {
			c.Next()
			return
		}
		
		// Use the regular request/response logger
		l.RequestResponseLogger()(c)
	})
}

// SlowRequestLogger logs requests that exceed a certain duration threshold
func (l *LoggingMiddleware) SlowRequestLogger(threshold time.Duration) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		startTime := time.Now()
		
		c.Next()
		
		duration := time.Since(startTime)
		
		if duration > threshold {
			l.logger.Warn("Slow request detected",
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"duration", duration,
				"threshold", threshold,
				"status_code", c.Writer.Status(),
				"user_agent", c.Request.UserAgent(),
			)
		}
	})
}

// APIMetricsLogger logs API usage metrics
func (l *LoggingMiddleware) APIMetricsLogger() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		startTime := time.Now()
		
		c.Next()
		
		duration := time.Since(startTime)
		
		// Log metrics in a structured format for monitoring systems
		l.logger.Info("API Metrics",
			"timestamp", startTime.Unix(),
			"method", c.Request.Method,
			"endpoint", c.Request.URL.Path,
			"status_code", c.Writer.Status(),
			"duration_ms", duration.Milliseconds(),
			"request_size", c.Request.ContentLength,
			"response_size", c.Writer.Size(),
			"user_agent", c.Request.UserAgent(),
			"remote_addr", c.ClientIP(),
		)
	})
} 