package product

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductResponse struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CategoryID         primitive.ObjectID `json:"category_id" bson:"category_id"`
	ProductName        string             `json:"product_name" bson:"product_name"`
	ProductDescription string             `json:"product_description" bson:"product_description"`
	RatingAverage      float64            `json:"rating_average" bson:"rating_average"`
	ReviewsCount       int                `json:"reviews_count" bson:"reviews_count"`
	Variants           []ProductVariant   `json:"variants" bson:"variants"`
	CreatedAt          time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at" bson:"updated_at"`
}