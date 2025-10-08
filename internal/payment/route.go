package payment

import (
	"net/http"

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
		paymentGroup.POST("/repurchase", handler.RepurchaseOrder)
		paymentGroup.GET("/vnpay/return", handler.HandleVNPayReturn)
		paymentGroup.POST("/vnpay/ipn", handler.HandleVNPayIpn)

		paymentGroup.HEAD("/vnpay/ipn", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})
		paymentGroup.GET("/vnpay/ipn", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})
		// VNPay
	}
}
