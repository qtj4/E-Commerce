package handler

import (
	"context"
	"testing"

	"E-Commerce/inventory-service/internal/entity"
	"E-Commerce/inventory-service/internal/service"
	pb "E-Commerce/inventory-service/proto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Mock repository
type mockProductRepository struct {
	mock.Mock
}

func (m *mockProductRepository) Create(p *entity.Product) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *mockProductRepository) Get(id uuid.UUID) (*entity.Product, error) {
	args := m.Called(id)
	if p := args.Get(0); p != nil {
		return p.(*entity.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockProductRepository) Update(p *entity.Product) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *mockProductRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockProductRepository) List(categoryID string, page, pageSize int) ([]*entity.Product, int, error) {
	args := m.Called(categoryID, page, pageSize)
	return args.Get(0).([]*entity.Product), args.Int(1), args.Error(2)
}

func (m *mockProductRepository) CheckStock(productID uuid.UUID, quantity int) (bool, error) {
	args := m.Called(productID, quantity)
	return args.Bool(0), args.Error(1)
}

func (m *mockProductRepository) UpdateStock(productID uuid.UUID, quantity int, orderID string) error {
	args := m.Called(productID, quantity, orderID)
	return args.Error(0)
}

type InventoryTestSuite struct {
	suite.Suite
	repo   *mockProductRepository
	server *InventoryGRPCServer
}

func (s *InventoryTestSuite) SetupTest() {
	s.repo = new(mockProductRepository)
	svc := service.NewInventoryService(s.repo)
	s.server = NewInventoryGRPCServer(svc)
}

func (s *InventoryTestSuite) TestProductStockManagement() {
	ctx := context.Background()
	productID := uuid.New()

	// Setup mock for Create
	s.repo.On("Create", mock.AnythingOfType("*entity.Product")).Return(nil).Run(func(args mock.Arguments) {
		product := args.Get(0).(*entity.Product)
		product.ID = productID
	})

	// First create a test product
	createReq := &pb.CreateProductRequest{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		Stock:       100,
		CategoryId:  "test-category",
	}

	createResp, err := s.server.CreateProduct(ctx, createReq)
	s.NoError(err)
	s.NotNil(createResp)
	s.Equal(createReq.Name, createResp.Product.Name)
	s.Equal(int32(100), createResp.Product.Stock)

	// Setup mock for CheckStock
	s.repo.On("CheckStock", productID, 50).Return(true, nil)

	// Test CheckStock
	checkStockReq := &pb.CheckStockRequest{
		ProductId: productID.String(),
		Quantity:  50,
	}

	checkResp, err := s.server.CheckStock(ctx, checkStockReq)
	s.NoError(err)
	s.True(checkResp.Available)

	// Setup mock for UpdateStock
	s.repo.On("UpdateStock", productID, -30, mock.AnythingOfType("string")).Return(nil)

	// Test UpdateStock
	updateStockReq := &pb.UpdateStockRequest{
		ProductId: productID.String(),
		Quantity:  -30,
		OrderId:   uuid.New().String(),
	}

	updateResp, err := s.server.UpdateStock(ctx, updateStockReq)
	s.NoError(err)
	s.True(updateResp.Success)

	// Setup mock for Get
	s.repo.On("Get", productID).Return(&entity.Product{
		ID:          productID,
		Name:        createReq.Name,
		Description: createReq.Description,
		Price:       float64(createReq.Price),
		Stock:       70,
		CategoryID:  createReq.CategoryId,
	}, nil)

	// Verify the stock was updated correctly
	getReq := &pb.GetProductRequest{
		Id: productID.String(),
	}
	getResp, err := s.server.GetProduct(ctx, getReq)
	s.NoError(err)
	s.Equal(int32(70), getResp.Product.Stock)

	// Verify all mocked calls were made
	s.repo.AssertExpectations(s.T())
}

func (s *InventoryTestSuite) TestProductCRUD() {
	ctx := context.Background()
	productID := uuid.New()

	// Setup mock for Create
	s.repo.On("Create", mock.AnythingOfType("*entity.Product")).Return(nil).Run(func(args mock.Arguments) {
		product := args.Get(0).(*entity.Product)
		product.ID = productID
	})

	// Test Create
	createReq := &pb.CreateProductRequest{
		Name:        "CRUD Test Product",
		Description: "CRUD Test Description",
		Price:       149.99,
		Stock:       50,
		CategoryId:  "test-category",
	}

	createResp, err := s.server.CreateProduct(ctx, createReq)
	s.NoError(err)
	s.NotNil(createResp)
	s.Equal(createReq.Name, createResp.Product.Name)

	// Setup mock for Get
	s.repo.On("Get", productID).Return(&entity.Product{
		ID:          productID,
		Name:        createReq.Name,
		Description: createReq.Description,
		Price:       float64(createReq.Price),
		Stock:       50,
		CategoryID:  createReq.CategoryId,
	}, nil)

	// Test Get
	getReq := &pb.GetProductRequest{
		Id: productID.String(),
	}
	getResp, err := s.server.GetProduct(ctx, getReq)
	s.NoError(err)
	s.NotNil(getResp)
	s.Equal(createReq.Name, getResp.Product.Name)
	s.Equal(createReq.Description, getResp.Product.Description)
	s.Equal(float32(149.99), getResp.Product.Price)
	s.Equal(int32(50), getResp.Product.Stock)

	// Setup mock for Update
	s.repo.On("Update", mock.AnythingOfType("*entity.Product")).Return(nil)

	// Test Update
	updateReq := &pb.UpdateProductRequest{
		Id:          productID.String(),
		Name:        "Updated CRUD Test Product",
		Description: "Updated CRUD Test Description",
		Price:       199.99,
		Stock:       75,
		CategoryId:  "test-category",
	}

	updateResp, err := s.server.UpdateProduct(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updateResp)
	s.Equal(updateReq.Name, updateResp.Product.Name)
	s.Equal(updateReq.Description, updateResp.Product.Description)
	s.Equal(float32(199.99), updateResp.Product.Price)
	s.Equal(int32(75), updateResp.Product.Stock)

	// Setup mock for Delete
	s.repo.On("Delete", productID).Return(nil)

	// Test Delete
	deleteReq := &pb.DeleteProductRequest{
		Id: productID.String(),
	}
	deleteResp, err := s.server.DeleteProduct(ctx, deleteReq)
	s.NoError(err)
	s.True(deleteResp.Success)

	// Setup mock for Get after delete
	s.repo.On("Get", productID).Return(nil, ErrProductNotFound)

	// Verify deletion
	_, err = s.server.GetProduct(ctx, getReq)
	s.Error(err)
	st, ok := status.FromError(err)
	s.True(ok)
	s.Equal(codes.NotFound, st.Code())
	s.Equal("product not found", st.Message())

	// Verify all mocked calls were made
	s.repo.AssertExpectations(s.T())
}

func TestInventoryService(t *testing.T) {
	suite.Run(t, new(InventoryTestSuite))
}
