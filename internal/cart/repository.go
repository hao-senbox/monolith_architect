package cart

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CartRepository interface {
	Create(ctx context.Context, cart *Cart) error
	FindCartByUserID(ctx context.Context, userID primitive.ObjectID) (*Cart, error)
	AddToCart(ctx context.Context, productItem *Product, userID primitive.ObjectID, quantity int) error	
}

type cartRepository struct {
	collection *mongo.Collection
}

func NewCartRepository(collection *mongo.Collection) CartRepository {
	return &cartRepository{
		collection: collection,
	}
}

func (r *cartRepository) Create(ctx context.Context, cart *Cart) error {
	_, err := r.collection.InsertOne(ctx, cart)
	if err != nil {
		return err
	}
	return nil
}

func (r *cartRepository) FindCartByUserID(ctx context.Context, userID primitive.ObjectID) (*Cart, error) {
	
	filter := bson.M{"user_id": userID}

	var cart *Cart

	err := r.collection.FindOne(ctx, filter).Decode(&cart)
	if err == mongo.ErrNoDocuments {

		cart = &Cart{
			UserID:     userID,
			CartItems:  []CartItem{},
			TotalPrice: 0.0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		err = r.Create(ctx, cart)
		if err != nil {
			return nil, err
		}

		return cart, nil
	}
	
	return cart, nil
}

func (r *cartRepository) AddToCart(ctx context.Context, productItem *Product, userID primitive.ObjectID, quantity int) error {

	cart, err := r.FindCartByUserID(ctx, userID)
	if err != nil {
		return err
	}

	cart.CartItems = append(cart.CartItems, CartItem{
		Products:   []Product{},
		TotalPrice: 0.0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})

	return nil
}

