package config

import (
	"log"

	"github.com/streadway/amqp"
	"E-Commerce/product-service/proto"
	"google.golang.org/grpc"
)

// Config holds the configuration for the consumer-service
type Config struct {
	RabbitMQConn  *amqp.Connection
	ProductClient proto.ProductServiceClient
}

// NewConfig initializes and returns the configuration
func NewConfig() *Config {
	// Connect to RabbitMQ
	rabbitConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	// Connect to product-service via gRPC
	grpcConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to product-service: %v", err)
	}

	return &Config{
		RabbitMQConn:  rabbitConn,
		ProductClient: proto.NewProductServiceClient(grpcConn),
	}
}