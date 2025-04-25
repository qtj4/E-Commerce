package repository

import (
	"context"

	"E-Commerce/inventory-service/proto"
)

// GRPCRepository defines the interface for gRPC operations
type GRPCRepository interface {
	UpdateStock(productID string, quantity int) error
}

// grpcRepository implements GRPCRepository
type grpcRepository struct {
	client proto.InventoryServiceClient
}

// NewGRPCRepository creates a new gRPC repository
func NewGRPCRepository(client proto.InventoryServiceClient) GRPCRepository {
	return &grpcRepository{client: client}
}

// UpdateStock calls the inventory-service to update stock
func (r *grpcRepository) UpdateStock(productID string, quantity int) error {
	_, err := r.client.UpdateStock(context.Background(), &proto.UpdateStockRequest{
		ProductId: productID,
		Quantity:  int32(quantity),
	})
	return err
}