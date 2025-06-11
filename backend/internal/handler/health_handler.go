package handler

import (
	"context"
	"database/sql"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	db        *sql.DB
	startTime time.Time
	version   string
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *sql.DB, version string) *HealthHandler {
	return &HealthHandler{
		db:        db,
		startTime: time.Now(),
		version:   version,
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Version   string                 `json:"version,omitempty"`
	Uptime    string                 `json:"uptime,omitempty"`
	Checks    map[string]HealthCheck `json:"checks,omitempty"`
	System    *SystemInfo            `json:"system,omitempty"`
}

// HealthCheck represents an individual health check
type HealthCheck struct {
	Status  string        `json:"status"`
	Message string        `json:"message,omitempty"`
	Latency time.Duration `json:"latency,omitempty"`
	Error   string        `json:"error,omitempty"`
}

// SystemInfo represents system information
type SystemInfo struct {
	GoVersion    string `json:"go_version"`
	NumGoroutine int    `json:"num_goroutine"`
	MemoryMB     uint64 `json:"memory_mb"`
	NumCPU       int    `json:"num_cpu"`
}

// BasicHealth provides a simple health check endpoint
// @Summary Basic health check
// @Description Returns basic health status
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *HealthHandler) BasicHealth(c *gin.Context) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   h.version,
		Uptime:    time.Since(h.startTime).String(),
	}

	c.JSON(http.StatusOK, response)
}

// DetailedHealth provides comprehensive health information
// @Summary Detailed health check
// @Description Returns detailed health status including system info and dependency checks
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 503 {object} HealthResponse
// @Router /health/detailed [get]
func (h *HealthHandler) DetailedHealth(c *gin.Context) {
	checks := make(map[string]HealthCheck)
	overallStatus := "healthy"

	// Database health check
	dbCheck := h.checkDatabase()
	checks["database"] = dbCheck
	if dbCheck.Status != "healthy" {
		overallStatus = "unhealthy"
	}

	// System information
	sysInfo := h.getSystemInfo()

	response := HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Version:   h.version,
		Uptime:    time.Since(h.startTime).String(),
		Checks:    checks,
		System:    sysInfo,
	}

	statusCode := http.StatusOK
	if overallStatus != "healthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// ReadinessCheck checks if the service is ready to accept traffic
// @Summary Readiness check
// @Description Checks if the service is ready to handle requests
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 503 {object} HealthResponse
// @Router /ready [get]
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	checks := make(map[string]HealthCheck)
	overallStatus := "ready"

	// Check all critical dependencies
	dbCheck := h.checkDatabase()
	checks["database"] = dbCheck
	if dbCheck.Status != "healthy" {
		overallStatus = "not_ready"
	}

	response := HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Checks:    checks,
	}

	statusCode := http.StatusOK
	if overallStatus != "ready" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// LivenessCheck checks if the service is alive
// @Summary Liveness check
// @Description Checks if the service is alive and responding
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /live [get]
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	response := HealthResponse{
		Status:    "alive",
		Timestamp: time.Now(),
		Uptime:    time.Since(h.startTime).String(),
	}

	c.JSON(http.StatusOK, response)
}

// checkDatabase performs a database connectivity check
func (h *HealthHandler) checkDatabase() HealthCheck {
	start := time.Now()
	
	if h.db == nil {
		return HealthCheck{
			Status:  "unhealthy",
			Message: "Database connection not initialized",
			Error:   "db is nil",
			Latency: time.Since(start),
		}
	}

	// Test database connectivity with a simple query
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.db.PingContext(ctx)
	latency := time.Since(start)

	if err != nil {
		return HealthCheck{
			Status:  "unhealthy",
			Message: "Database connection failed",
			Error:   err.Error(),
			Latency: latency,
		}
	}

	// Test with a simple query
	var count int
	err = h.db.QueryRowContext(ctx, "SELECT 1").Scan(&count)
	if err != nil {
		return HealthCheck{
			Status:  "unhealthy",
			Message: "Database query failed",
			Error:   err.Error(),
			Latency: time.Since(start),
		}
	}

	return HealthCheck{
		Status:  "healthy",
		Message: "Database connection successful",
		Latency: time.Since(start),
	}
}

// getSystemInfo retrieves system information
func (h *HealthHandler) getSystemInfo() *SystemInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return &SystemInfo{
		GoVersion:    runtime.Version(),
		NumGoroutine: runtime.NumGoroutine(),
		MemoryMB:     m.Alloc / 1024 / 1024, // Convert bytes to MB
		NumCPU:       runtime.NumCPU(),
	}
}

// Metrics provides basic application metrics
// @Summary Application metrics
// @Description Returns basic application metrics for monitoring
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /metrics [get]
func (h *HealthHandler) Metrics(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	metrics := map[string]interface{}{
		"uptime_seconds":    time.Since(h.startTime).Seconds(),
		"go_version":        runtime.Version(),
		"goroutines":        runtime.NumGoroutine(),
		"memory_alloc_mb":   m.Alloc / 1024 / 1024,
		"memory_total_mb":   m.TotalAlloc / 1024 / 1024,
		"memory_sys_mb":     m.Sys / 1024 / 1024,
		"gc_runs":           m.NumGC,
		"cpu_count":         runtime.NumCPU(),
		"timestamp":         time.Now().Unix(),
	}

	c.JSON(http.StatusOK, metrics)
}

// RegisterRoutes registers health check routes
func (h *HealthHandler) RegisterRoutes(r *gin.Engine) {
	// Basic health check
	r.GET("/health", h.BasicHealth)
	r.GET("/healthz", h.BasicHealth) // Kubernetes style

	// Detailed health check
	r.GET("/health/detailed", h.DetailedHealth)
	
	// Kubernetes health checks
	r.GET("/ready", h.ReadinessCheck)
	r.GET("/live", h.LivenessCheck)
	
	// Metrics endpoint
	r.GET("/metrics", h.Metrics)
} 