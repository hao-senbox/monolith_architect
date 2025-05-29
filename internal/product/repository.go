package product

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
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