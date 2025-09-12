package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	FindAll(ctx context.Context) ([]*User, error)
	Create(ctx context.Context, user *User) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, userId primitive.ObjectID) (*UserWithProfile, error)
	FindByUserID(ctx context.Context, userID primitive.ObjectID) (*User, error)
	UpdateByID(ctx context.Context, userID primitive.ObjectID, updateFields bson.M) error
	DeleteByID(ctx context.Context, userID primitive.ObjectID) error
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) UserRepository {
	return &userRepository{collection: collection}
}

func (r *userRepository) DeleteByID(ctx context.Context, userID primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": userID})
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) FindAll(ctx context.Context) ([]*User, error) {

	var userAll []*User

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		userAll = append(userAll, &user)
	}

	return userAll, nil
	
}

func (r *userRepository) FindByID(ctx context.Context, userID primitive.ObjectID) (*UserWithProfile, error) {

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: userID}}}}

	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "profiles"},
		{Key: "localField", Value: "_id"},
		{Key: "foreignField", Value: "user_id"},
		{Key: "as", Value: "profile"},
	}}}


	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{
		{Key: "path", Value: "$profile"},
		{Key: "preserveNullAndEmptyArrays", Value: true},
	}}}

	cursor, err := r.collection.Aggregate(ctx, mongo.Pipeline{matchStage, lookupStage, unwindStage})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var user UserWithProfile
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		return &user, nil
	}
	return nil, nil
}

func (r *userRepository) Create(ctx context.Context, user *User) (*User, error) {
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	
	filter := bson.M{"email": email}

	var user User

	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
	
}

func (r *userRepository) UpdateByID(ctx context.Context, userID primitive.ObjectID, updateFields bson.M) error {

	filter := bson.M{"_id": userID}
	update := bson.M{"$set": updateFields}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
	
}

func (r *userRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID) (*User, error) {

	filter := bson.M{"_id": userID}

	var user User

	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
	
}