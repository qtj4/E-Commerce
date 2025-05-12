package config

import (
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	DB    *sqlx.DB
	Redis *redis.Client
}

func NewConfig() *Config {
	// Initialize PostgreSQL
	dsn := os.Getenv("POSTGRES_URL")
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &Config{
		DB:    db,
		Redis: rdb,
	}
}
