package product

import "mime/multipart"

type CreateProductRequest struct {
	ProductName        string                     `json:"product_name"`
	ProductDescription string                     `json:"product_description"`
	CategoryID         string                     `json:"category_id"`
	Color              string                     `json:"color"`
	Price              float64                    `json:"price"`
	Discount           float64                    `json:"discount"`
	Currency           string                     `json:"currency"`
	Sizes              []CreateSizeOptionsRequest `json:"sizes"`
}

type ProductFiles struct {
	MainImage *multipart.FileHeader   `json:"main_image"`
	SubImages []*multipart.FileHeader `json:"sub_images"`
}

type CreateSizeOptionsRequest struct {
	Size  string `json:"size"`
	Stock int    `json:"stock"`
}

type UpdateProductRequest struct {
	CategoryID         string                     `json:"category_id" bson:"category_id"`
	ProductName        string                     `json:"product_name" bson:"product_name"`
	ProductDescription string                     `json:"product_description" bson:"product_description"`
	Color              string                     `json:"color" bson:"color"`
	Price              float64                    `json:"price"`
	Discount           float64                    `json:"discount"`
	Currency           string                     `json:"currency"`
	Sizes              []CreateSizeOptionsRequest `json:"sizes" bson:"sizes"`
}
