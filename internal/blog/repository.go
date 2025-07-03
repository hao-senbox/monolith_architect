package blog

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BlogRepository interface{
	Create(ctx context.Context, blog *Blog) error
	FindAll(ctx context.Context) ([]*Blog, error)
	FindID(ctx context.Context, id primitive.ObjectID) (*Blog, error)
	UpdateByID(ctx context.Context, id primitive.ObjectID, blog *Blog) error
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	IncrementViews(ctx context.Context, id primitive.ObjectID) error
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

func (r *blogRepository) FindID(ctx context.Context, id primitive.ObjectID) (*Blog, error) {

	filter := bson.M{"_id": id}

	var blog *Blog

	err := r.collection.FindOne(ctx, filter).Decode(&blog)
	if err != nil {
		return nil, err
	}

	return blog, nil
	
}

func (r *blogRepository) UpdateByID(ctx context.Context, id primitive.ObjectID, blog *Blog) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.UpdateOne(ctx, filter, bson.M{"$set": blog})
	return err
}

func (r *blogRepository) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

func (r *blogRepository) IncrementViews(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.UpdateOne(ctx, filter, bson.M{"$inc": bson.M{"total_view": 1}})
	return err
}