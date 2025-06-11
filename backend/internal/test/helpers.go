package test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	_ "github.com/go-sql-driver/mysql"
)

// TestSuite provides a base test suite with common setup
type TestSuite struct {
	suite.Suite
	DB     *sql.DB
	Mock   sqlmock.Sqlmock
	TestDB *sql.DB
}

// SetupSuite runs once before all tests in the suite
func (s *TestSuite) SetupSuite() {
	// Setup mock database
	db, mock, err := sqlmock.New()
	require.NoError(s.T(), err)
	
	s.DB = db
	s.Mock = mock
	
	// Setup test database if needed
	s.setupTestDB()
}

// TearDownSuite runs once after all tests in the suite
func (s *TestSuite) TearDownSuite() {
	if s.DB != nil {
		s.DB.Close()
	}
	if s.TestDB != nil {
		s.TestDB.Close()
	}
}

// SetupTest runs before each test
func (s *TestSuite) SetupTest() {
	// Reset mock expectations
	s.Mock.ExpectationsWereMet()
}

// TearDownTest runs after each test
func (s *TestSuite) TearDownTest() {
	// Verify all mock expectations were met
	err := s.Mock.ExpectationsWereMet()
	assert.NoError(s.T(), err)
}

// setupTestDB creates a test database connection
func (s *TestSuite) setupTestDB() {
	testDBURL := os.Getenv("TEST_DB_URL")
	if testDBURL == "" {
		// Skip real database tests if no test DB configured
		return
	}
	
	db, err := sql.Open("mysql", testDBURL)
	if err != nil {
		log.Printf("Failed to connect to test database: %v", err)
		return
	}
	
	s.TestDB = db
}

// TestHelpers provides utility functions for testing
type TestHelpers struct{}

// NewTestHelpers creates a new test helpers instance
func NewTestHelpers() *TestHelpers {
	return &TestHelpers{}
}

// AssertJSONResponse validates JSON response structure
func (h *TestHelpers) AssertJSONResponse(t *testing.T, body []byte, expectedKeys ...string) {
	// Simple JSON validation - in production use proper JSON parsing
	bodyStr := string(body)
	for _, key := range expectedKeys {
		assert.Contains(t, bodyStr, fmt.Sprintf(`"%s"`, key))
	}
}

// AssertErrorResponse validates error response format
func (h *TestHelpers) AssertErrorResponse(t *testing.T, body []byte, expectedError string) {
	bodyStr := string(body)
	assert.Contains(t, bodyStr, `"error"`)
	assert.Contains(t, bodyStr, expectedError)
}

// MockUser creates a mock user for testing
func (h *TestHelpers) MockUser() map[string]interface{} {
	return map[string]interface{}{
		"id":       1,
		"openid":   "test_openid_123",
		"nickname": "Test User",
		"avatar":   "https://example.com/avatar.jpg",
		"credits":  100,
	}
}

// MockTemplate creates a mock template for testing
func (h *TestHelpers) MockTemplate() map[string]interface{} {
	return map[string]interface{}{
		"id":          1,
		"name":        "Test Template",
		"description": "A test template",
		"preview_url": "https://example.com/preview.jpg",
		"category":    "portrait",
		"credits":     1,
	}
}

// MockTransaction creates a mock transaction for testing
func (h *TestHelpers) MockTransaction() map[string]interface{} {
	return map[string]interface{}{
		"id":          "txn_123",
		"user_id":     1,
		"type":        "purchase",
		"amount":      10,
		"credits":     100,
		"status":      "completed",
		"created_at":  "2024-01-01T00:00:00Z",
	}
}

// DatabaseTestSuite provides database-specific testing utilities
type DatabaseTestSuite struct {
	TestSuite
}

// SetupDatabase creates test tables and data
func (s *DatabaseTestSuite) SetupDatabase() {
	if s.TestDB == nil {
		s.T().Skip("Test database not configured")
		return
	}
	
	// Create test tables
	s.createTestTables()
	s.seedTestData()
}

// CleanupDatabase removes test data
func (s *DatabaseTestSuite) CleanupDatabase() {
	if s.TestDB == nil {
		return
	}
	
	// Clean up test data
	tables := []string{"transactions", "templates", "users"}
	for _, table := range tables {
		_, err := s.TestDB.Exec(fmt.Sprintf("DELETE FROM %s WHERE id LIKE 'test_%%'", table))
		if err != nil {
			log.Printf("Failed to clean table %s: %v", table, err)
		}
	}
}

// createTestTables creates necessary test tables
func (s *DatabaseTestSuite) createTestTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			openid VARCHAR(255) UNIQUE NOT NULL,
			nickname VARCHAR(255),
			avatar VARCHAR(500),
			credits INT DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS templates (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			preview_url VARCHAR(500),
			category VARCHAR(100),
			credits INT DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS transactions (
			id VARCHAR(255) PRIMARY KEY,
			user_id INT NOT NULL,
			type ENUM('purchase', 'generation', 'refund') NOT NULL,
			amount DECIMAL(10,2),
			credits INT,
			status ENUM('pending', 'completed', 'failed') DEFAULT 'pending',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
	}
	
	for _, query := range queries {
		_, err := s.TestDB.Exec(query)
		require.NoError(s.T(), err, "Failed to create test table")
	}
}

// seedTestData inserts test data
func (s *DatabaseTestSuite) seedTestData() {
	// Insert test user
	_, err := s.TestDB.Exec(`
		INSERT INTO users (openid, nickname, avatar, credits) 
		VALUES ('test_openid_123', 'Test User', 'https://example.com/avatar.jpg', 100)
		ON DUPLICATE KEY UPDATE nickname = VALUES(nickname)
	`)
	require.NoError(s.T(), err)
	
	// Insert test template
	_, err = s.TestDB.Exec(`
		INSERT INTO templates (name, description, preview_url, category, credits) 
		VALUES ('Test Template', 'A test template', 'https://example.com/preview.jpg', 'portrait', 1)
		ON DUPLICATE KEY UPDATE name = VALUES(name)
	`)
	require.NoError(s.T(), err)
}

// APITestSuite provides API testing utilities
type APITestSuite struct {
	TestSuite
	BaseURL string
}

// SetupAPI initializes API testing
func (s *APITestSuite) SetupAPI() {
	s.BaseURL = "http://localhost:8080/api/v1"
}

// MockHTTPRequest creates a mock HTTP request for testing
func (s *APITestSuite) MockHTTPRequest(method, path string, body []byte) map[string]interface{} {
	return map[string]interface{}{
		"method": method,
		"url":    s.BaseURL + path,
		"body":   string(body),
		"headers": map[string]string{
			"Content-Type": "application/json",
		},
	}
}

// Performance testing utilities
type PerformanceTestSuite struct {
	TestSuite
}

// BenchmarkFunction provides a wrapper for benchmark testing
func (s *PerformanceTestSuite) BenchmarkFunction(b *testing.B, fn func()) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fn()
	}
}

// LoadTestConfig represents load testing configuration
type LoadTestConfig struct {
	Concurrency int
	Duration    int // seconds
	RPS         int // requests per second
}

// RunLoadTest simulates load testing
func (s *PerformanceTestSuite) RunLoadTest(config LoadTestConfig, testFunc func()) {
	// Simplified load test simulation
	// In production, use proper load testing tools
	for i := 0; i < config.Concurrency; i++ {
		go func() {
			for j := 0; j < config.Duration; j++ {
				testFunc()
			}
		}()
	}
}

// Test utilities for common assertions
func AssertValidUUID(t *testing.T, uuid string) {
	assert.Regexp(t, `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, uuid)
}

func AssertValidTimestamp(t *testing.T, timestamp string) {
	assert.Regexp(t, `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}`, timestamp)
}

func AssertPositiveInteger(t *testing.T, value int) {
	assert.Greater(t, value, 0)
}

func AssertValidEmail(t *testing.T, email string) {
	assert.Regexp(t, `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
}

func AssertValidURL(t *testing.T, url string) {
	assert.Regexp(t, `^https?://`, url)
} 