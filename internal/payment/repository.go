package payment

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *Payment) (string, error)
	FindByOrderID(ctx context.Context, orderID primitive.ObjectID) (*Payment, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*Payment, error)
	FindByStripePaymentID(ctx context.Context, stripePaymentID string) (*Payment, error)
	UpdateStatus(ctx context.Context, paymentID primitive.ObjectID, status PaymentStatus) error
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

func (r *paymentRepository) FindByOrderID(ctx context.Context, orderID primitive.ObjectID) (*Payment, error) {

	var payment Payment

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