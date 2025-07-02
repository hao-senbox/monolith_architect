package blog

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BlogRepository interface{
	Create(ctx context.Context, blog *Blog) error
	FindAll(ctx context.Context) ([]*Blog, error)
}

type blogRepository struct{
	collection *mongo.Collection
}

func NewBlogRepository(collection *mongo.Collection) BlogRepository {
	return &blogRepository{collection: collection}
}

func (r *blogRepository) Create(ctx context.Context, blog *Blog) error {
	_, err := r.collection.InsertOne(ctx, blog)
	return err
}

func ( r *blogRepository) FindAll(ctx context.Context) ([]*Blog, error) {
	
	var blogs []*Blog

	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &blogs); err != nil {
		return nil, err
	}
	
	return blogs, nil
	
}