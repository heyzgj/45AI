package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// NewSQLiteConnection creates a new SQLite database connection for development
func NewSQLiteConnection(dbPath string) (*DB, error) {
	// Open SQLite connection
	sqlDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %w", err)
	}

	// Configure connection pool (SQLite doesn't need many connections)
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(0) // No limit for SQLite
	sqlDB.SetConnMaxIdleTime(0) // No limit for SQLite

	// Wrap in our DB type
	db := &DB{
		DB: sqlDB,
		config: Config{
			Database: dbPath,
		},
	}

	// Test connection
	if err := db.Ping(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to ping SQLite database: %w", err)
	}

	return db, nil
} 