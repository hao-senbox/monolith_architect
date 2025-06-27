package coupon

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CouponType string

const (
	Public  CouponType = "public"
	Private CouponType = "private"
)

type Coupon struct {
	ID           primitive.ObjectID    `json:"id" bson:"_id"`
	Name         string                `json:"name" bson:"name"`
	CodeCoupon   string                `json:"code_coupon" bson:"code_coupon"`
	Discount     float64               `json:"discount" bson:"discount"`
	MaximumUse   *int                  `json:"maximum_use" bson:"maximum_use"`
	UserIsUsed   []primitive.ObjectID `json:"user_is_used" bson:"user_is_used"`
	AllowedUsers []primitive.ObjectID `json:"allowed_users" bson:"allowed_users"`
	Type         string                `json:"type" bson:"type"`
	ExpiredAt    time.Time             `json:"expired_at" bson:"expired_at"`
	CreatedAt    time.Time             `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at" bson:"updated_at"`
}
