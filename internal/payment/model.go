package payment

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentStatus string
type PaymentProvider string

const (
	Pending   PaymentStatus = "pending"
	Success   PaymentStatus = "success"
	Failed    PaymentStatus = "failed"
	Cancelled PaymentStatus = "cancelled"
)

const (
	Stripe PaymentProvider = "stripe"
	VNPay  PaymentProvider = "vnpay"
)

type Payment struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id"`
	OrderID             primitive.ObjectID `json:"order_id" bson:"order_id"`
	// Stripe Payment
	StripePaymentID     *string             `json:"stripe_payment_id" bson:"stripe_payment_id"`
	StripePaymentSecret *string             `json:"stripe_payment_secret" bson:"stripe_payment_secret"`
	// Strpie Payment

	// VNPay
	VNPayTransactionNo   *string            `json:"vn_pay_transaction_no" bson:"vn_pay_transaction_no"`
	VNPayTransactionRef  *string            `json:"vn_pay_transaction_ref" bson:"vn_pay_transaction_ref"`
	VNPayResponseCode    *string            `json:"vn_pay_response_code" bson:"vn_pay_response_code"`
	VNPayBankCode        *string            `json:"vn_pay_bank_code" bson:"vn_pay_bank_code"`
	VNPayTransactionInfo *string            `json:"vn_pay_transaction_info" bson:"vn_pay_transaction_info"`
	// VNPay
	Amount              float64            `json:"amount" bson:"amount"`
	Currency            string             `json:"currency" bson:"currency"`
	Status              PaymentStatus      `json:"status" bson:"status"`
	PaymentMethod       string             `json:"payment_method" bson:"payment_method"`
	ExpiredAt           time.Time          `json:"expired_at" bson:"expired_at"`
	CreatedAt           time.Time          `json:"created_at" bson:"created_at"`
	UpdateAt            time.Time          `json:"updated_at" bson:"updated_at"`
}