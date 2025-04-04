package config

import (
	"log"

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
	db, err := sqlx.Connect("postgres", "user=postgres dbname=inventory sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return &Config{
		DB: db,
	}
}