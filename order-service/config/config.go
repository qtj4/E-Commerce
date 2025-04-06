package config

import (
	"log"
	"os"

	pbInventory "E-Commerce/inventory-service/proto"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

// Config holds the configuration for the Order Service
type Config struct {
	DB              *sqlx.DB
	InventoryClient pbInventory.InventoryServiceClient
}

// NewConfig initializes and returns the configuration
func NewConfig() *Config {
	// Connect to PostgreSQL database
	dsn := os.Getenv("POSTGRES_URL")
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Connect to Inventory Service via gRPC
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Inventory Service: %v", err)
	}
	inventoryClient := pbInventory.NewInventoryServiceClient(conn)

	return &Config{
		DB:              db,
		InventoryClient: inventoryClient,
	}
}