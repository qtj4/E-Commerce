package handler

import (
	"context"
	"log"

	"E-Commerce/producer-service/internal/service"
	pb "E-Commerce/producer-service/proto"
)

type ProducerHandler struct {
	pb.UnimplementedProducerServiceServer
	svc service.ProducerService
}

func NewProducerHandler(svc service.ProducerService) *ProducerHandler {
	return &ProducerHandler{svc: svc}
}

func (h *ProducerHandler) NotifyOrderCreated(ctx context.Context, req *pb.OrderCreatedRequest) (*pb.OrderCreatedResponse, error) {
	log.Printf("Received order.created notification for order %s with products: %v", req.OrderId, req.ProductIds)
	err := h.svc.PublishOrderCreated(req.OrderId, req.ProductIds)
	if err != nil {
		log.Printf("Failed to publish order.created event: %v", err)
		return &pb.OrderCreatedResponse{Success: false}, err
	}
	log.Printf("Successfully published order.created event for order %s", req.OrderId)
	return &pb.OrderCreatedResponse{Success: true}, nil
}