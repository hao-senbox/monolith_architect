package cart

import (
	"fmt"
	"modular_monolith/internal/product"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartService interface {
	CreateCart(ctx *gin.Context, req *AddtoCartRequest) error
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

func (s *cartService) CreateCart(c *gin.Context, req *AddtoCartRequest) error {

	productID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		return fmt.Errorf("invalid product id: %v", err)
	}

	product, err := s.productRepo.FindByID(c, productID)
	if err != nil || product == nil {
		return fmt.Errorf("product not found")
	}

	productItem := &Product{
		ProductID:   product.ID,
		ProductName: product.ProductName,
		Quantity:    req.Quantity,
		TotalPrice:  product.Price * float64(req.Quantity),
		Price:       product.Price,
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

	return s.repo.AddToCart(c, productItem, userID, quantity)

}
