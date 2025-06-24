package cart

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *CartHandler) {
	cartGroup := r.Group("/api/v1/cart")
	{
		cartGroup.POST("", handler.CreateCart)
		cartGroup.GET("/:user-id", handler.GetCart)
		cartGroup.PUT("", handler.UpdateCart)
		cartGroup.DELETE("", handler.DeleteItemCart)
		cartGroup.DELETE("/:user-id", handler.DeleteCart)
	}
}
