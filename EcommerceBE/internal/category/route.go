package category

import (
	"modular_monolith/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(route *gin.Engine, categoryHandler *CategoryHandler) {
	categoriesGroup := route.Group("/api/v1/category")
	{
		categoriesGroup.Use(middleware.JWTAuthMiddleware())
		categoriesGroup.POST("", categoryHandler.CreateCategory)
		categoriesGroup.GET("", categoryHandler.GetCategories)
		categoriesGroup.GET("/:id", categoryHandler.GetCategory)
		categoriesGroup.PUT("/:id", categoryHandler.UpdateCategory)
		categoriesGroup.DELETE("/:id", categoryHandler.DeleteCategory)
	}
}
