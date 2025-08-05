package reviews

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReviewRepository interface {
	Create(ctx context.Context, review *Reviews) error
	FindAll(ctx context.Context, productID primitive.ObjectID) ([]*Reviews, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*Reviews, error)
	UpdateByID(ctx context.Context, id primitive.ObjectID, req *UpdateReviewRequest) error
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
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

func (r *reviewRepository) FindAll(ctx context.Context, productID primitive.ObjectID) ([]*Reviews, error) {

	var reviews []*Reviews
	filter := bson.M{
		"product_id": productID,
	}

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})
	
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &reviews); err != nil {
		return nil, err
	}

	return reviews, nil

}

func (r *reviewRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*Reviews, error) {

	var review Reviews

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&review)
	if err != nil {
		return nil, err
	}

	return &review, nil

}

func (r *reviewRepository) UpdateByID(ctx context.Context, id primitive.ObjectID, req *UpdateReviewRequest) error {

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": req})
	if err != nil {
		return err
	}

	return nil

}

func (r *reviewRepository) DeleteByID(ctx context.Context, id primitive.ObjectID) error {

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil

}
