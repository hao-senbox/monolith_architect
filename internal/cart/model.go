package cart

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cart struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id"`
	CartItems  []CartItem         `json:"cart_items" bson:"cart_items"`
	TotalPrice float64            `json:"total_price" bson:"total_price"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}

type CartItem struct {
	Products   []Product `json:"product" bson:"product"`
	TotalPrice float64   `json:"total_price" bson:"total_price"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" bson:"updated_at"`
}

type Product struct {
	ProductID  primitive.ObjectID `json:"product_id" bson:"product_id"`
	ProductName string             `json:"product_name" bson:"product_name"`
	Quantity   int                `json:"quantity" bson:"quantity"`
	Price      float64            `json:"price" bson:"price"`
	TotalPrice float64            `json:"total_price" bson:"total_price"`
	ImageUrl   string             `json:"image_url" bson:"image_url"`
}
