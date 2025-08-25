package payment

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"modular_monolith/config"
	"modular_monolith/internal/order"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
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
	CreateVNPayPayment(ctx context.Context, req *VNPayRequest, clientIP string) (string, error)
	RepurchaseOrder(ctx context.Context, req *VNPayRequest, clientIP string) (*RepurchaseOrderResponse, error)
	HandleVNPayCallback(ctx context.Context, callback *VNPayCallback) error
	CronPaymentExpiration(ctx context.Context) error
}

type paymentService struct {
	orderRepository   order.OrderRepository
	paymentRepository PaymentRepository
	config            config.VNPayConfig
}

func NewPaymentService(paymentRepository PaymentRepository, orderRepository order.OrderRepository, config config.VNPayConfig) PaymentService {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	return &paymentService{
		paymentRepository: paymentRepository,
		orderRepository:   orderRepository,
		config:            config,
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
		StripePaymentID:     &pi.ID,
		StripePaymentSecret: &pi.ClientSecret,
		Amount:              existingOrder.TotalPrice,
		Currency:            "usd",
		Status:              Pending,
		PaymentMethod:       "stripe",
		CreatedAt:           time.Now(),
		UpdateAt:            time.Now(),
	}

	_, err = s.paymentRepository.Create(ctx, payment)
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

func (s *paymentService) CreateVNPayPayment(ctx context.Context, req *VNPayRequest, clientIP string) (string, error) {

	if req.OrderID == "" {
		return "", fmt.Errorf("order id is required")
	}

	orderID, err := primitive.ObjectIDFromHex(req.OrderID)
	if err != nil {
		return "", fmt.Errorf("invalid order id: %v", err)
	}

	existingOrder, _ := s.orderRepository.FindByID(ctx, orderID)
	if existingOrder == nil {
		return "", fmt.Errorf("order not found")
	}

	existingPayment, _ := s.paymentRepository.FindByOrderID(ctx, orderID)
	if existingPayment != nil {
		return "", fmt.Errorf("payment already exists")
	}

	payment := &Payment{
		ID:            primitive.NewObjectID(),
		OrderID:       orderID,
		Amount:        existingOrder.TotalPrice,
		Currency:      "vnd",
		Status:        Failed,
		PaymentMethod: "vnpay",
		ExpiredAt:     nowVN().Add(15 * time.Minute),
		CreatedAt:     time.Now(),
		UpdateAt:      time.Now(),
	}

	paymentID, err := s.paymentRepository.Create(ctx, payment)
	if err != nil {
		return "", err
	}

	params := s.buildVNPayParams(paymentID, existingOrder, clientIP)

	secureHash := s.createSecureHash(params)

	params["vnp_SecureHash"] = secureHash

	paymentURL := s.buildPaymentURL(params)

	return paymentURL, nil

}

func (s *paymentService) RepurchaseOrder(ctx context.Context, req *VNPayRequest, clientIP string) (*RepurchaseOrderResponse, error) {

	if req.OrderID == "" {
		return nil, fmt.Errorf("order id is required")
	}

	orderID, err := primitive.ObjectIDFromHex(req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("invalid order id: %v", err)
	}

	existingOrder, _ := s.orderRepository.FindByID(ctx, orderID)
	if existingOrder == nil {
		return nil, fmt.Errorf("order not found")
	}

	existingPayment, _ := s.paymentRepository.FindByOrderID(ctx, orderID)
	if existingPayment != nil {
		err := s.paymentRepository.DeletePayment(ctx, existingPayment.ID)
		if err != nil {
			return nil, err
		}
	}

	if existingOrder.Type == "cod" {

		err := s.orderRepository.UpdateByID(ctx, orderID, string(Pending))
		if err != nil {
			return nil, err
		}

		data := &RepurchaseOrderResponse{
			Link: "",
			Type: "cod",
		}

		return data, nil

	} else if existingOrder.Type == "vnpay" {
		
		payment := &Payment{
			ID:            primitive.NewObjectID(),
			OrderID:       orderID,
			Amount:        existingOrder.TotalPrice,
			Currency:      "vnd",
			Status:        Failed,
			PaymentMethod: "vnpay",
			ExpiredAt:     nowVN().Add(15 * time.Minute),
			CreatedAt:     time.Now(),
			UpdateAt:      time.Now(),
		}

		paymentID, err := s.paymentRepository.Create(ctx, payment)
		if err != nil {
			return nil, err
		}

		params := s.buildVNPayParams(paymentID, existingOrder, clientIP)

		secureHash := s.createSecureHash(params)

		params["vnp_SecureHash"] = secureHash

		paymentURL := s.buildPaymentURL(params)

		data := &RepurchaseOrderResponse{
			Link: paymentURL,
			Type: "vnpay",
		}

		return data, nil
		
	}

	return nil, nil
}

func (s *paymentService) buildVNPayParams(paymentID string, order *order.Order, clientIP string) map[string]string {

	create := nowVN()
	expire := create.Add(15 * time.Minute)

	var orderInfo string

	if orderInfo == "" {
		orderInfo = fmt.Sprintf("Thanh toan don hang #%s", order.ID.Hex())
	}

	params := map[string]string{
		"vnp_Version":    s.config.Version,
		"vnp_Command":    s.config.Command,
		"vnp_TmnCode":    s.config.TmnCode,
		"vnp_Amount":     strconv.FormatFloat(order.TotalPrice*100, 'f', 0, 64),
		"vnp_CreateDate": create.Format("20060102150405"),
		"vnp_CurrCode":   s.config.CurrCode,
		"vnp_IpAddr":     clientIP,
		"vnp_OrderType":  "other",
		"vnp_Locale":     "vn",
		"vnp_OrderInfo":  orderInfo,
		"vnp_ReturnUrl":  "https://monolith-architect.onrender.com/api/v1/payment/vnpay/callback",
		"vnp_ExpireDate": expire.Format("20060102150405"),
		"vnp_TxnRef":     paymentID,
		"vnp_BankCode":   "",
	}

	return params

}

func nowVN() time.Time {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		loc = time.FixedZone("ICT", 7*3600)
	}
	return time.Now().In(loc)
}

func (s *paymentService) createSecureHash(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "vnp_SecureHash" || k == "vnp_SecureHashType" {
			continue
		}

		if params[k] == "" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	for i, key := range keys {
		if i > 0 {
			b.WriteString("&")
		}
		b.WriteString(key)
		b.WriteString("=")
		b.WriteString(url.QueryEscape(params[key]))
	}
	mac := hmac.New(sha512.New, []byte(s.config.HashSecret))
	mac.Write([]byte(b.String()))
	return hex.EncodeToString(mac.Sum(nil))
}

func (s *paymentService) buildPaymentURL(params map[string]string) string {

	u, _ := url.Parse(s.config.PaymentUrl)

	q := u.Query()

	for key, value := range params {
		q.Set(key, value)
	}

	u.RawQuery = q.Encode()

	return u.String()
}

func (s *paymentService) HandleVNPayCallback(ctx context.Context, callback *VNPayCallback) error {

	isValid, err := s.VerifyCallback(callback)
	if err != nil {
		return fmt.Errorf("failed to verify callback: %w", err)
	}

	if !isValid {
		return fmt.Errorf("invalid callback signature")
	}

	paymentID, err := primitive.ObjectIDFromHex(callback.TransactionRef)
	if err != nil {
		return fmt.Errorf("invalid payment id: %w", err)
	}

	payment, err := s.paymentRepository.FindByID(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("failed to find payment: %w", err)
	}

	requestCallBack := &VNPayCallbackRequest{
		VNPayTransactionNo:   &callback.TransactionNo,
		VNPayTransactionRef:  &callback.TransactionRef,
		VNPayResponseCode:    &callback.ResponseCode,
		VNPayBankCode:        &callback.BankCode,
		VNPayTransactionInfo: &callback.OrderInfo,
	}

	err = s.paymentRepository.UpdateVnPay(ctx, paymentID, requestCallBack)
	if err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	switch callback.ResponseCode {
	case "00":
		payment.Status = Success
		err = s.orderRepository.UpdateByID(ctx, payment.OrderID, string(Success))
		if err != nil {
			return fmt.Errorf("failed to update order: %w", err)
		}
	case "24":
		payment.Status = Cancelled
		err = s.orderRepository.UpdateByID(ctx, payment.OrderID, string(Cancelled))
		if err != nil {
			return fmt.Errorf("failed to update order: %w", err)
		}
	default:
		payment.Status = Failed
		err = s.orderRepository.UpdateByID(ctx, payment.OrderID, string(Failed))
		if err != nil {
			return fmt.Errorf("failed to update order: %w", err)
		}
	}

	return s.paymentRepository.UpdateStatus(ctx, paymentID, payment.Status)
}

func (s *paymentService) VerifyCallback(callback *VNPayCallback) (bool, error) {

	params := map[string]string{
		"vnp_Amount":            callback.Amount,
		"vnp_BankCode":          callback.BankCode,
		"vnp_BankTranNo":        callback.BankTranNo,
		"vnp_CardType":          callback.CardType,
		"vnp_OrderInfo":         callback.OrderInfo,
		"vnp_PayDate":           callback.PayDate,
		"vnp_ResponseCode":      callback.ResponseCode,
		"vnp_TmnCode":           callback.TmnCode,
		"vnp_TransactionNo":     callback.TransactionNo,
		"vnp_TransactionStatus": callback.TransactionStatus,
		"vnp_TxnRef":            callback.TransactionRef,
	}

	expectedHash := s.createSecureHash(params)

	return expectedHash == callback.SecureHash, nil

}

func (s *paymentService) CronPaymentExpiration(ctx context.Context) error {

	payments, err := s.paymentRepository.FindByStatus(ctx)
	if err != nil {
		log.Printf("failed to find pending payments: %v", err)
	}

	for _, payment := range payments {
		if payment.ExpiredAt.Before(nowVN()) {
			err = s.paymentRepository.DeletePayment(ctx, payment.ID)
			if err != nil {
				log.Printf("failed to update payment status: %v", err)
			}

			err := s.orderRepository.DeleteByID(ctx, payment.OrderID)
			if err != nil {
				log.Printf("failed to delete order: %v", err)
			}
		}
	}

	return nil
}
