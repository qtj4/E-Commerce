package service

import (
	"E-Commerce/producer-service/internal/entity"
	"E-Commerce/producer-service/internal/repository"
)

type ProducerService interface {
	PublishOrderCreated(orderID string, products []string) error
}

type producerService struct {
	repo repository.RabbitMQRepository
}

func NewProducerService(repo repository.RabbitMQRepository) ProducerService {
	return &producerService{repo: repo}
}

func (s *producerService) PublishOrderCreated(orderID string, products []string) error {
	event := entity.OrderCreatedEvent{
		OrderID:  orderID,
		Products: products,
	}
	return s.repo.PublishOrderCreated(event)
}