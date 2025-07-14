package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"modular_monolith/internal/order"
	"os"
	"time"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/webhook"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentService interface {
	CreatePaymentIntent(ctx context.Context, req *CreatePaymentIntentRequest) (*PaymentIntentResponse, error)
	ConfirmPayment(ctx context.Context, paymentIntentID string) error
	HandleWebhook(ctx context.Context, payload []byte, signature string) error
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

	amount := int64(existingOrder.TotalPrice)
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

func (s *paymentService) ConfirmPayment(ctx context.Context, paymentIntentID string) error {

	payment, err := s.paymentRepository.FindByStripePaymentID(ctx, paymentIntentID)
	if err != nil {
		return err
	}

	pi, err := paymentintent.Get(paymentIntentID, nil)
	if err != nil {
		return err
	}

	var newStatus PaymentStatus
	switch pi.Status {
	case stripe.PaymentIntentStatusSucceeded:
		newStatus = Success
	case stripe.PaymentIntentStatusRequiresPaymentMethod:
		newStatus = Cancelled
	case stripe.PaymentIntentStatusRequiresAction:
		newStatus = Failed
	default:
		newStatus = Pending
	}

	err = s.paymentRepository.UpdateStatus(ctx, payment.ID, newStatus)
	if err != nil {
		return err
	}
	err = s.orderRepository.UpdateByID(ctx, payment.OrderID, string(newStatus))
	if err != nil {
		return err
	}

	return nil
}

func (s *paymentService) HandleWebhook(ctx context.Context, payload []byte, signature string) error {

	event, err := webhook.ConstructEvent(payload, signature, os.Getenv("STRIPE_WEBHOOK_SECRET"))
	if err != nil {
		return fmt.Errorf("webhook signature verification failed: %w", err)
	}
	fmt.Println("Received event:", event.Type)

	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			return fmt.Errorf("failed to parse webhook: %w", err)
		}

		return s.ConfirmPayment(ctx, paymentIntent.ID)

	case "payment_intent.payment_failed":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			return fmt.Errorf("failed to parse webhook: %w", err)
		}

		payment, err := s.paymentRepository.FindByStripePaymentID(ctx, paymentIntent.ID)
		if err != nil {
			return err
		}
		return s.paymentRepository.UpdateStatus(ctx, payment.ID, Failed)
	}

	return nil

}
