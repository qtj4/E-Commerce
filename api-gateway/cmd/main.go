package main

import (
    "log"

    "E-Commerce/api-gateway/internal/middleware"
    pbInventory "E-Commerce/inventory-service/proto"
    "E-Commerce/api-gateway/internal/handler"
	pb "E-Commerce/order-service/proto"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)
func main() {
    inventoryConn, err := grpc.Dial("inventory-service:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect to Inventory Service: %v", err)
    }
    defer inventoryConn.Close()

    orderConn, err := grpc.Dial("order-service:50052", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect to Order Service: %v", err)
    }
    defer orderConn.Close()

    inventoryClient := pbInventory.NewInventoryServiceClient(inventoryConn)
    orderClient := pb.NewOrderServiceClient(orderConn)

    r := gin.Default()
    r.Use(middleware.Auth()) // Add authentication middleware

    h := handler.NewRESTHandler(inventoryClient, orderClient)

    // Inventory Routes
    r.POST("/products", h.CreateProduct)
    r.GET("/products/:id", h.GetProduct)
    r.PATCH("/products/:id", h.UpdateProduct)
    r.DELETE("/products/:id", h.DeleteProduct)
    r.GET("/products", h.ListProducts)

    // Order Routes
    r.POST("/orders", h.CreateOrder)
    r.GET("/orders/:id", h.GetOrder)
    r.PATCH("/orders/:id", h.UpdateOrder)
    r.GET("/orders", h.ListOrders)

    if err := r.Run(":8080"); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}