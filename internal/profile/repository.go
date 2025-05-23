package profile

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProfileRepository interface {
	Create(ctx context.Context, profile *Profile) error
	FindByUserID(ctx context.Context, userID primitive.ObjectID) (*Profile, error)
}

type profileRepository struct {
	collection *mongo.Collection
}

func NewProfileRepository (collection *mongo.Collection) ProfileRepository {
	return &profileRepository{collection: collection}
}

func (r *profileRepository) Create(ctx context.Context, profile *Profile) error {

	_, err := r.collection.InsertOne(ctx, profile)
	if err != nil {
		return err
	}

	return nil
}

func (r *profileRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID) (*Profile, error) {

	filter := bson.M{"user_id": userID}

	var profile Profile

	err := r.collection.FindOne(ctx, filter).Decode(&profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil

}