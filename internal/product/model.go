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
	Variants           []ProductVariant   `json:"variants" bson:"variants"`
	CreatedAt          time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at" bson:"updated_at"`
}

type ProductVariant struct {
	SKU       string     `json:"sku" bson:"sku"`
	Color     string     `json:"color" bson:"color"`
	Size      string     `json:"size" bson:"size"`
	Stock     int        `json:"stock" bson:"stock"`
	Price     float64    `json:"price" bson:"price"`
	Discount  float64    `json:"discount" bson:"discount"`
	Currency  string     `json:"currency" bson:"currency"`
	MainImage string     `json:"main_image" bson:"main_image"`
	SubImage  []SubImage `json:"sub_image" bson:"sub_image"`
}

type SubImage struct {
	Url  string `json:"url" bson:"url"`
	Sort int    `json:"sort" bson:"sort"`
}

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

type CreateProductRequest struct {
	CategoryID         primitive.ObjectID `json:"category_id" bson:"category_id"`
	ProductName        string             `json:"product_name" bson:"product_name"`
	ProductDescription string             `json:"product_description" bson:"product_description"`
	Variants           []ProductVariant   `json:"variants" bson:"variants"`
}

type UpdateProductRequest struct {
	CategoryID         primitive.ObjectID `json:"category_id" bson:"category_id"`
	ProductName        string             `json:"product_name" bson:"product_name"`
	ProductDescription string             `json:"product_description" bson:"product_description"`
	Variants           []ProductVariant   `json:"variants" bson:"variants"`
}