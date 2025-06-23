package cart

import (
	"modular_monolith/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes (r *gin.Engine, handler *CartHandler) {
	cartGroup := r.Group("/api/v1/cart")
	{
		cartGroup.Use(middleware.JWTAuthMiddleware())
		cartGroup.POST("", handler.CreateCart)
		cartGroup.GET("/:user-id", handler.GetCart)
		cartGroup.PUT("", handler.UpdateCart)
		cartGroup.DELETE("", handler.DeleteItemCart)
		cartGroup.DELETE("/:user-id", handler.DeleteCart)
	}
}