package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Migration represents a database migration
type Migration struct {
	Version   string
	Name      string
	UpSQL     string
	DownSQL   string
	AppliedAt *time.Time
}

// MigrationRunner handles database migrations
type MigrationRunner struct {
	db             *sql.DB
	migrationsPath string
}

// NewMigrationRunner creates a new migration runner
func NewMigrationRunner(db *sql.DB, migrationsPath string) *MigrationRunner {
	return &MigrationRunner{
		db:             db,
		migrationsPath: migrationsPath,
	}
}

// Initialize creates the migrations table if it doesn't exist
func (m *MigrationRunner) Initialize() error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	
	_, err := m.db.Exec(query)
	return err
}

// Migrate runs all pending migrations
func (m *MigrationRunner) Migrate() error {
	if err := m.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize migrations table: %w", err)
	}

	// Get applied migrations
	applied, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Get migration files
	migrations, err := m.getMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Apply pending migrations
	for _, migration := range migrations {
		if _, ok := applied[migration.Version]; ok {
			continue
		}

		fmt.Printf("Applying migration %s: %s\n", migration.Version, migration.Name)
		
		if err := m.applyMigration(migration); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
		}
	}

	return nil
}

// Rollback rolls back the last n migrations
func (m *MigrationRunner) Rollback(n int) error {
	// Implementation would go here
	// This is a placeholder for rollback functionality
	return fmt.Errorf("rollback not implemented yet")
}

// getAppliedMigrations retrieves the list of applied migrations
func (m *MigrationRunner) getAppliedMigrations() (map[string]bool, error) {
	query := "SELECT version FROM schema_migrations"
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

// getMigrationFiles reads migration files from the migrations directory
func (m *MigrationRunner) getMigrationFiles() ([]Migration, error) {
	files, err := ioutil.ReadDir(m.migrationsPath)
	if err != nil {
		return nil, err
	}

	var migrations []Migration
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".up.sql") {
			continue
		}

		// Extract version and name from filename
		// Expected format: 001_create_users_table.up.sql
		parts := strings.Split(file.Name(), "_")
		if len(parts) < 2 {
			continue
		}

		version := parts[0]
		name := strings.TrimSuffix(file.Name(), ".up.sql")

		// Read up migration
		upPath := filepath.Join(m.migrationsPath, file.Name())
		upSQL, err := ioutil.ReadFile(upPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read up migration %s: %w", upPath, err)
		}

		// Read down migration if exists
		downPath := strings.Replace(upPath, ".up.sql", ".down.sql", 1)
		downSQL, _ := ioutil.ReadFile(downPath)

		migrations = append(migrations, Migration{
			Version: version,
			Name:    name,
			UpSQL:   string(upSQL),
			DownSQL: string(downSQL),
		})
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// applyMigration applies a single migration
func (m *MigrationRunner) applyMigration(migration Migration) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute migration
	if _, err := tx.Exec(migration.UpSQL); err != nil {
		return err
	}

	// Record migration
	query := "INSERT INTO schema_migrations (version) VALUES (?)"
	if _, err := tx.Exec(query, migration.Version); err != nil {
		return err
	}

	return tx.Commit()
} 