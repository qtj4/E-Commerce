package config

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Config holds the configuration for the Inventory Service
type Config struct {
	DB *sqlx.DB
}

// NewConfig initializes and returns the configuration
func NewConfig() *Config {
	// Connect to PostgreSQL database
	dsn := os.Getenv("POSTGRES_URL")
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return &Config{
		DB: db,
	}
}