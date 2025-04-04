package config

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	pbInventory "github.com/qtj4/E-Commerce/inventory-service/proto"
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
	db, err := sqlx.Connect("postgres", "user=postgres dbname=orders sslmode=disable")
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