package payment

type PaymentHandler struct {
	PaymentService PaymentService
}

func NewPaymentHandler(paymentService PaymentService) *PaymentHandler {
	return &PaymentHandler{
		PaymentService: paymentService,
	}
}
