package payment

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentStatus string

const (
	Pending   PaymentStatus = "pending"
	Success   PaymentStatus = "success"
	Failed    PaymentStatus = "failed"
	Cancelled PaymentStatus = "cancelled"
)
type Payment struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id"`
	OrderID             primitive.ObjectID `json:"order_id" bson:"order_id"`
	StripePaymentID     string             `json:"stripe_payment_id" bson:"stripe_payment_id"`
	StripePaymentSecret string             `json:"stripe_payment_secret" bson:"stripe_payment_secret"`
	Amount              float64            `json:"amount" bson:"amount"`
	Currency            string             `json:"currency" bson:"currency"`
	Status              PaymentStatus      `json:"status" bson:"status"`
	PaymentMethod       string             `json:"payment_method" bson:"payment_method"`
	CreatedAt           time.Time          `json:"created_at" bson:"created_at"`
	UpdateAt            time.Time          `json:"updated_at" bson:"updated_at"`
}