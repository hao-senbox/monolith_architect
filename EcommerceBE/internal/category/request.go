package category

type CreateCategoryRequest struct {
	CategoryName string `json:"category_name" bson:"category_name"`
	ParentID     *string `json:"parent_id" bson:"parent_id"`
}

type UpdateCategoryRequest struct {
	CategoryName string `json:"category_name" bson:"category_name"`
	ParentID     *string `json:"parent_id" bson:"parent_id"`
}
