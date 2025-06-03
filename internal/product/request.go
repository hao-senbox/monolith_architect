package product

import "mime/multipart"

type CreateProductRequest struct {
	ProductName        string                        `json:"product_name"`
	ProductDescription string                        `json:"product_description"`
	CategoryID         string                        `json:"category_id"`
	Variants           []CreateProductVariantRequest `json:"variants"`
}

type VariantFiles struct {
	MainImage *multipart.FileHeader
	SubImages []*multipart.FileHeader
}
type CreateProductVariantRequest struct {
	Color string `json:"color"`
	Sizes []CreateSizeOptionsRequest
}

type CreateSizeOptionsRequest struct {
	SKU      string  `json:"sku"`
	Size     string  `json:"size"`
	Stock    int     `json:"stock"`
	Price    float64 `json:"price"`
	Discount float64 `json:"discount"`
	Currency string  `json:"currency"`
}

type UpdateProductRequest struct {
	CategoryID         string                        `json:"category_id" bson:"category_id"`
	ProductName        string                        `json:"product_name" bson:"product_name"`
	ProductDescription string                        `json:"product_description" bson:"product_description"`
	Variants           []CreateProductVariantRequest `json:"variants" bson:"variants"`
}
