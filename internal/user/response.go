package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserWithProfile struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email     string             `bson:"email,omitempty" json:"email"`
	FirstName string             `bson:"first_name,omitempty" json:"first_name"`
	LastName  string             `bson:"last_name,omitempty" json:"last_name"`
	Phone     string             `bson:"phone,omitempty" json:"phone"`
	UserType  string             `bson:"user_type,omitempty" json:"user_type"`
	Profile   *Profile           `bson:"profile,omitempty" json:"profile,omitempty"`
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
