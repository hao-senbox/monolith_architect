package payment

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *PaymentHandler) {
	paymentGroup := r.Group("/api/v1/payment")
	{
		// Stripe
		paymentGroup.POST("/create-intent", handler.CreatePaymentIntent)
		paymentGroup.POST("/webhook", handler.StripeWebhook)
		// Stripe

		// VNPay
		paymentGroup.POST("/vnpay", handler.CreateVNPayPayment)
		paymentGroup.POST("/vnpay/callback", handler.HandleVNPayCallback)
		paymentGroup.POST("/vnpay/ipn", handler.HandleVNPayIPN)
		// VNPay
	}
}
