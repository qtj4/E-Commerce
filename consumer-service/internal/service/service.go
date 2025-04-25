package service

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
	"E-Commerce/consumer-service/internal/entity"
	"E-Commerce/consumer-service/internal/repository"
)

// ConsumerService defines the interface for the consumer service
type ConsumerService interface {
	ProcessOrderCreated(msgs <-chan amqp.Delivery)
}

// consumerService implements ConsumerService
type consumerService struct {
	rabbitRepo repository.RabbitMQRepository
	grpcRepo   repository.GRPCRepository
}

// NewConsumerService creates a new consumer service
func NewConsumerService(rabbitRepo repository.RabbitMQRepository, grpcRepo repository.GRPCRepository) ConsumerService {
	return &consumerService{rabbitRepo: rabbitRepo, grpcRepo: grpcRepo}
}

// ProcessOrderCreated processes order.created events from RabbitMQ
func (s *consumerService) ProcessOrderCreated(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		var event entity.OrderCreatedEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			continue
		}

		log.Printf("Received order.created event for order %s with products: %v", event.OrderID, event.Products)

		for _, productID := range event.Products {
			// Decrease stock by 1 (adjust quantity as needed)
			err := s.grpcRepo.UpdateStock(productID, -1)
			if err != nil {
				log.Printf("Failed to update stock for product %s: %v", productID, err)
			} else {
				log.Printf("Successfully updated stock for product %s", productID)
			}
		}
	}
}