package order

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"modular_monolith/internal/cart"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req *CreateOrderRequest) (string, error)
	GetAllOrders(ctx context.Context) ([]Order, error)
	GetOrderByID(ctx context.Context, id string) (*Order, error)
	UpdateOrder(ctx context.Context, req *UpdateOrderRequest, id string) error
	DeleteOrder(ctx context.Context, id string) error
}

type orderService struct {
	orderRepo   OrderRepository
	cartService cart.CartService
}

func NewOrderService(orderRepo OrderRepository, cartService cart.CartService) OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		cartService: cartService,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, req *CreateOrderRequest) (string, error) {

	if req.UserID == "" {
		return "", fmt.Errorf("user_id is required")
	}

	if req.Address == "" {
		return "", fmt.Errorf("address is required")
	}

	if req.Email == "" {
		return "", fmt.Errorf("email is required")
	}

	if req.Name == "" {
		return "", fmt.Errorf("name is required")
	}

	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return "", fmt.Errorf("invalid user_id: %v", err)
	}

	carts, err := s.cartService.GetCartByUserID(ctx, req.UserID)
	if err != nil {
		return "", err
	}

	var orderItems []OrderItem
	for _, cart := range carts.CartItems {
		orderItem := &OrderItem{
			ProductID:    cart.ProductID,
			ProductName:  cart.ProductName,
			Quantity:     cart.Quantity,
			Price:        cart.Price,
			TotalPrice:   cart.TotalPrice,
			ProductImage: cart.ImageUrl,
			Size:         cart.Size,
		}
		orderItems = append(orderItems, *orderItem)
	}

	order := &Order{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		OrderCode: s.generateOrderCode(),
		ShippingAddress: ShippingAddress{
			Name:    req.Name,
			Email:   req.Email,
			Phone:   req.Phone,
			Address: req.Address,
		},
		Status: Pending,
		TotalPrice: carts.TotalPrice,
		OrderItems: orderItems,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	id, err := s.orderRepo.Create(ctx, order)
	if err != nil {
		return "", err
	}

	err = s.cartService.DeleteCart(ctx, req.UserID)
	if err != nil {
		return "", err
	}

	return id, nil

}

func (s *orderService) generateOrderCode() string {

	timestamp := time.Now().Format("20060102-150405")

	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	randomStr := hex.EncodeToString(b)

	return fmt.Sprintf("ORD-%s-%s", timestamp, randomStr)

}

func (s *orderService) GetAllOrders(ctx context.Context) ([]Order, error) {
	return s.orderRepo.FindAll(ctx)
}

func (s *orderService) GetOrderByID(ctx context.Context, id string) (*Order, error) {

	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %v", err)
	}

	return s.orderRepo.FindByID(ctx, objectID)
}

func (s *orderService) UpdateOrder(ctx context.Context, req *UpdateOrderRequest, id string) error {

	if id == "" {
		return fmt.Errorf("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	return s.orderRepo.UpdateByID(ctx, objectID, req.Status)
	
}

func (s *orderService) DeleteOrder(ctx context.Context, id string) error {

	if id == "" {
		return fmt.Errorf("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	return s.orderRepo.DeleteByID(ctx, objectID)
	
}