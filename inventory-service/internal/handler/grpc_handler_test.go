package handler

import (
	"E-Commerce/inventory-service/config"
	"E-Commerce/inventory-service/internal/repository"
	"E-Commerce/inventory-service/internal/service"
	pb "E-Commerce/inventory-service/proto"
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupTestServer(t *testing.T) *InventoryGRPCServer {
	// Ensure we have the test database URL
	os.Setenv("POSTGRES_URL", "postgres://postgres:postgres@localhost:5432/inventory_test?sslmode=disable")

	cfg := config.NewConfig()
	repo := repository.NewProductRepository(cfg.DB, cfg.Redis)
	svc := service.NewInventoryService(repo)
	return NewInventoryGRPCServer(svc)
}

func TestProductStockManagement(t *testing.T) {
	server := setupTestServer(t)
	ctx := context.Background()

	// First create a test product
	createReq := &pb.CreateProductRequest{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		Stock:       100,
		CategoryId:  "test-category",
	}

	createResp, err := server.CreateProduct(ctx, createReq)
	assert.NoError(t, err)
	assert.NotNil(t, createResp)
	assert.Equal(t, createReq.Name, createResp.Product.Name)
	assert.Equal(t, int32(100), createResp.Product.Stock)

	productID := createResp.Product.Id

	// Test CheckStock
	checkStockReq := &pb.CheckStockRequest{
		ProductId: productID,
		Quantity:  50,
	}

	checkResp, err := server.CheckStock(ctx, checkStockReq)
	assert.NoError(t, err)
	assert.True(t, checkResp.Available)

	// Test UpdateStock
	updateStockReq := &pb.UpdateStockRequest{
		ProductId: productID,
		Quantity:  -30,
		OrderId:   uuid.New().String(),
	}

	updateResp, err := server.UpdateStock(ctx, updateStockReq)
	assert.NoError(t, err)
	assert.True(t, updateResp.Success)

	// Verify the stock was updated correctly
	getReq := &pb.GetProductRequest{
		Id: productID,
	}
	getResp, err := server.GetProduct(ctx, getReq)
	assert.NoError(t, err)
	assert.Equal(t, int32(70), getResp.Product.Stock)
}

func TestProductCRUD(t *testing.T) {
	server := setupTestServer(t)
	ctx := context.Background()

	// Test Create
	createReq := &pb.CreateProductRequest{
		Name:        "CRUD Test Product",
		Description: "CRUD Test Description",
		Price:       149.99,
		Stock:       50,
		CategoryId:  "test-category",
	}

	createResp, err := server.CreateProduct(ctx, createReq)
	assert.NoError(t, err)
	assert.NotNil(t, createResp)
	assert.Equal(t, createReq.Name, createResp.Product.Name)
	productID := createResp.Product.Id

	// Test Get
	getReq := &pb.GetProductRequest{
		Id: productID,
	}
	getResp, err := server.GetProduct(ctx, getReq)
	assert.NoError(t, err)
	assert.NotNil(t, getResp)
	assert.Equal(t, createReq.Name, getResp.Product.Name)
	assert.Equal(t, createReq.Description, getResp.Product.Description)
	assert.Equal(t, float32(149.99), getResp.Product.Price)
	assert.Equal(t, int32(50), getResp.Product.Stock)

	// Test Update
	updateReq := &pb.UpdateProductRequest{
		Id:          productID,
		Name:        "Updated CRUD Test Product",
		Description: "Updated CRUD Test Description",
		Price:       199.99,
		Stock:       75,
		CategoryId:  "test-category",
	}

	updateResp, err := server.UpdateProduct(ctx, updateReq)
	assert.NoError(t, err)
	assert.NotNil(t, updateResp)
	assert.Equal(t, updateReq.Name, updateResp.Product.Name)
	assert.Equal(t, updateReq.Description, updateResp.Product.Description)
	assert.Equal(t, float32(199.99), updateResp.Product.Price)
	assert.Equal(t, int32(75), updateResp.Product.Stock)

	// Test Delete
	deleteReq := &pb.DeleteProductRequest{
		Id: productID,
	}
	deleteResp, err := server.DeleteProduct(ctx, deleteReq)
	assert.NoError(t, err)
	assert.True(t, deleteResp.Success)

	// Verify deletion
	_, err = server.GetProduct(ctx, getReq)
	assert.Error(t, err)
}
