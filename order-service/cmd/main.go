package main

import (
	"fmt"
	"log"
	"net"
	"os"

	pbInventory "E-Commerce/inventory-service/proto"
	"E-Commerce/order-service/internal/handler"
	"E-Commerce/order-service/internal/repository"
	"E-Commerce/order-service/internal/service"
	"E-Commerce/order-service/internal/utils"
	pb "E-Commerce/order-service/proto"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	dsn := os.Getenv("POSTGRES_URL")
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic(err)
	}

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	ic := pbInventory.NewInventoryServiceClient(conn)

	// Email configuration from environment variables
	emailConfig := utils.EmailConfig{
		SenderEmail:    "e_book_aitu@zohomail.com",
		SenderPassword: "gakon2006",
		SMTPHost:       "smtp.zoho.com",
		SMTPPort:       "587",
	}

	// Validate email configuration
	if emailConfig.SenderEmail == "" || emailConfig.SenderPassword == "" ||
		emailConfig.SMTPHost == "" || emailConfig.SMTPPort == "" {
		log.Fatal("Missing email configuration. Please set SMTP_EMAIL, SMTP_PASSWORD, SMTP_HOST, and SMTP_PORT environment variables")
	}

	repo := repository.NewOrderRepository(db)
	svc := service.NewOrderService(repo, ic, emailConfig)
	h := handler.NewOrderGRPCServer(svc)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, h)
	fmt.Println("Order Service running on :50052")
	s.Serve(lis)
}
