package main

import (
	"fmt"
	"net"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	pbInventory "github.com/qtj4/E-Commerce/inventory-service/proto"
	"github.com/qtj4/E-Commerce/order-service/internal/handler"
	"github.com/qtj4/E-Commerce/order-service/internal/repository"
	"github.com/qtj4/E-Commerce/order-service/internal/service"
	pb "github.com/qtj4/E-Commerce/order-service/proto"
	"google.golang.org/grpc"
)

func main() {
	db, err := sqlx.Connect("postgres", "user=postgres dbname=orders sslmode=disable")
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	ic := pbInventory.NewInventoryServiceClient(conn)
	repo := repository.NewOrderRepository(db)
	svc := service.NewOrderService(repo, ic)
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
