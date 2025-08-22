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

type VNPayCallbackRequest struct {
	VNPayTransactionNo   *string            `json:"vn_pay_transaction_no" bson:"vn_pay_transaction_no"`
	VNPayTransactionRef  *string            `json:"vn_pay_transaction_ref" bson:"vn_pay_transaction_ref"`
	VNPayResponseCode    *string            `json:"vn_pay_response_code" bson:"vn_pay_response_code"`
	VNPayBankCode        *string            `json:"vn_pay_bank_code" bson:"vn_pay_bank_code"`
	VNPayTransactionInfo *string            `json:"vn_pay_transaction_info" bson:"vn_pay_transaction_info"`
}