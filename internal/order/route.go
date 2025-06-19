package order

import (
	"modular_monolith/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *OrderHandler) {

	orderGroup := r.Group("/api/v1/order")
	{
		orderGroup.Use(middleware.JWTAuthMiddleware())
		orderGroup.POST("", handler.CreateOrder)
		orderGroup.GET("", handler.GetAllOrders)
		orderGroup.GET("/:id", handler.GetOrderByID)
		orderGroup.PUT("/:id", handler.UpdateOrder)
		orderGroup.DELETE("/:id", handler.DeleteOrder)
	}
}