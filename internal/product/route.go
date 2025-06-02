package product

import (
	"modular_monolith/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes (r *gin.Engine, handler *ProductHandler) {
	productGroup := r.Group("/api/v1/product")
	{
		productGroup.Use(middleware.JWTAuthMiddleware())
		productGroup.POST("/", handler.CreateProduct)
		productGroup.GET("/", handler.GetAllProducts)
		productGroup.GET("/:id", handler.GetProductByID)
		// productGroup.PUT("/:id", handler.UpdateProduct)
		// productGroup.DELETE("/:product_id", handler.DeleteProduct)
	}
}