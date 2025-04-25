package config

import (
	"log"

	"github.com/streadway/amqp"
	"E-Commerce/inventory-service/proto"
	"google.golang.org/grpc"
)

// Config holds the configuration for the consumer-service
type Config struct {
	RabbitMQConn    *amqp.Connection
	InventoryClient proto.InventoryServiceClient
}

// NewConfig initializes and returns the configuration
func NewConfig() *Config {
	// Connect to RabbitMQ
	rabbitConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	// Connect to inventory-service via gRPC
	grpcConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to inventory-service: %v", err)
	}

	return &Config{
		RabbitMQConn:    rabbitConn,
		InventoryClient: proto.NewInventoryServiceClient(grpcConn),
	}
}