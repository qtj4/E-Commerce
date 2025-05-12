package main

import (
	"log"
	"net"

	"E-Commerce/inventory-service/config"
	"E-Commerce/inventory-service/internal/handler"
	"E-Commerce/inventory-service/internal/repository"
	"E-Commerce/inventory-service/internal/service"
	pb "E-Commerce/inventory-service/proto"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.NewConfig()
	defer cfg.DB.Close()
	defer cfg.Redis.Close()

	repo := repository.NewProductRepository(cfg.DB, cfg.Redis)
	svc := service.NewInventoryService(repo)
	h := handler.NewInventoryGRPCServer(svc)

	lis, err := net.Listen("tcp", ":50055")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterInventoryServiceServer(s, h)

	log.Println("inventory-service started, listening on :50055")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
