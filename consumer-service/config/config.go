package config

import (
	"log"

	"github.com/streadway/amqp"
	"E-Commerce/inventory-service/proto"
	"google.golang.org/grpc"
)

type Config struct {
	RabbitMQConn    *amqp.Connection
	InventoryClient proto.InventoryServiceClient
}

func NewConfig() *Config {
	rabbitConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	grpcConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to inventory-service: %v", err)
	}

	return &Config{
		RabbitMQConn:    rabbitConn,
		InventoryClient: proto.NewInventoryServiceClient(grpcConn),
	}
}