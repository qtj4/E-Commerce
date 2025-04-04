package config

import (
	pbInventory "github.com/yourusername/inventory-service/proto"
	pbOrder "github.com/yourusername/order-service/proto"
	"google.golang.org/grpc"
)

type Config struct {
	InventoryClient pbInventory.InventoryServiceClient
	OrderClient     pbOrder.OrderServiceClient
}

func NewConfig() *Config {
	inventoryConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	orderConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return &Config{
		InventoryClient: pbInventory.NewInventoryServiceClient(inventoryConn),
		OrderClient:     pbOrder.NewOrderServiceClient(orderConn),
	}
}