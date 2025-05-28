package product

import "go.mongodb.org/mongo-driver/mongo"

type ProductRepository interface {

}

type productRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(collection *mongo.Collection) ProductRepository {
	return &productRepository{
		collection: collection,
	}

}