package order

type CreateOrderRequest struct {
	UserID     string  `json:"user_id" bson:"user_id"`
	Name       string  `json:"name" bson:"name"`
	Email      string  `json:"email" bson:"email"`
	Phone      string  `json:"phone" bson:"phone"`
	Address    string  `json:"address" bson:"address"`
	CouponCode *string `json:"coupon_code" bson:"coupon_code"`
}

type UpdateOrderRequest struct {
	Status string `json:"status" bson:"status"`
}
