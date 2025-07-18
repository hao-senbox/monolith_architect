package product

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CategoryID         primitive.ObjectID `json:"category_id" bson:"category_id"`
	ProductName        string             `json:"product_name" bson:"product_name"`
	ProductDescription string             `json:"product_description" bson:"product_description"`
	RatingAverage      float64            `json:"rating_average" bson:"rating_average"`
	ReviewsCount       int                `json:"reviews_count" bson:"reviews_count"`
	Color              string             `json:"color" bson:"color"`
	MainImagePublicID  string             `json:"main_image_public_id" bson:"main_image_public_id"`
	MainImage          string             `json:"main_image" bson:"main_image"`
	SubImages          []SubImage         `json:"sub_image" bson:"sub_image"`
	Price              float64            `json:"price" bson:"price"`
	Discount           float64            `json:"discount" bson:"discount"`
	Currency           string             `json:"currency" bson:"currency"`
	Sizes              []SizeOptions      `json:"sizes" bson:"sizes"`
	CreatedAt          time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at" bson:"updated_at"`
}

type SizeOptions struct {
	Size  string `json:"size" bson:"size"`
	Stock int    `json:"stock" bson:"stock"`
}

type SubImage struct {
	Url              string `json:"url" bson:"url"`
	SubImagePublicID string `json:"sub_image_public_id" bson:"sub_image_public_id"`
}
