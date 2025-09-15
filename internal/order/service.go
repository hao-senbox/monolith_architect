package order

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"modular_monolith/internal/cart"
	"modular_monolith/internal/coupon"
	"modular_monolith/internal/product"
	"modular_monolith/internal/shared/model"
	"modular_monolith/internal/shared/ports"
	"modular_monolith/pkg/email"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req *CreateOrderRequest) (string, error)
	GetAllOrders(ctx context.Context) ([]*OrderResponse, error)
	GetOrderByID(ctx context.Context, id string) (*OrderResponse, error)
	UpdateOrder(ctx context.Context, req *UpdateOrderRequest, id string) error
	DeleteOrder(ctx context.Context, id string) error
	GetOrderByUserID(ctx context.Context, userID string) ([]*OrderResponse, error)
}

type orderService struct {
	orderRepo         OrderRepository
	cartService       cart.CartService
	couponRepository  coupon.CouponRepository
	paymentRepository ports.PaymentRepository
	productRepository product.ProductRepository
	EmailService      *email.EmailService
}

func NewOrderService(orderRepo OrderRepository, cartService cart.CartService, couponRepository coupon.CouponRepository, paymentRepository ports.PaymentRepository, productRepository product.ProductRepository) OrderService {
	emailService := email.NewEmailService()
	return &orderService{
		orderRepo:         orderRepo,
		cartService:       cartService,
		couponRepository:  couponRepository,
		paymentRepository: paymentRepository,
		productRepository: productRepository,
		EmailService:      emailService,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, req *CreateOrderRequest) (string, error) {

	var orderData *Order

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

	if req.Type == "" {
		return "", fmt.Errorf("type is required")
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

		err = s.productRepository.UpdateQuantityByID(ctx, cart.ProductID, cart.Size, -cart.Quantity)
		if err != nil {
			return "", fmt.Errorf("product %s (size %s) is out of stock or insufficient", cart.ProductName, cart.Size)
		}

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

	if req.CouponCode != nil {
		coupon, err := s.couponRepository.FindByCode(ctx, *req.CouponCode)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return "", fmt.Errorf("invalid coupon code")
			}
			return "", err
		}
		priceDiscount := carts.TotalPrice - (carts.TotalPrice * coupon.Discount / 100)
		orderData = &Order{
			ID:        primitive.NewObjectID(),
			UserID:    userID,
			Type:      req.Type,
			OrderCode: s.generateOrderCode(),
			ShippingAddress: ShippingAddress{
				Name:    req.Name,
				Email:   req.Email,
				Phone:   req.Phone,
				Address: req.Address,
			},
			Status:     Pending,
			TotalPrice: priceDiscount,
			OrderItems: orderItems,
			Discount:   &coupon.Discount,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

	} else {
		orderData = &Order{
			ID:        primitive.NewObjectID(),
			UserID:    userID,
			Type:      req.Type,
			OrderCode: s.generateOrderCode(),
			ShippingAddress: ShippingAddress{
				Name:    req.Name,
				Email:   req.Email,
				Phone:   req.Phone,
				Address: req.Address,
			},
			Status:     Pending,
			TotalPrice: carts.TotalPrice,
			OrderItems: orderItems,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

	}

	id, err := s.orderRepo.Create(ctx, orderData)
	if err != nil {
		return "", err
	}

	if strings.EqualFold(req.Type, "cod") {
		html := BuildOrderEmailHTML(*orderData,
			"Football Shop",
		)
		_ = s.EmailService.SendEmail(req.Email, "Order successful #"+orderData.OrderCode, html)
	}

	if req.CouponCode != nil {
		err = s.couponRepository.AddUserIsUsed(ctx, userID, *req.CouponCode)
		if err != nil {
			return "", err
		}
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

func (s *orderService) GetAllOrders(ctx context.Context) ([]*OrderResponse, error) {
	var data []*OrderResponse

	orders, err := s.orderRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, order := range orders {

		var payment *model.Payment

		payment, err = s.paymentRepository.FindByOrderID(ctx, order.ID)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				payment = nil
			} else {
				return nil, err
			}
		}

		data = append(data, &OrderResponse{
			ID:        order.ID,
			UserID:    order.UserID,
			Type:      order.Type,
			OrderCode: order.OrderCode,
			ShippingAddress: ShippingAddress{
				Name:    order.ShippingAddress.Name,
				Email:   order.ShippingAddress.Email,
				Phone:   order.ShippingAddress.Phone,
				Address: order.ShippingAddress.Address,
			},
			Status:     order.Status,
			TotalPrice: order.TotalPrice,
			OrderItems: order.OrderItems,
			CreatedAt:  order.CreatedAt,
			UpdatedAt:  order.UpdatedAt,
			Payment:    payment,
		})
	}

	return data, nil
}

func (s *orderService) GetOrderByID(ctx context.Context, id string) (*OrderResponse, error) {

	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %v", err)
	}

	order, err := s.orderRepo.FindByID(ctx, objectID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	var payment *model.Payment

	payment, err = s.paymentRepository.FindByOrderID(ctx, order.ID)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	data := &OrderResponse{
		ID:        order.ID,
		UserID:    order.UserID,
		Type:      order.Type,
		OrderCode: order.OrderCode,
		ShippingAddress: ShippingAddress{
			Name:    order.ShippingAddress.Name,
			Email:   order.ShippingAddress.Email,
			Phone:   order.ShippingAddress.Phone,
			Address: order.ShippingAddress.Address,
		},
		Status:     order.Status,
		TotalPrice: order.TotalPrice,
		OrderItems: order.OrderItems,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
		Payment:    payment,
	}

	return data, nil
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

func (s *orderService) GetOrderByUserID(ctx context.Context, userID string) ([]*OrderResponse, error) {

	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %v", err)
	}

	orders, err := s.orderRepo.FindByUserID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	var data []*OrderResponse

	for _, order := range orders {

		var payment *model.Payment

		payment, err = s.paymentRepository.FindByOrderID(ctx, order.ID)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				payment = nil
			} else {
				return nil, err
			}
		}

		data = append(data, &OrderResponse{
			ID:        order.ID,
			UserID:    order.UserID,
			Type:      order.Type,
			OrderCode: order.OrderCode,
			ShippingAddress: ShippingAddress{
				Name:    order.ShippingAddress.Name,
				Email:   order.ShippingAddress.Email,
				Phone:   order.ShippingAddress.Phone,
				Address: order.ShippingAddress.Address,
			},
			Status:     order.Status,
			TotalPrice: order.TotalPrice,
			OrderItems: order.OrderItems,
			CreatedAt:  order.CreatedAt,
			UpdatedAt:  order.UpdatedAt,
			Payment:    payment,
		})
	}

	return data, nil
}
