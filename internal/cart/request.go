package cart

type AddtoCartRequest struct {
	ProductID string `json:"product_id" bson:"product_id"`
	UserID    string `json:"user_id" bson:"user_id"`
	Quantity  int    `json:"quantity" bson:"quantity"`
}
