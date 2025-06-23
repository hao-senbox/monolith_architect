package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	LastName     string             `json:"last_name" bson:"last_name"`
	FristName    string             `json:"first_name" bson:"first_name"`
	Email        string             `json:"email" bson:"email"`
	Password     string             `json:"password" bson:"password"`
	Phone        string             `json:"phone" bson:"phone"`
	Token        string             `json:"token" bson:"token"`
	RefreshToken string             `json:"refresh_token" bson:"refresh_token"`
	UserType     string             `json:"user_type" bson:"user_type"`
	CreatedAt    string             `json:"created_at" bson:"created_at"`
	UpdatedAt    string             `json:"updated_at" bson:"updated_at"`
}
