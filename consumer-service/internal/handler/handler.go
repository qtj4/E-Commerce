package handler

import (
	"log"

	"E-Commerce/consumer-service/internal/service"
)

// ConsumerHandler handles RabbitMQ message processing
type ConsumerHandler struct {
	svc service.ConsumerService
}

// NewConsumerHandler creates a new consumer handler
func NewConsumerHandler(svc service.ConsumerService) *ConsumerHandler {
	return &ConsumerHandler{svc: svc}
}

// ProcessOrderCreated starts processing order.created events
func (h *ConsumerHandler) ProcessOrderCreated() {
	err := h.svc.ProcessOrderCreated()
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}
}