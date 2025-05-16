package service

import (
	pbInventory "E-Commerce/inventory-service/proto"
	"E-Commerce/order-service/internal/entity"
	"E-Commerce/order-service/internal/repository"
	"E-Commerce/order-service/internal/utils"
	pbProducer "E-Commerce/producer-service/proto"
	pbUser "E-Commerce/user-service/proto"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type OrderService interface {
	CreateOrder(userID string, items []*entity.OrderItem) (*entity.Order, error)
	GetOrder(id uuid.UUID) (*entity.Order, error)
	UpdateOrderStatus(id uuid.UUID, status string) error
	ListOrders(userID string, page, pageSize int) ([]*entity.Order, int, error)
}

type orderService struct {
	repo            repository.OrderRepository
	inventoryClient pbInventory.InventoryServiceClient
	producerClient  pbProducer.ProducerServiceClient
	userClient      pbUser.UserServiceClient
	emailConfig     utils.EmailConfig
}

func NewOrderService(repo repository.OrderRepository, inventoryClient pbInventory.InventoryServiceClient, emailConfig utils.EmailConfig) OrderService {
	// Connect to producer service
	producerConn, err := grpc.Dial("localhost:50054", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to producer-service: %v", err)
	}
	producerClient := pbProducer.NewProducerServiceClient(producerConn)

	// Connect to user service
	userConn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to user-service: %v", err)
	}
	userClient := pbUser.NewUserServiceClient(userConn)

	return &orderService{
		repo:            repo,
		inventoryClient: inventoryClient,
		producerClient:  producerClient,
		userClient:      userClient,
		emailConfig:     emailConfig,
	}
}

func (s *orderService) CreateOrder(userID string, items []*entity.OrderItem) (*entity.Order, error) {
	// Get user details for email
	userResp, err := s.userClient.GetUserProfile(context.Background(), &pbUser.GetUserProfileRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user details: %v", err)
	}
	userEmail := userResp.Email

	// Validate and check stock
	for _, item := range items {
		pid := item.ProductID
		p, err := s.inventoryClient.GetProduct(context.Background(), &pbInventory.GetProductRequest{Id: pid.String()})
		if err != nil {
			return nil, fmt.Errorf("failed to get product: %v", err)
		}
		item.Price = float64(p.Product.Price)

		resp, err := s.inventoryClient.CheckStock(context.Background(), &pbInventory.CheckStockRequest{
			ProductId: pid.String(),
			Quantity:  int32(item.Quantity),
		})
		if err != nil || !resp.Available {
			return nil, fmt.Errorf("insufficient stock for product %s", pid)
		}
	}

	// Update inventory
	for _, item := range items {
		_, err := s.inventoryClient.UpdateStock(context.Background(), &pbInventory.UpdateStockRequest{
			ProductId: item.ProductID.String(),
			Quantity:  -int32(item.Quantity),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to update stock: %v", err)
		}
	}

	// Calculate total
	total := 0.0
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}

	order := &entity.Order{
		ID:          uuid.New(),
		UserID:      userID,
		Status:      "pending",
		TotalAmount: total,
		CreatedAt:   time.Now(),
		Items:       items,
	}

	// Create order in database
	if err := s.repo.Create(order); err != nil {
		return nil, fmt.Errorf("failed to create order: %v", err)
	}

	// Notify producer service
	productIDs := make([]string, len(items))
	for i, item := range items {
		productIDs[i] = item.ProductID.String()
	}
	_, err = s.producerClient.NotifyOrderCreated(context.Background(), &pbProducer.OrderCreatedRequest{
		OrderId:    order.ID.String(),
		ProductIds: productIDs,
	})
	if err != nil {
		log.Printf("Failed to notify producer-service: %v", err)
	}

	// Prepare and send order confirmation email
	emailItems := make([]utils.OrderItemData, len(items))
	for i, item := range items {
		p, err := s.inventoryClient.GetProduct(context.Background(), &pbInventory.GetProductRequest{Id: item.ProductID.String()})
		if err != nil {
			log.Printf("Failed to get product details for email: %v", err)
			continue
		}
		emailItems[i] = utils.OrderItemData{
			ProductName: p.Product.Name,
			Quantity:    int(item.Quantity),
			Price:       float64(item.Price),
			Subtotal:    float64(item.Price) * float64(item.Quantity),
		}
	}

	emailData := utils.OrderEmailData{
		OrderID:     order.ID.String(),
		UserEmail:   userEmail,
		Items:       emailItems,
		TotalAmount: total,
		OrderStatus: order.Status,
	}

	if err := utils.SendOrderConfirmationEmail(s.emailConfig, emailData); err != nil {
		log.Printf("Failed to send order confirmation email: %v", err)
	}

	return order, nil
}

func (s *orderService) GetOrder(id uuid.UUID) (*entity.Order, error) {
	return s.repo.Get(id)
}

func (s *orderService) UpdateOrderStatus(id uuid.UUID, status string) error {
	return s.repo.UpdateStatus(id, status)
}

func (s *orderService) ListOrders(userID string, page, pageSize int) ([]*entity.Order, int, error) {
	return s.repo.List(userID, page, pageSize)
}
