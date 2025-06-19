package payment

import "go.mongodb.org/mongo-driver/mongo"

type PaymentRepository interface{

}

type paymentRepository struct{
	collection *mongo.Collection
}

func NewPaymentRepository(collection *mongo.Collection) PaymentRepository {
	return &paymentRepository{
		collection: collection,
	}
}
