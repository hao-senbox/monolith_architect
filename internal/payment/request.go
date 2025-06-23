package payment

type CreatePaymentIntentRequest struct {
	OrderID string `json:"order_id"`
}