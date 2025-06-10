package reviews

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewService interface {
	CreateReview(ctx *gin.Context, req *CreateReviewRequest) error
}

type reviewService struct {
	reviewRepo ReviewRepository
}

func NewReviewService(reviewRepo ReviewRepository) ReviewService {
	return &reviewService{reviewRepo: reviewRepo}
}

func (r *reviewService) CreateReview(ctx *gin.Context, req *CreateReviewRequest) error {

	if req.ProductID == "" {
		return fmt.Errorf("product_id is required")
	}

	if req.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	if req.Rating <= 0 {
		return fmt.Errorf("rating must be greater than 0")
	}

	if req.Review == "" {
		return fmt.Errorf("review is required")
	}

	objectUserID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return fmt.Errorf("invalid user id: %v", err)
	}

	objectProductID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		return fmt.Errorf("invalid product id: %v", err)
	}

	review := &Reviews{
		ID        : primitive.NewObjectID(),
		ProductID : objectProductID,
		UserID    : objectUserID,
		Rating    : req.Rating,
		Review    : req.Review,
		CreatedAt : time.Now(),
		UpdatedAt : time.Now(),
	}
	
	return r.reviewRepo.Create(ctx, review)
	
}