package category

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	FindAll(ctx context.Context) ([]*Category, error)
	FindByID(ctx context.Context, categoryID primitive.ObjectID) (*Category, error)
	UpdateByID(ctx context.Context, category *Category, categoryID primitive.ObjectID) error
	DeleteByID(ctx context.Context, categoryID primitive.ObjectID) error
}

type categoryRepository struct {
	collection *mongo.Collection
}

func NewCategoryRepository(collection *mongo.Collection) CategoryRepository {
	return &categoryRepository{
		collection: collection,
	}
}

func (r *categoryRepository) FindByID(ctx context.Context, categoryID primitive.ObjectID) (*Category, error) {
	
	filter := bson.M{"_id": categoryID}

	var category Category

	err := r.collection.FindOne(ctx, filter).Decode(&category)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepository) FindAll(ctx context.Context) ([]*Category, error) {

	var categories []*Category

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}	
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var category Category
		if err := cursor.Decode(&category); err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	return categories, nil
}

func (r *categoryRepository) Create(ctx context.Context, category *Category) error {
	_, err := r.collection.InsertOne(ctx, category)
	if err != nil {
		return err
	}
	return nil
}

func (r *categoryRepository) UpdateByID(ctx context.Context, category *Category, categoryID primitive.ObjectID) error {

	filter := bson.M{"_id": categoryID}
	update := bson.M{"$set": category}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil

}

func (r *categoryRepository) DeleteByID(ctx context.Context, categoryID primitive.ObjectID) error {

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": categoryID})
	if err != nil {
		return err
	}

	return nil
	
}
