package cart

import (
	"context"
	"fmt"
	"modular_monolith/internal/product"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartService interface {
	CreateCart(ctx context.Context, req *AddtoCartRequest) error
	GetCartByUserID(ctx context.Context, userID string) (*Cart, error)
	UpdateCart(ctx context.Context, req *UpdateCartRequest) error
	DeleteItemCart(ctx context.Context, req *DeleteItemCartRequest) error
	DeleteCart(ctx context.Context, userID string) error
}

type cartService struct {
	repo        CartRepository
	productRepo product.ProductRepository
}

func NewCartService(repo CartRepository, productRepo product.ProductRepository) CartService {
	return &cartService{
		repo:        repo,
		productRepo: productRepo,
	}
}

func (s *cartService) CreateCart(c context.Context, req *AddtoCartRequest) error {

	productID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		return fmt.Errorf("invalid product id: %v", err)
	}

	product, err := s.productRepo.FindByID(c, productID)
	if err != nil || product == nil {
		return fmt.Errorf("product not found")
	}

	cartItem := &CartItem{
		ProductID:   product.ID,
		ProductName: product.ProductName,
		Quantity:    req.Quantity,
		TotalPrice:  product.Price * float64(req.Quantity),
		Price:       product.Price,
		Size:        req.Size,
		ImageUrl:    product.MainImage,
	}

	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return fmt.Errorf("invalid user id: %v", err)
	}

	quantity := req.Quantity

	if quantity <= 0 {
		return fmt.Errorf("quantity must be greater than 0")
	}

	return s.repo.AddToCart(c, cartItem, userID)

}

func (s *cartService) GetCartByUserID(c context.Context, userID string) (*Cart, error) {

	if userID == "" {
		return nil, fmt.Errorf("user id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %v", err)
	}

	cart, err := s.repo.FindCartByUserID(c, objectID)
	if err != nil {
		return nil, err
	}

	return cart, nil

}

func (s *cartService) UpdateCart(c context.Context, req *UpdateCartRequest) error {

	if req.ProductID == "" {
		return fmt.Errorf("product id is required")
	}

	if req.UserID == "" {
		return fmt.Errorf("user id is required")
	}

	if req.Quantity <= 0 {
		return fmt.Errorf("quantity must be greater than 0")
	}

	if req.Types == "" {
		return fmt.Errorf("types is required")
	}

	objectProductID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		return fmt.Errorf("invalid product id: %v", err)
	}

	objectUserID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return fmt.Errorf("invalid user id: %v", err)
	}

	return s.repo.UpdateCart(c, objectProductID, objectUserID, req.Quantity, req.Types)

}

func (s *cartService) DeleteItemCart(c context.Context, req *DeleteItemCartRequest) error {

	if req.ProductID == "" {
		return fmt.Errorf("product id is required")
	}

	if req.UserID == "" {
		return fmt.Errorf("user id is required")
	}

	objectProductID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		return fmt.Errorf("invalid product id: %v", err)
	}

	objectUserID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return fmt.Errorf("invalid user id: %v", err)
	}

	return s.repo.DeleteItemCart(c, objectProductID, objectUserID)
}

func (s *cartService) DeleteCart(c context.Context, userID string) error {

	if userID == "" {
		return fmt.Errorf("user id is required")
	}

	objectUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user id: %v", err)
	}

	return s.repo.DeleteCart(c, objectUserID)
	
}