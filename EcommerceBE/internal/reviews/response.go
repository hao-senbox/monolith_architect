package reviews

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewsResponse struct {
	ReviewsResponse   []*ReviewResponse `json:"reviews" bson:"reviews"`
	TotalReviewsCount int               `json:"total_reviews_count" bson:"total_reviews_count"`
	AvarageRating     float64           `json:"average_rating" bson:"average_rating"`
	CreatedAt         time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at" bson:"updated_at"`
}

type ReviewResponse struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id"`
	Rating    int                `json:"rating" bson:"rating"`
	Review    string             `json:"review" bson:"review"`
	UserInfo  UserInfo           `json:"user_info" bson:"user_info"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type UserInfo struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	FullName  string             `json:"full_name" bson:"full_name"`
	Avatar    string             `json:"avatar" bson:"avatar"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
