package payment

import (
	"modular_monolith/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *PaymentHandler) {
	paymentGroup := r.Group("/api/v1/payment")
	{
		paymentGroup.Use(middleware.JWTAuthMiddleware())
		paymentGroup.POST("/create-intent", handler.CreatePaymentIntent)
		// paymentGroup.GET("/:id", handler.GetPaymentByOrderID)
		// paymentGroup.POST("/webhook", handler.StripeWebhook)
	}
}