package order

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderStatus string

const (
	Pending    OrderStatus = "pending"
	Paid       OrderStatus = "paid"
	Processing OrderStatus = "processing"
	Cancelled  OrderStatus = "cancelled"
)

type Order struct {
	ID              primitive.ObjectID  `json:"id" bson:"_id"`
	UserID          primitive.ObjectID  `json:"user_id" bson:"user_id"`
	OrderCode       string              `json:"order_code" bson:"order_code"`
	OrderItems      []OrderItem         `json:"order_items" bson:"order_items"`
	TotalPrice      float64             `json:"total_price" bson:"total_price"`
	Status          OrderStatus         `json:"status" bson:"status"`
	Discount        *float64            `json:"discount" bson:"discount"`
	ShippingAddress ShippingAddress     `json:"shipping_address" bson:"shipping_address"`
	CustomerNote    *string             `json:"customer_note" bson:"customer_note"`
	CreatedAt       time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at" bson:"updated_at"`
}

type OrderItem struct {
	ProductID    primitive.ObjectID `json:"product_id" bson:"product_id"`
	ProductName  string             `json:"product_name" bson:"product_name"`
	ProductImage string             `json:"product_image" bson:"product_image"`
	Quantity     int                `json:"quantity" bson:"quantity"`
	Price        float64            `json:"price" bson:"price"`
	Size         string             `json:"size" bson:"size"`
	TotalPrice   float64            `json:"total_price" bson:"total_price"`
}

type ShippingAddress struct {
	Name    string `json:"name" bson:"name"`
	Email   string `json:"email" bson:"email"`
	Phone   string `json:"phone" bson:"phone"`
	Address string `json:"address" bson:"address"`
}
