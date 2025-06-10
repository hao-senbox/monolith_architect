package reviews

type CreateReviewRequest struct {
	ProductID string `json:"product_id" bson:"product_id"`
	UserID    string `json:"user_id" bson:"user_id"`
	Rating    int                `json:"rating" bson:"rating"`
	Review    string             `json:"review" bson:"review"`
}

type UpdateReviewRequest struct {
	Rating    int                `json:"rating" bson:"rating"`
	Review    string             `json:"review" bson:"review"`
}