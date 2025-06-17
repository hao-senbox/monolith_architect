package order

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepository interface{
	Create(ctx context.Context, order *Order) error
	FindAll(ctx context.Context) ([]Order, error)
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

func (r *orderRepository) FindAll(ctx context.Context) ([]Order, error) {

	var orders []Order

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	
	return orders, nil
	
}