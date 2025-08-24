package order

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepository interface{
	Create(ctx context.Context, order *Order) (string, error)
	FindAll(ctx context.Context) ([]Order, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*Order, error)
	UpdateByID(ctx context.Context, id primitive.ObjectID, status string) error
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]*Order, error)
}

type orderRepository struct{
	collection *mongo.Collection
}

func NewOrderRepository(collection *mongo.Collection) OrderRepository {
	return &orderRepository{collection: collection}
}

func (r *orderRepository) Create(ctx context.Context, order *Order) (string, error) {
	result, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		return "", err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()
	return id, err
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

func (r *orderRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*Order, error) {

	filter := bson.M{"_id": id}

	var order Order

	if err := r.collection.FindOne(ctx, filter).Decode(&order); err != nil {
		return nil, err
	}
	
	return &order, nil

}

func (r *orderRepository) UpdateByID(ctx context.Context, id primitive.ObjectID, status string) error {

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": status}}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	
	return err

}

func (r *orderRepository) DeleteByID(ctx context.Context, id primitive.ObjectID) error {

	filter := bson.M{"_id": id}

	_, err := r.collection.DeleteOne(ctx, filter)
	
	return err
	
}

func (r *orderRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]*Order, error) {

	filter := bson.M{"user_id": userID}

	var orders []*Order

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	
	return orders, nil

}