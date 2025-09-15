package reviews

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reviews struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	ProductID  primitive.ObjectID `json:"product_id" bson:"product_id"`
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id"`
	Rating     int                `json:"rating" bson:"rating"`
	Review     string             `json:"review" bson:"review"`
	LikeReview []LikeReview       `json:"like_review" bson:"like_review"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}

type LikeReview struct {
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
