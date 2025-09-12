package cart

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *CartHandler) {
	cartGroup := r.Group("/api/v1/cart")
	{
		cartGroup.POST("", handler.CreateCart)
		cartGroup.GET("/:user_id", handler.GetCart)
		cartGroup.PUT("", handler.UpdateCart)
		cartGroup.DELETE("", handler.DeleteItemCart)
		cartGroup.DELETE("/:user_id", handler.DeleteCart)
	}
}
