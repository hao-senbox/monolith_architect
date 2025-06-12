package reviews

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReviewRepository interface {
	Create(ctx context.Context, review *Reviews) error
	FindAll(ctx context.Context) ([]*Reviews, error)
}

type reviewRepository struct {
	collection *mongo.Collection
}

func NewReviewRepository(collection *mongo.Collection) ReviewRepository {
	return &reviewRepository{collection: collection}
}

func (r *reviewRepository) Create(ctx context.Context, review *Reviews) error {
	_, err := r.collection.InsertOne(ctx, review)
	if err != nil {
		return err
	}
	return nil
}

func (r *reviewRepository) FindAll(ctx context.Context) ([]*Reviews, error) {

	var reviews []*Reviews
	
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &reviews); err != nil {
		return nil, err
	}

	return reviews, nil
	
}