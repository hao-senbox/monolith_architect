package profile

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

type CreateProfileRequest struct {
	UserID   string `form:"user_id" bson:"user_id"`
	FullName string `form:"full_name" bson:"full_name"`
	Gender   string `form:"gender" bson:"gender"`
	BirthDay string `form:"birth_day" bson:"birth_day"`
	Address  string `form:"address" bson:"address"`
	Bio      string `form:"bio" bson:"bio"`
}

type UpdateProfileRequest struct {
	UserID   primitive.ObjectID `json:"user_id" bson:"user_id"`
	FullName string             `json:"full_name" bson:"full_name"`
	Gender   string             `json:"gender" bson:"gender"`
	BirthDay string             `json:"birth_day" bson:"birth_day"`
	Avatar   string             `json:"avatar" bson:"avatar"`
	Address  string             `json:"address" bson:"address"`
	Bio      string             `json:"bio" bson:"bio"`
}
