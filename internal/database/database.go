package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	DB *sql.DB
}

// NewDatabase creates a new SQLite database connection
func NewDatabase() (*Database, error) {
	// Get database path from environment variable or use default
	dbPath := getEnv("DB_PATH", getDefaultDBPath())

	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("error creating database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	database := &Database{DB: db}

	// Initialize schema
	if err := database.initSchema(); err != nil {
		return nil, fmt.Errorf("error initializing schema: %w", err)
	}

	return database, nil
}

// initSchema creates the database tables
func (d *Database) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date DATETIME NOT NULL,
		description TEXT NOT NULL,
		amount REAL NOT NULL,
		category TEXT NOT NULL,
		type TEXT NOT NULL CHECK(type IN ('income', 'expense')),
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(date);
	CREATE INDEX IF NOT EXISTS idx_transactions_category ON transactions(category);
	CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(type);

	CREATE TABLE IF NOT EXISTS budgets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		category TEXT NOT NULL,
		amount REAL NOT NULL,
		period TEXT NOT NULL CHECK(period IN ('custom')),
		start_date DATETIME NOT NULL,
		end_date DATETIME NOT NULL,
		UNIQUE(category, start_date, end_date)
	);

	CREATE INDEX IF NOT EXISTS idx_budgets_category ON budgets(category);
	CREATE INDEX IF NOT EXISTS idx_budgets_category_period ON budgets(category, period);
	`

	_, err := d.DB.Exec(schema)
	return err
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.DB.Close()
}

// getDefaultDBPath returns the default database path
func getDefaultDBPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./atad.db"
	}
	return filepath.Join(homeDir, ".atad", "atad.db")
}

// Helper function to get environment variables with default values
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
