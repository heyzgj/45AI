package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/45ai/backend/internal/config"
	"github.com/45ai/backend/internal/handler"
	"github.com/45ai/backend/internal/middleware"
	"github.com/45ai/backend/internal/repository"
	"github.com/45ai/backend/internal/service"
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
	var db *database.DB
	if cfg.App.Environment == "development" {
		// Use SQLite for development
		db, err = database.NewSQLiteConnection("./database.db")
		if err != nil {
			log.Fatal("Failed to connect to SQLite database:", err)
		}
	} else {
		// Use MySQL for production
		db, err = database.NewConnection(cfg.Database)
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}
	}
	defer db.Close()

	// Run migrations for all environments
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
	router.Use(middleware.GetCORSMiddleware(cfg.App.Environment))

	// Serve static files for generated images with security and performance headers
	staticFileHandler := func(c *gin.Context) {
		filename := c.Param("filename")
		
		// Security: Only allow PNG files
		if len(filename) < 4 || filename[len(filename)-4:] != ".png" {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}
		
		// Set cache headers for better performance
		c.Header("Cache-Control", "public, max-age=86400") // 24 hours
		c.Header("Content-Type", "image/png")
		
		// Serve the file
		c.File("uploads/generated/" + filename)
	}
	
	router.GET("/uploads/generated/:filename", staticFileHandler)
	router.HEAD("/uploads/generated/:filename", staticFileHandler)

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

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB)
	wechatRepo := repository.NewWechatRepository(cfg.WeChat)
	templateRepo := repository.NewTemplateRepository(db.DB)
	transactionRepo := repository.NewTransactionRepository(db.DB)
	generationRepo := repository.NewGenerationRepository(db.DB)
	
	// Use Gemini for image generation in development, otherwise use the mock
	var comfyuiRepo repository.ComfyUIRepository
	if cfg.App.Environment == "development" {
		var err error
		comfyuiRepo, err = repository.NewGeminiRepository()
		if err != nil {
			log.Println("WARNING: Could not initialize Gemini Repository, falling back to mock. Error:", err)
			comfyuiRepo = repository.NewMockComfyUIRepository()
		} else {
			log.Println("Using Gemini Repository for image generation.")
		}
	} else {
		comfyuiRepo = repository.NewMockComfyUIRepository()
	}

	// Initialize services
	authService := service.NewAuthService(cfg.JWT, userRepo, wechatRepo)
	templateService := service.NewTemplateService(templateRepo)
	userService := service.NewUserService(userRepo)
	transactionService := service.NewTransactionService(transactionRepo)
	contentSafetyService := service.NewMockContentSafetyService()
	logger := slog.Default()
	queueService := service.NewChannelQueueService(generationRepo)
	generationService := service.NewGenerationService(contentSafetyService, userRepo, transactionRepo, templateRepo, comfyuiRepo, generationRepo)
	creditService := service.NewCreditService(userRepo, transactionRepo, logger)
	wechatPayService := service.NewMockWechatPayService(creditService, logger)
	appleIAPService := service.NewMockAppleIAPService(creditService, logger)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	templateHandler := handler.NewTemplateHandler(templateService)
	userHandler := handler.NewUserHandler(userService, transactionService)
	generationHandler := handler.NewGenerationHandler(generationService, queueService)
	paymentHandler := handler.NewPaymentHandler(wechatPayService, appleIAPService)

	// Initialize middleware
	authMiddleware := middleware.AuthMiddleware(authService)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
		}

		templates := v1.Group("/templates")
		{
			templates.GET("", templateHandler.GetAll)
			templates.GET("/:id", templateHandler.GetByID)
		}

		me := v1.Group("/me")
		me.Use(authMiddleware)
		{
			me.GET("", userHandler.GetProfile)
			me.PUT("", userHandler.UpdateProfile)
			me.GET("/transactions", userHandler.GetTransactions)
		}

		generation := v1.Group("/generate")
		generation.Use(authMiddleware)
		{
			generation.POST("", generationHandler.GenerateImage)
			generation.GET("/:job_id/status", generationHandler.GetGenerationStatus)
			generation.GET("/:job_id", generationHandler.GetGenerationResult)
		}

		billing := v1.Group("/billing")
		billing.Use(authMiddleware)
		{
			billing.POST("/purchase", paymentHandler.SimplePurchase)
		}
	}

	// Start queue workers
	ctx := context.Background()
	queueService.StartWorkers(ctx, 2, generationService) // Start 2 workers

	return router
} 