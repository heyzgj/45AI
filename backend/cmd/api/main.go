package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/45ai/backend/internal/config"
	"github.com/45ai/backend/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Connect to database
	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Run migrations
	migrationRunner := database.NewMigrationRunner(db.DB, "./migrations")
	if err := migrationRunner.Migrate(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize Gin router with dependencies
	router := setupRouter(cfg, db)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.App.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %d in %s mode", cfg.App.Port, cfg.App.Environment)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func setupRouter(cfg *config.Config, db *database.DB) *gin.Engine {
	// Set Gin mode based on environment
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Apply global middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Health check endpoints
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"service": "45AI Backend API",
			"environment": cfg.App.Environment,
		})
	})

	// Database health check
	router.GET("/health/db", func(c *gin.Context) {
		if err := db.HealthCheck(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"database": "connected",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Placeholder for future routes
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
		
		// TODO: Add more routes here
		// auth := v1.Group("/auth")
		// templates := v1.Group("/templates")
		// generation := v1.Group("/generate")
		// billing := v1.Group("/billing")
	}

	return router
} 