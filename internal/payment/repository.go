package payment

import (
	"context"
	"modular_monolith/internal/shared/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *Payment) (string, error)
	FindByOrderID(ctx context.Context, orderID primitive.ObjectID) (*model.Payment, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*Payment, error)
	FindByStatus(ctx context.Context, status PaymentStatus) ([]*Payment, error)
	FindByStripePaymentID(ctx context.Context, stripePaymentID string) (*Payment, error)
	UpdateStatus(ctx context.Context, paymentID primitive.ObjectID, status PaymentStatus) error
	UpdateVnPay(ctx context.Context, paymentID primitive.ObjectID, req *VNPayCallbackRequest) error
	DeletePayment(ctx context.Context, paymentID primitive.ObjectID) error
}

type paymentRepository struct {
	collection *mongo.Collection
}

func NewPaymentRepository(collection *mongo.Collection) PaymentRepository {
	return &paymentRepository{
		collection: collection,
	}
}

func (r *paymentRepository) Create(ctx context.Context, payment *Payment) (string, error) {

	result, err := r.collection.InsertOne(ctx, payment)
	if err != nil {
		return "", err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()
	return id, nil

}

func (r *paymentRepository) FindByOrderID(ctx context.Context, orderID primitive.ObjectID) (*model.Payment, error) {

	var payment model.Payment

	err := r.collection.FindOne(ctx, bson.M{"order_id": orderID}).Decode(&payment)
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *paymentRepository) FindByStripePaymentID(ctx context.Context, stripePaymentID string) (*Payment, error) {

	var payment Payment

	err := r.collection.FindOne(ctx, bson.M{"stripe_payment_id": stripePaymentID}).Decode(&payment)
	if err != nil {
		return nil, err
	}

	return &payment, nil

}

func (r *paymentRepository) UpdateStatus(ctx context.Context, paymentID primitive.ObjectID, status PaymentStatus) error {

	filter := bson.M{"_id": paymentID}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil

}

func (r *paymentRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*Payment, error) {

	var payment Payment

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&payment)
	if err != nil {
		return nil, err
	}

	return &payment, nil

}

func (r *paymentRepository) UpdateVnPay(ctx context.Context, paymentID primitive.ObjectID, req *VNPayCallbackRequest) error {

	filter := bson.M{"_id": paymentID}
	
	update := bson.M{
		"$set": bson.M{
			"vn_pay_transaction_no":   req.VNPayTransactionNo,
			"vn_pay_transaction_ref":  req.VNPayTransactionRef,
			"vn_pay_response_code":    req.VNPayResponseCode,
			"vn_pay_bank_code":        req.VNPayBankCode,
			"vn_pay_transaction_info": req.VNPayTransactionInfo,
			"updated_at":              time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *paymentRepository) FindByStatus(ctx context.Context, status PaymentStatus) ([]*Payment, error) {

	var payments []*Payment

	filter := bson.M{"status": status}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &payments); err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *paymentRepository) DeletePayment(ctx context.Context, paymentID primitive.ObjectID) error {
	filter := bson.M{"_id": paymentID}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}