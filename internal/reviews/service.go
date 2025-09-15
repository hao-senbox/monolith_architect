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
	LikeReview(ctx context.Context, req *LikeReviewRequest, id string) error
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

	ratingCount := map[int]int{
		1: 0,
		2: 0,
		3: 0,
		4: 0,
		5: 0,
	}

	var percent []PercentRating
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

		totalRating += v.Rating

		ratingCount[v.Rating]++

		user, err := r.userRepo.FindByID(ctx, v.UserID)
		if err != nil {
			return nil, err
		}

		var avatar *string
		if user.Profile != nil {
			avatar = &user.Profile.Avatar
		}

		reviews = append(reviews, &ReviewResponse{
			ID:        v.ID,
			ProductID: v.ProductID,
			Rating:    v.Rating,
			Review:    v.Review,
			UserInfo: UserInfo{
				ID:       user.ID,
				FullName: user.LastName + user.FirstName,
				Avatar:   avatar,
			},
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}

	average := 0.0
	if len(reviewList) > 0 {
		average = float64(totalRating) / float64(len(reviewList))
	}

	for i := 1; i <= 5; i++ {
		count := ratingCount[i]
		percentValue := 0
		if len(reviewList) > 0 {
			percentValue = int(float64(count) / float64(len(reviewList)) * 100)
		}
		percent = append(percent, PercentRating{
			Rating:  i,
			Count:   count,
			Percent: fmt.Sprintf("%d%%", percentValue),
		})
	}

	reviewsRes := &ReviewsResponse{
		ReviewsResponse:   reviews,
		TotalReviewsCount: len(reviews),
		AvarageRating:     average,
		CreatedAt:         time.Now(),
		Percent:           percent,
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

	var avatar *string
	if user.Profile != nil {
		avatar = &user.Profile.Avatar
	}

	reviewRes := &ReviewResponse{
		ID:        review.ID,
		ProductID: review.ProductID,
		Rating:    review.Rating,
		Review:    review.Review,
		UserInfo: UserInfo{
			ID:       user.ID,
			FullName: user.LastName + user.FirstName,
			Avatar:   avatar,
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

func (r *reviewService) LikeReview(ctx context.Context, req *LikeReviewRequest, id string) error {

	if id == "" {
		return fmt.Errorf("id is required")
	}
	if req.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	if req.Type == "" {
		return fmt.Errorf("type is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid review id: %v", err)
	}

	userObjID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return fmt.Errorf("invalid user id: %v", err)
	}

	review, err := r.reviewRepo.FindByID(ctx, objectID)
	if err != nil {
		return err
	}

	if review == nil {
		return fmt.Errorf("review not found")
	}

	if req.Type == "like" {
		alreadyLiked := false
		for _, v := range review.LikeReview {
			if v.UserID == userObjID {
				alreadyLiked = true
				break
			}
		}
		if !alreadyLiked {
			review.LikeReview = append(review.LikeReview, LikeReview{
				UserID:    userObjID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		}
	} else if req.Type == "unlike" {
		for i, v := range review.LikeReview {
			if v.UserID == userObjID {
				review.LikeReview = append(review.LikeReview[:i], review.LikeReview[i+1:]...)
				break
			}
		}
	} else {
		return fmt.Errorf("invalid type, must be 'like' or 'unlike'")
	}

	updateReq := &UpdateReviewRequest{
		Rating:      review.Rating,
		Review:      review.Review,
		LikeReview:  review.LikeReview, 
	}
	return r.reviewRepo.UpdateByID(ctx, objectID, updateReq)
}
