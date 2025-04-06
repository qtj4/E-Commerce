package main

import (
	"log"

	"E-Commerce/api-gateway/internal/handler"
	"E-Commerce/api-gateway/internal/middleware"
	"E-Commerce/api-gateway/internal/repository"
	pbInventory "E-Commerce/inventory-service/proto"
	pb "E-Commerce/order-service/proto"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	inventoryConn, err := grpc.Dial("inventory-service:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Inventory Service: %v", err)
	} else {
		log.Println("Connected to Inventory Service")
	}
	defer inventoryConn.Close()

	orderConn, err := grpc.Dial("order-service:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Order Service: %v", err)
	} else {
		log.Println("Connected to Order Service")
	}
	defer orderConn.Close()

	inventoryClient := pbInventory.NewInventoryServiceClient(inventoryConn)
	orderClient := pb.NewOrderServiceClient(orderConn)

	r := gin.Default()

	authRepo := repository.NewAuthRepository() 
	authHandler := handler.NewAuthHandler(authRepo)

	// Auth routes (unprotected)
	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)

	h := handler.NewRESTHandler(inventoryClient, orderClient)

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.Auth())
	{
		// Admin only routes
		admin := protected.Group("/")
		admin.Use(middleware.RequireRole("admin"))
		{
			admin.POST("/products", h.CreateProduct)
			admin.PATCH("/products/:id", h.UpdateProduct)
			admin.DELETE("/products/:id", h.DeleteProduct)
		}

		// User routes
		protected.GET("/products/:id", h.GetProduct)
		protected.GET("/products", h.ListProducts)
		protected.POST("/orders", h.CreateOrder)
		protected.GET("/orders/:id", h.GetOrder)
		protected.GET("/orders", h.ListOrders)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
