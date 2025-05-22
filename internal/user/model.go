package user

import (
	"time"

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

type UserWithProfile struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email    string             `bson:"email,omitempty" json:"email"`
	Phone    string             `bson:"phone,omitempty" json:"phone"`
	UserType string             `bson:"user_type,omitempty" json:"user_type"`
	Profile  *Profile           `bson:"profile,omitempty" json:"profile,omitempty"`
}

type Profile struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	FullName  string             `json:"full_name" bson:"full_name"`
	Gender    string             `json:"gender" bson:"gender"`
	BirthDay  time.Time          `json:"birth_day" bson:"birth_day"`
	Avatar    string             `json:"avatar" bson:"avatar"`
	Address   string             `json:"address" bson:"address"`
	Bio       *string            `json:"bio" bson:"bio"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type RegisterRequest struct {
	FristName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
	Phone     string `json:"phone" bson:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}
