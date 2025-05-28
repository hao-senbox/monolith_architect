package product

import "github.com/gin-gonic/gin"

type ProductHandler struct {
	ProductService ProductService
}

func NewProductHandler(productService ProductService) *ProductHandler {
	return &ProductHandler{
		ProductService: productService,
	}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	
}	