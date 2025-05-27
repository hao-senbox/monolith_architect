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
	ProductSize        []ProductSize      `json:"product_size" bson:"product_size"`
	Color              []string           `json:"color" bson:"color"`
	Price              []Price            `json:"price" bson:"price"`
	Attachment         []Attachment       `json:"image" bson:"image"`
	CreatedAt          time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at" bson:"updated_at"`
}

type ProductSize struct {
	Size     string `json:"size" bson:"size"`
	Quantity int    `json:"quantity" bson:"quantity"`
}

type Price struct {
	Base     string `json:"base" bson:"base"`
	Discount string `json:"discount" bson:"discount"`
	Currency string `json:"currency" bson:"currency"`
}

type Attachment struct {
	Url string `json:"url" bson:"url"`
	ImageMain bool `json:"image_main" bson:"image_main"`
}
