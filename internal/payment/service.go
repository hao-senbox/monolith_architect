package payment

import (
	"context"
	"fmt"
	"modular_monolith/internal/order"
	"os"
	"time"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentService interface {
	CreatePaymentIntent(ctx context.Context, req *CreatePaymentIntentRequest) (*PaymentIntentResponse, error)
}

type paymentService struct {
	orderRepository   order.OrderRepository
	paymentRepository PaymentRepository
}

func NewPaymentService(paymentRepository PaymentRepository, orderRepository order.OrderRepository) PaymentService {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	return &paymentService{
		paymentRepository: paymentRepository,
		orderRepository:   orderRepository,
	}
}

func (s *paymentService) CreatePaymentIntent(ctx context.Context, req *CreatePaymentIntentRequest) (*PaymentIntentResponse, error) {

	if req.OrderID == "" {
		return nil, fmt.Errorf("order_id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("invalid order id: %v", err)
	}

	existingOrder, _ := s.orderRepository.FindByID(ctx, objectID)
	if existingOrder == nil {
		return nil, fmt.Errorf("order not found")
	}

	existingPaymenr, _ := s.paymentRepository.FindByOrderID(ctx, objectID)
	if existingPaymenr != nil {
		return nil, fmt.Errorf("payment already exists")
	}

	amount := int64(existingOrder.TotalPrice * 100)
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String("usd"),
		Metadata: map[string]string{
			"order_id": objectID.Hex(),
		},
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, err
	}

	payment := &Payment{
		ID:                  primitive.NewObjectID(),
		OrderID:             objectID,
		StripePaymentID:     pi.ID,
		StripePaymentSecret: pi.ClientSecret,
		Amount:              existingOrder.TotalPrice,
		Currency:            "usd",
		Status:              Pending,
		PaymentMethod:       "stripe",
		CreatedAt:           time.Now(),
		UpdateAt:            time.Now(),
	}

	err = s.paymentRepository.Create(ctx, payment)
	if err != nil {
		return nil, err
	}

	paymentRes := &PaymentIntentResponse{
		PaymentIntentID: pi.ID,
		ClientSecret:    pi.ClientSecret,
		Amount:          int(pi.Amount),
	}

	return paymentRes, nil

}
