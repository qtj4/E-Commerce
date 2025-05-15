package handler

import (
	pbInventory "E-Commerce/inventory-service/proto"
	pbOrder "E-Commerce/order-service/proto"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Mock inventory service client
type mockInventoryServiceClient struct {
	mock.Mock
}

func (m *mockInventoryServiceClient) CreateProduct(ctx context.Context, req *pbInventory.CreateProductRequest, opts ...grpc.CallOption) (*pbInventory.CreateProductResponse, error) {
	args := m.Called(ctx, req)
	if resp := args.Get(0); resp != nil {
		return resp.(*pbInventory.CreateProductResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockInventoryServiceClient) GetProduct(ctx context.Context, req *pbInventory.GetProductRequest, opts ...grpc.CallOption) (*pbInventory.GetProductResponse, error) {
	args := m.Called(ctx, req)
	if resp := args.Get(0); resp != nil {
		return resp.(*pbInventory.GetProductResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockInventoryServiceClient) UpdateProduct(ctx context.Context, req *pbInventory.UpdateProductRequest, opts ...grpc.CallOption) (*pbInventory.UpdateProductResponse, error) {
	args := m.Called(ctx, req)
	if resp := args.Get(0); resp != nil {
		return resp.(*pbInventory.UpdateProductResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockInventoryServiceClient) DeleteProduct(ctx context.Context, req *pbInventory.DeleteProductRequest, opts ...grpc.CallOption) (*pbInventory.DeleteProductResponse, error) {
	args := m.Called(ctx, req)
	if resp := args.Get(0); resp != nil {
		return resp.(*pbInventory.DeleteProductResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockInventoryServiceClient) ListProducts(ctx context.Context, req *pbInventory.ListProductsRequest, opts ...grpc.CallOption) (*pbInventory.ListProductsResponse, error) {
	args := m.Called(ctx, req)
	if resp := args.Get(0); resp != nil {
		return resp.(*pbInventory.ListProductsResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockInventoryServiceClient) CheckStock(ctx context.Context, req *pbInventory.CheckStockRequest, opts ...grpc.CallOption) (*pbInventory.CheckStockResponse, error) {
	args := m.Called(ctx, req)
	if resp := args.Get(0); resp != nil {
		return resp.(*pbInventory.CheckStockResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockInventoryServiceClient) UpdateStock(ctx context.Context, req *pbInventory.UpdateStockRequest, opts ...grpc.CallOption) (*pbInventory.UpdateStockResponse, error) {
	args := m.Called(ctx, req)
	if resp := args.Get(0); resp != nil {
		return resp.(*pbInventory.UpdateStockResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

// Mock order service client
type mockOrderServiceClient struct {
	mock.Mock
}

func (m *mockOrderServiceClient) CreateOrder(ctx context.Context, req *pbOrder.CreateOrderRequest, opts ...grpc.CallOption) (*pbOrder.CreateOrderResponse, error) {
	args := m.Called(ctx, req)
	if resp := args.Get(0); resp != nil {
		return resp.(*pbOrder.CreateOrderResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockOrderServiceClient) GetOrder(ctx context.Context, req *pbOrder.GetOrderRequest, opts ...grpc.CallOption) (*pbOrder.GetOrderResponse, error) {
	args := m.Called(ctx, req)
	if resp := args.Get(0); resp != nil {
		return resp.(*pbOrder.GetOrderResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockOrderServiceClient) UpdateOrder(ctx context.Context, req *pbOrder.UpdateOrderRequest, opts ...grpc.CallOption) (*pbOrder.UpdateOrderResponse, error) {
	args := m.Called(ctx, req)
	if resp := args.Get(0); resp != nil {
		return resp.(*pbOrder.UpdateOrderResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockOrderServiceClient) ListOrders(ctx context.Context, req *pbOrder.ListOrdersRequest, opts ...grpc.CallOption) (*pbOrder.ListOrdersResponse, error) {
	args := m.Called(ctx, req)
	if resp := args.Get(0); resp != nil {
		return resp.(*pbOrder.ListOrdersResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

// Test setup helper
func setupTest() (*RESTHandler, *mockInventoryServiceClient, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	mockInventory := new(mockInventoryServiceClient)
	mockOrder := new(mockOrderServiceClient)
	handler := NewRESTHandler(mockInventory, mockOrder)
	router := gin.New()
	return handler, mockInventory, router
}

func TestCreateProduct(t *testing.T) {
	handler, mockInventory, router := setupTest()
	router.POST("/products", handler.CreateProduct)

	tests := []struct {
		name         string
		request      map[string]interface{}
		mockResponse *pbInventory.CreateProductResponse
		mockError    error
		expectedCode int
		expectedBody interface{}
	}{
		{
			name: "Success",
			request: map[string]interface{}{
				"name":        "Test Product",
				"description": "Test Description",
				"price":       10.99,
				"stock":       100,
				"category_id": "category123",
			},
			mockResponse: &pbInventory.CreateProductResponse{
				Product: &pbInventory.Product{
					Id:          "product123",
					Name:        "Test Product",
					Description: "Test Description",
					Price:       10.99,
					Stock:       100,
					CategoryId:  "category123",
				},
			},
			mockError:    nil,
			expectedCode: http.StatusCreated,
		},
		{
			name: "Invalid Request - Missing Required Fields",
			request: map[string]interface{}{
				"description": "Test Description",
			},
			mockResponse: nil,
			mockError:    nil,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup request
			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest("POST", "/products", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Setup mock expectation
			if tt.mockResponse != nil {
				mockInventory.On("CreateProduct", mock.Anything, mock.MatchedBy(func(req *pbInventory.CreateProductRequest) bool {
					return req.Name == tt.request["name"].(string)
				})).Return(tt.mockResponse, tt.mockError).Once()
			}

			// Perform request
			router.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedCode == http.StatusCreated {
				var response *pbInventory.Product
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockResponse.Product.Id, response.Id)
			}
		})
	}
}

func TestGetProduct(t *testing.T) {
	handler, mockInventory, router := setupTest()
	router.GET("/products/:id", handler.GetProduct)

	tests := []struct {
		name         string
		productID    string
		mockResponse *pbInventory.GetProductResponse
		mockError    error
		expectedCode int
	}{
		{
			name:      "Success",
			productID: "product123",
			mockResponse: &pbInventory.GetProductResponse{
				Product: &pbInventory.Product{
					Id:          "product123",
					Name:        "Test Product",
					Description: "Test Description",
					Price:       10.99,
					Stock:       100,
					CategoryId:  "category123",
				},
			},
			mockError:    nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Product Not Found",
			productID:    "nonexistent",
			mockResponse: nil,
			mockError:    status.Error(codes.NotFound, "product not found"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/products/"+tt.productID, nil)
			w := httptest.NewRecorder()

			mockInventory.On("GetProduct", mock.Anything, &pbInventory.GetProductRequest{
				Id: tt.productID,
			}).Return(tt.mockResponse, tt.mockError).Once()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedCode == http.StatusOK {
				var response *pbInventory.Product
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockResponse.Product.Id, response.Id)
			}
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	handler, mockInventory, router := setupTest()
	router.PATCH("/products/:id", handler.UpdateProduct)

	tests := []struct {
		name         string
		productID    string
		request      map[string]interface{}
		mockResponse *pbInventory.UpdateProductResponse
		mockError    error
		expectedCode int
	}{
		{
			name:      "Success",
			productID: "product123",
			request: map[string]interface{}{
				"name":  "Updated Product",
				"price": 15.99,
			},
			mockResponse: &pbInventory.UpdateProductResponse{
				Product: &pbInventory.Product{
					Id:          "product123",
					Name:        "Updated Product",
					Description: "Test Description",
					Price:       15.99,
					Stock:       100,
					CategoryId:  "category123",
				},
			},
			mockError:    nil,
			expectedCode: http.StatusOK,
		},
		{
			name:      "Product Not Found",
			productID: "nonexistent",
			request: map[string]interface{}{
				"name": "Updated Product",
			},
			mockResponse: nil,
			mockError:    status.Error(codes.NotFound, "product not found"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest("PATCH", "/products/"+tt.productID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			if tt.mockResponse != nil || tt.mockError != nil {
				mockInventory.On("UpdateProduct", mock.Anything, mock.MatchedBy(func(req *pbInventory.UpdateProductRequest) bool {
					return req.Id == tt.productID
				})).Return(tt.mockResponse, tt.mockError).Once()
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedCode == http.StatusOK {
				var response *pbInventory.Product
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockResponse.Product.Id, response.Id)
				assert.Equal(t, tt.mockResponse.Product.Name, response.Name)
				assert.Equal(t, tt.mockResponse.Product.Price, response.Price)
			}
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	handler, mockInventory, router := setupTest()
	router.DELETE("/products/:id", handler.DeleteProduct)

	tests := []struct {
		name         string
		productID    string
		mockResponse *pbInventory.DeleteProductResponse
		mockError    error
		expectedCode int
	}{
		{
			name:      "Success",
			productID: "product123",
			mockResponse: &pbInventory.DeleteProductResponse{
				Success: true,
			},
			mockError:    nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Product Not Found",
			productID:    "nonexistent",
			mockResponse: nil,
			mockError:    status.Error(codes.NotFound, "product not found"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/products/"+tt.productID, nil)
			w := httptest.NewRecorder()

			mockInventory.On("DeleteProduct", mock.Anything, &pbInventory.DeleteProductRequest{
				Id: tt.productID,
			}).Return(tt.mockResponse, tt.mockError).Once()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedCode == http.StatusOK {
				var response map[string]bool
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.True(t, response["success"])
			}
		})
	}
}
