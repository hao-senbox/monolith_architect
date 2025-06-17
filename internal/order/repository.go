package order

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepository interface{
	Create(ctx context.Context, order *Order) error
}

type orderRepository struct{
	collection *mongo.Collection
}

func NewOrderRepository(collection *mongo.Collection) OrderRepository {
	return &orderRepository{collection: collection}
}

func (r *orderRepository) Create(ctx context.Context, order *Order) error {
	_, err := r.collection.InsertOne(ctx, order)
	return err
}