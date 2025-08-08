package payment

type CreatePaymentIntentRequest struct {
	OrderID string `json:"order_id"`
}

type VNPayRequest struct {
	OrderID     string  `json:"order_id"`
	OrderInfo   string  `json:"order_info"`
	ReturnURL   string  `json:"return_url"`
	BankCode    string  `json:"bank_code,omitempty"`
	Locale      string  `json:"locale,omitempty"`
}