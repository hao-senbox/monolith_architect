package product

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *ProductHandler) {
	productGroup := r.Group("/api/v1/product")
	{
		productGroup.POST("", handler.CreateProduct)
		productGroup.GET("", handler.GetAllProducts)
		productGroup.GET("/:id", handler.GetProductByID)
		productGroup.PUT("/:id", handler.UpdateProduct)
		productGroup.DELETE("/:id", handler.DeleteProduct)
	}
}
