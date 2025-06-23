package payment

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentRepository interface{
	Create(ctx context.Context, payment *Payment) error
	FindByOrderID(ctx context.Context, orderID primitive.ObjectID) (*Payment, error)
}

type paymentRepository struct{
	collection *mongo.Collection
}

func NewPaymentRepository(collection *mongo.Collection) PaymentRepository {
	return &paymentRepository{
		collection: collection,
	}
}

func (r *paymentRepository) Create(ctx context.Context, payment *Payment) error {

	_, err := r.collection.InsertOne(ctx, payment)
	if err != nil {
		return err
	}

	return nil

}

func (r *paymentRepository) FindByOrderID(ctx context.Context, orderID primitive.ObjectID) (*Payment, error) {

	var payment Payment

	err := r.collection.FindOne(ctx, bson.M{"order_id": orderID}).Decode(&payment)
	if err != nil {
		return nil, err
	}
	
	return &payment, nil
}
