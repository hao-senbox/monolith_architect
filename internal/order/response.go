package order

import (
	"modular_monolith/internal/shared/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderResponse struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	UserID          primitive.ObjectID `json:"user_id" bson:"user_id"`
	Type            string             `json:"type" bson:"type"`
	OrderCode       string             `json:"order_code" bson:"order_code"`
	OrderItems      []OrderItem        `json:"order_items" bson:"order_items"`
	TotalPrice      float64            `json:"total_price" bson:"total_price"`
	Status          OrderStatus        `json:"status" bson:"status"`
	Discount        *float64           `json:"discount" bson:"discount"`
	ShippingAddress ShippingAddress    `json:"shipping_address" bson:"shipping_address"`
	CustomerNote    *string            `json:"customer_note" bson:"customer_note"`
	Payment         *model.Payment     `json:"payment" bson:"payment"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
}
