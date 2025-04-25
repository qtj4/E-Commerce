package main

import (
	"log"

	"E-Commerce/consumer-service/config"
	"E-Commerce/consumer-service/internal/handler"
	"E-Commerce/consumer-service/internal/repository"
	"E-Commerce/consumer-service/internal/service"
)

func main() {
	cfg := config.NewConfig()

	rabbitRepo := repository.NewRabbitMQRepository(cfg.RabbitMQConn)
	grpcRepo := repository.NewGRPCRepository(cfg.ProductClient)
	svc := service.NewConsumerService(rabbitRepo, grpcRepo)
	h := handler.NewConsumerHandler(svc)

	log.Println("consumer-service started, listening for order.created events...")
	h.ProcessOrderCreated()

	select {} // Keep the service running
}