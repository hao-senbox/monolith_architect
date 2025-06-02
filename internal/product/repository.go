package product

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	FindAll(ctx context.Context) ([]*Product, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*Product, error)
	UpdateByID(ctx context.Context, id primitive.ObjectID, product *Product) error	
}

type productRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(collection *mongo.Collection) ProductRepository {
	return &productRepository{
		collection: collection,
	}
}

func (r *productRepository) Create(ctx context.Context, product *Product) error {

	_, err := r.collection.InsertOne(ctx, product)
	if err != nil {
		return err
	}
	
	return nil
}

func (r *productRepository) FindAll(ctx context.Context) ([]*Product, error) {

	var products []*Product

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	
	return products, nil

}

func (r *productRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*Product, error) {

	filter := bson.M{"_id": id}

	var product *Product

	err := r.collection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *productRepository) UpdateByID(ctx context.Context, id primitive.ObjectID, product *Product) error {

	filter := bson.M{"_id": id}
	update := bson.M{"$set": product}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	
	return nil
	
}