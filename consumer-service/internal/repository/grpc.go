package repository

import (
	"context"

	"E-Commerce/inventory-service/proto"
)

type GRPCRepository interface {
	UpdateStock(productID string, quantity int) error
}

type grpcRepository struct {
	client proto.InventoryServiceClient
}

func NewGRPCRepository(client proto.InventoryServiceClient) GRPCRepository {
	return &grpcRepository{client: client}
}

func (r *grpcRepository) UpdateStock(productID string, quantity int) error {
	_, err := r.client.UpdateStock(context.Background(), &proto.UpdateStockRequest{
		ProductId: productID,
		Quantity:  int32(quantity),
	})
	return err
}