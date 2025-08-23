package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentStatus string
type PaymentProvider string

const (
	PaymentPending   PaymentStatus = "pending"
	PaymentSuccess   PaymentStatus = "success"
	PaymentFailed    PaymentStatus = "failed"
	PaymentCancelled PaymentStatus = "cancelled"
)

const (
	ProviderStripe PaymentProvider = "stripe"
	ProviderVNPay  PaymentProvider = "vnpay"
)

type Payment struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id"`
	OrderID             primitive.ObjectID `json:"order_id" bson:"order_id"`

	// Stripe
	StripePaymentID     *string `json:"stripe_payment_id" bson:"stripe_payment_id"`
	StripePaymentSecret *string `json:"stripe_payment_secret" bson:"stripe_payment_secret"`

	// VNPay
	VNPayTransactionNo   *string `json:"vn_pay_transaction_no" bson:"vn_pay_transaction_no"`
	VNPayTransactionRef  *string `json:"vn_pay_transaction_ref" bson:"vn_pay_transaction_ref"`
	VNPayResponseCode    *string `json:"vn_pay_response_code" bson:"vn_pay_response_code"`
	VNPayBankCode        *string `json:"vn_pay_bank_code" bson:"vn_pay_bank_code"`
	VNPayTransactionInfo *string `json:"vn_pay_transaction_info" bson:"vn_pay_transaction_info"`

	Amount        float64       `json:"amount" bson:"amount"`
	Currency      string        `json:"currency" bson:"currency"`
	Status        PaymentStatus `json:"status" bson:"status"`
	PaymentMethod string        `json:"payment_method" bson:"payment_method"`
	CreatedAt     time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at" bson:"updated_at"`
}
