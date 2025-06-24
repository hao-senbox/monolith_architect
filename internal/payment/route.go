package payment

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *PaymentHandler) {
	paymentGroup := r.Group("/api/v1/payment")
	{
		paymentGroup.POST("/create-intent", handler.CreatePaymentIntent)
		r.POST("/api/v1/payment/webhook", handler.StripeWebhook)
	}
}
