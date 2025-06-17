package order

type CreateOrderRequest struct {
	UserID string `json:"user_id" bson:"user_id"`
	Name   string `json:"name" bson:"name"`
	Email  string `json:"email" bson:"email"`
	Phone  string `json:"phone" bson:"phone"`
	Address string `json:"address" bson:"address"`
}