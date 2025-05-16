package handler

import (
	"E-Commerce/inventory-service/internal/entity"
	"E-Commerce/inventory-service/internal/service"
	pb "E-Commerce/inventory-service/proto"
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrProductNotFound = errors.New("product not found")

type InventoryGRPCServer struct {
	pb.UnimplementedInventoryServiceServer
	svc service.InventoryService
}

func NewInventoryGRPCServer(svc service.InventoryService) *InventoryGRPCServer {
	return &InventoryGRPCServer{svc: svc}
}

func (s *InventoryGRPCServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	p := &entity.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       float64(req.Price),
		Stock:       int(req.Stock),
		CategoryID:  req.CategoryId,
	}

	err := s.svc.CreateProduct(p)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create product")
	}

	return &pb.CreateProductResponse{
		Product: &pb.Product{
			Id:          p.ID.String(),
			Name:        p.Name,
			Description: p.Description,
			Price:       float32(p.Price),
			Stock:       int32(p.Stock),
			CategoryId:  p.CategoryID,
		},
	}, nil
}

func (s *InventoryGRPCServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}
	p, err := s.svc.GetProduct(id)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		return nil, status.Error(codes.Internal, "failed to get product")
	}
	if p == nil {
		return nil, status.Error(codes.NotFound, "product not found")
	}
	return &pb.GetProductResponse{
		Product: &pb.Product{
			Id:          p.ID.String(),
			Name:        p.Name,
			Description: p.Description,
			Price:       float32(p.Price),
			Stock:       int32(p.Stock),
			CategoryId:  p.CategoryID,
		},
	}, nil
}

func (s *InventoryGRPCServer) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}

	p := &entity.Product{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Price:       float64(req.Price),
		Stock:       int(req.Stock),
		CategoryID:  req.CategoryId,
	}

	err = s.svc.UpdateProduct(p)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		return nil, status.Error(codes.Internal, "failed to update product")
	}

	return &pb.UpdateProductResponse{
		Product: &pb.Product{
			Id:          p.ID.String(),
			Name:        p.Name,
			Description: p.Description,
			Price:       float32(p.Price),
			Stock:       int32(p.Stock),
			CategoryId:  p.CategoryID,
		},
	}, nil
}

func (s *InventoryGRPCServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}
	err = s.svc.DeleteProduct(id)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		return nil, status.Error(codes.Internal, "failed to delete product")
	}
	return &pb.DeleteProductResponse{Success: true}, nil
}

func (s *InventoryGRPCServer) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	products, total, err := s.svc.ListProducts(req.CategoryId, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list products")
	}

	pbProducts := make([]*pb.Product, len(products))
	for i, p := range products {
		pbProducts[i] = &pb.Product{
			Id:          p.ID.String(),
			Name:        p.Name,
			Description: p.Description,
			Price:       float32(p.Price),
			Stock:       int32(p.Stock),
			CategoryId:  p.CategoryID,
		}
	}

	return &pb.ListProductsResponse{
		Products: pbProducts,
		Total:    int32(total),
	}, nil
}

func (s *InventoryGRPCServer) CheckStock(ctx context.Context, req *pb.CheckStockRequest) (*pb.CheckStockResponse, error) {
	pid, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}
	available, err := s.svc.CheckStock(pid, int(req.Quantity))
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		return nil, status.Error(codes.Internal, "failed to check stock")
	}
	return &pb.CheckStockResponse{Available: available}, nil
}

func (s *InventoryGRPCServer) UpdateStock(ctx context.Context, req *pb.UpdateStockRequest) (*pb.UpdateStockResponse, error) {
	pid, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}
	err = s.svc.UpdateStock(pid, int(req.Quantity), req.OrderId)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		return nil, status.Error(codes.Internal, "failed to update stock")
	}
	return &pb.UpdateStockResponse{Success: true}, nil
}
