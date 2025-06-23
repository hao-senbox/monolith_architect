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
	PublicID  string             `json:"public_id" bson:"public_id"`
	Address   string             `json:"address" bson:"address"`
	Bio       *string            `json:"bio" bson:"bio"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

