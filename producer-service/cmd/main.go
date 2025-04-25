package main

import (
	"log"
	"net"

	"E-Commerce/producer-service/config"
	"E-Commerce/producer-service/internal/handler"
	"E-Commerce/producer-service/internal/repository"
	"E-Commerce/producer-service/internal/service"
	pb "E-Commerce/producer-service/proto"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.NewConfig()

	repo := repository.NewRabbitMQRepository(cfg.RabbitMQConn)
	svc := service.NewProducerService(repo)
	h := handler.NewProducerHandler(svc)

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterProducerServiceServer(s, h)

	log.Println("producer-service started, listening on :50053")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}