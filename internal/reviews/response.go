package reviews

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewResponse struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id"`
	UserInfo  UserInfo           `json:"user_info" bson:"user_info"`
	Rating    int                `json:"rating" bson:"rating"`
	Review    string             `json:"review" bson:"review"`
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