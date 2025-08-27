package product

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	FindAll(ctx context.Context, filter *ProductFilter) ([]*Product, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*Product, error)
	UpdateByID(ctx context.Context, id primitive.ObjectID, product *Product) error
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
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

func (r *productRepository) FindAll(ctx context.Context, filter *ProductFilter) ([]*Product, error) {

	query := bson.M{}

	if filter.Name != "" {
		query["product_name"] = bson.M{"$regex": filter.Name, "$options": "i"}
	}

	if filter.CategoryID != "" {
		objID, err := primitive.ObjectIDFromHex(filter.CategoryID)
		if err == nil {
			query["category_id"] = objID
		}
	}

	if filter.MinPrice > 0 || filter.MaxPrice > 0 {
		priceQuery := bson.M{}
		if filter.MinPrice > 0 {
			priceQuery["$gte"] = filter.MinPrice
		}
		if filter.MaxPrice > 0 {
			priceQuery["$lte"] = filter.MaxPrice
		}
		query["price"] = priceQuery
	}

	if filter.Size != "" {
		query["sizes.size"] = filter.Size
	}
	if filter.Surface != "" {
		query["surface"] = filter.Surface
	}

	sortQuery := bson.M{}

	switch filter.Sort {
	case "price-asc":
		sortQuery["price"] = 1
	case "price-desc":
		sortQuery["price"] = -1
	case "name-asc":
		sortQuery["product_name"] = 1
	case "name-desc":
		sortQuery["product_name"] = -1
	}

	opts := options.Find().SetSort(sortQuery)

    cursor, err := r.collection.Find(ctx, query, opts)
    if err != nil {
        return nil, err
    }

    var products []*Product
    if err := cursor.All(ctx, &products); err != nil {
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

func (r *productRepository) DeleteByID(ctx context.Context, id primitive.ObjectID) error {

	filter := bson.M{"_id": id}

	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil

}
