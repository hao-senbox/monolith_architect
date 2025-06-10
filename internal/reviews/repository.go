package reviews

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type ReviewRepository interface {
	Create(ctx context.Context, review *Reviews) error
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