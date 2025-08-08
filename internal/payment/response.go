package payment

type PaymentIntentResponse struct {
	PaymentIntentID string `json:"payment_intent_id"`
	ClientSecret    string `json:"client_secret"`
	Amount          int    `json:"amount"`
}

type VNPayResponse struct {
	PaymentURL string `json:"payment_url"`
}

type VNPayCallback struct {
	Amount           string `form:"vnp_Amount"`
	BankCode         string `form:"vnp_BankCode"`
	BankTranNo       string `form:"vnp_BankTranNo"`
	CardType         string `form:"vnp_CardType"`
	OrderInfo        string `form:"vnp_OrderInfo"`
	PayDate          string `form:"vnp_PayDate"`
	ResponseCode     string `form:"vnp_ResponseCode"`
	TmnCode          string `form:"vnp_TmnCode"`
	TransactionNo    string `form:"vnp_TransactionNo"`
	TransactionRef   string `form:"vnp_TxnRef"`
	SecureHashType   string `form:"vnp_SecureHashType"`
	SecureHash       string `form:"vnp_SecureHash"`
}	
