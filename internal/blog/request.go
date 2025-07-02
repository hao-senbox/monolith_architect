package blog

type CreateBlogRequest struct {
	UserID  string `json:"user_id" bson:"user_id"`
	Title   string `json:"title" bson:"title"`
	Content string `json:"content" bson:"content"`
}