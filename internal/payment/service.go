package payment

type PaymentService interface{

} 

type paymentService struct{
	paymentRepository PaymentRepository	
}

func NewPaymentService(paymentRepository PaymentRepository) PaymentService {
	return &paymentService{
		paymentRepository: paymentRepository,
	}
}