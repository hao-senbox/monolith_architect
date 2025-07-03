package blog

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

	type Blog struct {
		ID            primitive.ObjectID   `json:"id" bson:"_id"`
		Title         string               `json:"title" bson:"title"`
		Content       string               `json:"content" bson:"content"`
		ImageURL      string               `json:"image_url" bson:"image_url"`
		ImagePublicID string               `json:"image_public_id" bson:"image_public_id"`
		TotalView     int                  `json:"total_view" bson:"total_view"`
		LikedUsers    []primitive.ObjectID `json:"liked_users" bson:"liked_users"`
		DislikedUsers []primitive.ObjectID `json:"disliked_users" bson:"disliked_users"`
		AuthorID      primitive.ObjectID   `json:"author_id" bson:"author_id"`
		Author        *UserInfor           `json:"author" bson:"author"`
		Created       time.Time            `json:"created" bson:"created"`
		Updated       time.Time            `json:"updated" bson:"updated"`
	}

type UserInfor struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	FullName string             `json:"full_name" bson:"full_name"`
	Avatar   *string            `json:"avatar" bson:"avatar"`
}
