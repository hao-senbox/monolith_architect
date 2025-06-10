package category

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID           primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	CategoryName string              `json:"category_name" bson:"category_name"`
	ParentID     *primitive.ObjectID `json:"parent_id" bson:"parent_id"`
	CreatedAt    time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at" bson:"updated_at"`
}


