package service

import (
    "encoding/json"
    "E-Commerce/consumer-service/internal/entity"
    "E-Commerce/consumer-service/internal/repository"
    "github.com/streadway/amqp"
    "log"
)

type ConsumerService interface {
    ConsumeOrderCreated() (<-chan amqp.Delivery, error)
    ProcessOrderCreated(msg amqp.Delivery) error
}

type consumerService struct {
    rabbitRepo repository.RabbitMQRepository
    grpcRepo   repository.GRPCRepository
}

func NewConsumerService(rabbitRepo repository.RabbitMQRepository, grpcRepo repository.GRPCRepository) ConsumerService {
    return &consumerService{
        rabbitRepo: rabbitRepo,
        grpcRepo:   grpcRepo,
    }
}

func (s *consumerService) ConsumeOrderCreated() (<-chan amqp.Delivery, error) {
    return s.rabbitRepo.ConsumeOrderCreated()
}

func (s *consumerService) ProcessOrderCreated(msg amqp.Delivery) error {
    var event entity.OrderCreatedEvent
    if err := json.Unmarshal(msg.Body, &event); err != nil {
        log.Printf("Error unmarshaling message: %v", err)
        return err
    }

    for _, productID := range event.Products {
        if err := s.grpcRepo.UpdateStock(productID, -1); err != nil {
            log.Printf("Error updating stock for product %s: %v", productID, err)
            return err
        }
    }

    return msg.Ack(false)
}