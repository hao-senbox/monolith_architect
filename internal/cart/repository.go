package cart

import (
	"context"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CartRepository interface {
	Create(ctx context.Context, cart *Cart) error
	FindCartByUserID(ctx context.Context, userID primitive.ObjectID) (*Cart, error)
	AddToCart(ctx context.Context, cartItem *CartItem, userID primitive.ObjectID) error
	UpdateCart(ctx context.Context, productID primitive.ObjectID, userID primitive.ObjectID, quantity int, types string) error
	DeleteItemCart(ctx context.Context, productID primitive.ObjectID, userID primitive.ObjectID) error
	DeleteCart(ctx context.Context, userID primitive.ObjectID) error
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
			ID:         primitive.NewObjectID(),
			UserID:     userID,
			CartItems:  []*CartItem{},
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

func (r *cartRepository) AddToCart(ctx context.Context, cartItem *CartItem, userID primitive.ObjectID) error {

	cart, err := r.FindCartByUserID(ctx, userID)
	if err != nil {
		return err
	}

	found := false

	for i, item := range cart.CartItems {
		if item.ProductID == cartItem.ProductID {
			cart.CartItems[i].Quantity += cartItem.Quantity
			found = true
			break
		}
	}

	if !found {
		cart.CartItems = append(cart.CartItems, cartItem)
	}

	return r.updateTotalPrice(ctx, cart)
}

func (r *cartRepository) updateCart(ctx context.Context, cart *Cart) error {

	filter := bson.M{"_id": cart.ID}

	update := bson.M{
		"$set": bson.M{
			"user_id":     cart.UserID,
			"cart_items":  cart.CartItems,
			"total_price": cart.TotalPrice,
			"created_at":  cart.CreatedAt,
			"updated_at":  cart.UpdatedAt,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil

}

func (r *cartRepository) UpdateCart(ctx context.Context, productID primitive.ObjectID, userID primitive.ObjectID, quantity int, types string) error {
	cart, err := r.FindCartByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if types == "add" {
		for i, item := range cart.CartItems {
			if item.ProductID == productID {
				cart.CartItems[i].Quantity += quantity
				break
			}
		}
	} else if types == "remove" {
		for i, item := range cart.CartItems {
			if item.ProductID == productID {
				if cart.CartItems[i].Quantity > quantity {
					cart.CartItems[i].Quantity -= quantity
				} else {
					err := r.DeleteItemCart(ctx, productID, userID)
					if err != nil {
						return err
					}
					return nil
				}
				break
			}
		}
	}

	return r.updateTotalPrice(ctx, cart)
}

func (r *cartRepository) updateTotalPrice(ctx context.Context, cart *Cart) error {

	totalPrice := 0.0

	for _, item := range cart.CartItems {
		totalPrice += item.Price * float64(item.Quantity)
	}

	cart.TotalPrice = math.Round(totalPrice*100) / 100
	cart.UpdatedAt = time.Now()

	return r.updateCart(ctx, cart)

}

func (r *cartRepository) DeleteItemCart(ctx context.Context, productID primitive.ObjectID, userID primitive.ObjectID) error {

	cart, err := r.FindCartByUserID(ctx, userID)
	if err != nil {
		return err
	}

	for i, item := range cart.CartItems {
		if item.ProductID == productID {
			cart.CartItems = append(cart.CartItems[:i], cart.CartItems[i+1:]...)
			break
		}
	}

	return r.updateTotalPrice(ctx, cart)

}

func (r *cartRepository) DeleteCart(ctx context.Context, userID primitive.ObjectID) error {

	filter := bson.M{"user_id": userID}

	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
	
}