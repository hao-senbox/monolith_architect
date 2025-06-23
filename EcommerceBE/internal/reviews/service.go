package reviews

import (
	"context"
	"fmt"
	"modular_monolith/internal/user"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewService interface {
	CreateReview(ctx context.Context, req *CreateReviewRequest) error
	GetAllReviews(ctx context.Context, productID string) (*ReviewsResponse, error)
	GetReviewByID(ctx context.Context, id string) (*ReviewResponse, error)
	UpdateReview(ctx context.Context, id string, req *UpdateReviewRequest) error
	DeleteReview(ctx context.Context, id string) error
}

type reviewService struct {
	reviewRepo ReviewRepository
	userRepo   user.UserRepository
}

func NewReviewService(reviewRepo ReviewRepository, userRepo user.UserRepository) ReviewService {
	return &reviewService{
		reviewRepo: reviewRepo,
		userRepo:   userRepo,
	}
}

func (r *reviewService) CreateReview(ctx context.Context, req *CreateReviewRequest) error {

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
		ID:        primitive.NewObjectID(),
		ProductID: objectProductID,
		UserID:    objectUserID,
		Rating:    req.Rating,
		Review:    req.Review,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return r.reviewRepo.Create(ctx, review)

}

func (r *reviewService) GetAllReviews(ctx context.Context, productID string) (*ReviewsResponse, error) {

	objectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return nil, err
	}

	reviewList, err := r.reviewRepo.FindAll(ctx, objectID)
	if err != nil {
		return nil, err
	}

	var reviews []*ReviewResponse
	totalRating := 0

	for _, v := range reviewList {
		user, err := r.userRepo.FindByID(ctx, v.UserID)
		if err != nil {
			return nil, err
		}

		reviews = append(reviews, &ReviewResponse{
			ID:        v.ID,
			ProductID: v.ProductID,
			Rating:    v.Rating,
			Review:    v.Review,
			UserInfo: UserInfo{
				ID:       user.ID,
				FullName: user.Profile.FullName,
				Avatar:   user.Profile.Avatar,
			},
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})

		totalRating += v.Rating
	}

	average := 0.0
	if len(reviewList) > 0 {
		average = float64(totalRating) / float64(len(reviewList))
	}

	reviewsRes := &ReviewsResponse{
		ReviewsResponse:   reviews,
		TotalReviewsCount: len(reviews),
		AvarageRating:     average,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	return reviewsRes, nil
}

func (r *reviewService) GetReviewByID(ctx context.Context, id string) (*ReviewResponse, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid review id: %v", err)
	}

	review, err := r.reviewRepo.FindByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	user, err := r.userRepo.FindByID(ctx, review.UserID)
	if err != nil {
		return nil, err
	}

	reviewRes := &ReviewResponse{
		ID:        review.ID,
		ProductID: review.ProductID,
		Rating:    review.Rating,
		Review:    review.Review,
		UserInfo: UserInfo{
			ID:       user.ID,
			FullName: user.Profile.FullName,
			Avatar:   user.Profile.Avatar,
		},
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
	}

	return reviewRes, nil
}

func (r *reviewService) UpdateReview(ctx context.Context, id string, req *UpdateReviewRequest) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid review id: %v", err)
	}

	return r.reviewRepo.UpdateByID(ctx, objectID, req)

}

func (r *reviewService) DeleteReview(ctx context.Context, id string) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid review id: %v", err)
	}

	return r.reviewRepo.DeleteByID(ctx, objectID)

}