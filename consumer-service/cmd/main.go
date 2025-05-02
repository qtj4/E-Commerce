package main

import (
    "log"
    "os"
    pbInventory "E-Commerce/inventory-service/proto"
    "E-Commerce/consumer-service/internal/handler"
    "E-Commerce/consumer-service/internal/repository"
    "E-Commerce/consumer-service/internal/service"
    "google.golang.org/grpc"
)

func main() {
    // Connect to RabbitMQ
    rabbitURL := os.Getenv("RABBITMQ_URL")
    if rabbitURL == "" {
        rabbitURL = "amqp://guest:guest@localhost:5672/"
    }

    // Connect to Inventory Service
    inventoryConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect to inventory service: %v", err)
    }
    defer inventoryConn.Close()
    inventoryClient := pbInventory.NewInventoryServiceClient(inventoryConn)

    // Initialize repositories
    rabbitRepo, err := repository.NewRabbitMQRepository(rabbitURL)
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
    }
    grpcRepo := repository.NewGRPCRepository(inventoryClient)

    // Initialize service and handler
    svc := service.NewConsumerService(rabbitRepo, grpcRepo)
    h := handler.NewConsumerHandler(svc)

    // Start consumer
    if err := h.Start(); err != nil {
        log.Fatalf("Failed to start consumer: %v", err)
    }
}