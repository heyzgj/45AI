package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NewCORSMiddleware creates a CORS middleware configured for the application
func NewCORSMiddleware() gin.HandlerFunc {
	config := cors.Config{
		// Allow specific origins for production
		AllowOrigins: []string{
			"http://localhost:3000",     // Local development
			"http://localhost:8080",     // Vue dev server
			"https://servicewechat.com", // WeChat Mini Program
			"https://*.weixin.qq.com",   // WeChat domains
			"capacitor://localhost",     // Capacitor iOS app
			"ionic://localhost",         // Ionic iOS app
			"http://localhost",          // Local mobile development
		},
		
		// Allow common HTTP methods
		AllowMethods: []string{
			"GET",
			"POST", 
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
			"HEAD",
		},
		
		// Allow headers required for authentication and API requests
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"Accept",
			"Cache-Control",
			"X-Requested-With",
			"X-API-Key",
			"X-Client-Version",
			"X-Platform",
			"X-User-Agent",
		},
		
		// Expose headers that frontend might need
		ExposeHeaders: []string{
			"Content-Length",
			"X-Request-ID",
			"X-Rate-Limit-Remaining",
			"X-Rate-Limit-Reset",
		},
		
		// Allow credentials for authenticated requests
		AllowCredentials: true,
		
		// Cache preflight requests for 1 hour
		MaxAge: time.Hour,
	}
	
	return cors.New(config)
}

// GetCORSMiddleware returns appropriate CORS middleware based on environment
func GetCORSMiddleware(environment string) gin.HandlerFunc {
	if environment == "development" {
		config := cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
			AllowHeaders:     []string{"*"},
			AllowCredentials: true,
			MaxAge:           time.Hour,
		}
		return cors.New(config)
	}
	
	return NewCORSMiddleware()
} 