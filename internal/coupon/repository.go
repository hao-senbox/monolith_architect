package coupon

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CouponRepository interface {
	Create(ctx context.Context, coupon *Coupon) error
	FindAll(ctx context.Context) ([]*Coupon, error)
	FindByCode(ctx context.Context, id string) (*Coupon, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
	FindAllCouponsByUserID(ctx context.Context, userID primitive.ObjectID) ([]*Coupon, error)
	CheckCodeCoupon(ctx context.Context, codeCoupon string) (bool, error)
	AddUserIsUsed(ctx context.Context, userID primitive.ObjectID, codeCoupon string) error
	
}

type couponRepository struct {
	collection *mongo.Collection
}

func NewCouponRepository(collection *mongo.Collection) CouponRepository {
	return &couponRepository{
		collection: collection,
	}
}

func (r *couponRepository) Create(ctx context.Context, coupon *Coupon) error {
	_, err := r.collection.InsertOne(ctx, coupon)
	return err
}

func (r *couponRepository) FindAll(ctx context.Context) ([]*Coupon, error) {

	var coupons []*Coupon

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &coupons); err != nil {
		return nil, err
	}

	return coupons, nil

}

func (r *couponRepository) FindByCode(ctx context.Context, id string) (*Coupon, error) {

	filter := bson.M{"code_coupon": id}

	var coupon *Coupon

	if err := r.collection.FindOne(ctx, filter).Decode(&coupon); err != nil {
		return nil, err
	}

	return coupon, nil
}

func (r *couponRepository) FindAllCouponsByUserID(ctx context.Context, userID primitive.ObjectID) ([]*Coupon, error) {

	filter := bson.M{"allowed_users.user_id": userID}

	var coupons []*Coupon

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &coupons); err != nil {
		return nil, err
	}

	return coupons, nil
}

func (r *couponRepository) CheckCodeCoupon(ctx context.Context, codeCoupon string) (bool, error) {

	filter := bson.M{"code_coupon": codeCoupon}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil

}

func (r *couponRepository) AddUserIsUsed(ctx context.Context, userID primitive.ObjectID, codeCoupon string) error {

	filter := bson.M{"code_coupon": codeCoupon}
	update := bson.M{"$push": bson.M{"user_is_used": userID}}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err

}

func (r *couponRepository) Delete(ctx context.Context, id primitive.ObjectID) error {

	filter := bson.M{"_id": id}

	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil

}