package service

import (
	"E-Commerce/producer-service/internal/entity"
	"E-Commerce/producer-service/internal/repository"
)

// ProducerService defines the interface for the producer service
type ProducerService interface {
	PublishOrderCreated(orderID string, products []string) error
}

// producerService implements ProducerService
type producerService struct {
	repo repository.RabbitMQRepository
}

// NewProducerService creates a new producer service
func NewProducerService(repo repository.RabbitMQRepository) ProducerService {
	return &producerService{repo: repo}
}

// PublishOrderCreated publishes an order.created event
func (s *producerService) PublishOrderCreated(orderID string, products []string) error {
	event := entity.OrderCreatedEvent{
		OrderID:  orderID,
		Products: products,
	}
	return s.repo.PublishOrderCreated(event)
}