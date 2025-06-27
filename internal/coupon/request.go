package coupon

type CreateCouponRequest struct {
	Name       string  `json:"name" bson:"name"`
	Discount   float64 `json:"discount" bson:"discount"`
	MaximumUse int     `json:"maximum_use" bson:"maximum_use"`
	Type       string  `json:"type" bson:"type"`
	AllowedUsers []string `json:"allowed_users" bson:"allowed_users"`
	ExpiredAt  string  `json:"expired_at" bson:"expired_at"`
}

type CanUseCouponRequest struct {
	CouponCode string `json:"coupon_code" bson:"coupon_code"`
	UserID string `json:"user_id" bson:"user_id"`
}	