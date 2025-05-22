package profile

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type ProfileRepository interface {
	Create(ctx context.Context, profile *Profile) error
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