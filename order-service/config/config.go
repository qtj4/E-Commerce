package config

import (
	"log"
	"os"

	pbInventory "E-Commerce/inventory-service/proto"
	"E-Commerce/order-service/internal/utils"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

type Config struct {
	DB              *sqlx.DB
	InventoryClient pbInventory.InventoryServiceClient
	EmailConfig     utils.EmailConfig
}

func NewConfig() *Config {
	dsn := os.Getenv("POSTGRES_URL")
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Inventory Service: %v", err)
	}
	inventoryClient := pbInventory.NewInventoryServiceClient(conn)

	emailConfig := utils.EmailConfig{
		SenderEmail:    "e_book_aitu@zohomail.com",
		SenderPassword: "gakon2006",
		SMTPHost:       "smtp.zoho.com",
		SMTPPort:       "587",
	}

	return &Config{
		DB:              db,
		InventoryClient: inventoryClient,
		EmailConfig:     emailConfig,
	}
}
