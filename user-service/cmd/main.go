package main

import (
	"log"
	"net"

	pb "E-Commerce/user-service/proto"
	"E-Commerce/user-service/internal/handler"
	"E-Commerce/user-service/internal/repository"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	// Connect to PostgreSQL
	db, err := sqlx.Connect("postgres", "user=postgres password=secret dbname=ecommerce sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repository and service
	repo := repository.NewUserRepository(db)
	service := handler.NewUserService(repo)

	// Set up gRPC server
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, service)

	log.Println("User Service is running on port 50053")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}