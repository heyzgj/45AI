package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/45ai/backend/pkg/database"
	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	App      AppConfig
	Database database.Config
	JWT      JWTConfig
	WeChat   WeChatConfig
	External ExternalConfig
	Payment  PaymentConfig
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Environment string
	Port        int
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

// WeChatConfig holds WeChat-related configuration
type WeChatConfig struct {
	AppID     string
	AppSecret string
}

// ExternalConfig holds external service configuration
type ExternalConfig struct {
	ContentSafetyAPIKey string
	ContentSafetyAPIURL string
	ComfyUIAPIURL       string
	ComfyUIAPIKey       string
}

// PaymentConfig holds payment-related configuration
type PaymentConfig struct {
	WeChatPayMerchantID string
	WeChatPayAPIKey     string
	AppleIAPSecret      string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if not in production
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			// It's okay if .env doesn't exist in production
			if os.Getenv("APP_ENV") == "" {
				fmt.Println("Warning: .env file not found")
			}
		}
	}

	cfg := &Config{}

	// App configuration
	cfg.App.Environment = getEnv("APP_ENV", "development")
	cfg.App.Port = getEnvInt("PORT", 8080)

	// Database configuration
	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnvInt("DB_PORT", 3306)
	cfg.Database.User = getEnv("DB_USER", "root")
	cfg.Database.Password = getEnv("DB_PASSWORD", "")
	cfg.Database.Database = getEnv("DB_NAME", "45ai_db")
	cfg.Database.MaxOpenConns = getEnvInt("DB_MAX_OPEN_CONNS", 25)
	cfg.Database.MaxIdleConns = getEnvInt("DB_MAX_IDLE_CONNS", 5)
	cfg.Database.ConnMaxLifetime = getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute)
	cfg.Database.ConnMaxIdleTime = getEnvDuration("DB_CONN_MAX_IDLE_TIME", 10*time.Minute)

	// JWT configuration
	cfg.JWT.Secret = getEnv("JWT_SECRET", "")
	if cfg.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	cfg.JWT.Expiry = getEnvDuration("JWT_EXPIRY", 24*time.Hour)

	// WeChat configuration
	cfg.WeChat.AppID = getEnv("WECHAT_APP_ID", "")
	cfg.WeChat.AppSecret = getEnv("WECHAT_APP_SECRET", "")

	// External services
	cfg.External.ContentSafetyAPIKey = getEnv("CONTENT_SAFETY_API_KEY", "")
	cfg.External.ContentSafetyAPIURL = getEnv("CONTENT_SAFETY_API_URL", "")
	cfg.External.ComfyUIAPIURL = getEnv("COMFYUI_API_URL", "http://localhost:8188")
	cfg.External.ComfyUIAPIKey = getEnv("COMFYUI_API_KEY", "")

	// Payment configuration
	cfg.Payment.WeChatPayMerchantID = getEnv("WECHAT_PAY_MERCHANT_ID", "")
	cfg.Payment.WeChatPayAPIKey = getEnv("WECHAT_PAY_API_KEY", "")
	cfg.Payment.AppleIAPSecret = getEnv("APPLE_IAP_SHARED_SECRET", "")

	return cfg, nil
}

// Helper functions for environment variables

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
} 