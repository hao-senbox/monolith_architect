package payment

import (
	"fmt"
	"io"
	"modular_monolith/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	PaymentService PaymentService
}

func NewPaymentHandler(paymentService PaymentService) *PaymentHandler {
	return &PaymentHandler{
		PaymentService: paymentService,
	}
}

func (h *PaymentHandler) CreatePaymentIntent(c *gin.Context) {

	var req CreatePaymentIntentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	paymentRes, err := h.PaymentService.CreatePaymentIntent(c, &req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", paymentRes)

}

func (h *PaymentHandler) StripeWebhook(c *gin.Context) {

	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	signature := c.Request.Header.Get("Stripe-Signature")
	if signature == "" {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	err = h.PaymentService.HandleWebhook(c, payload, signature)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", nil)

}

func (h *PaymentHandler) CreateVNPayPayment(c *gin.Context) {

	var req VNPayRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	clientIP := helper.GetClientIP(c)

	paymentRes, err := h.PaymentService.CreateVNPayPayment(c, &req, clientIP)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", paymentRes)
}

func (h *PaymentHandler) HandleVNPayCallback(c *gin.Context) {

	var callback VNPayCallback

	if err := c.ShouldBind(&callback); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	err := h.PaymentService.HandleVNPayCallback(c, &callback)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	var path string
	pathOrigin := "http://localhost:5173"
	switch callback.ResponseCode {
	case "00":
		path = "/payment/success"
	case "24":
		path = "/payment/cancelled"
	default:
		path = "/payment/failed"
	}
	redirectURL := fmt.Sprintf("%s%s?txn_ref=%s", pathOrigin, path, callback.TransactionRef)
	fmt.Printf("redirectURL: %s\n", redirectURL)
	c.Redirect(http.StatusSeeOther, redirectURL)

}

func (h *PaymentHandler) RepurchaseOrder(c *gin.Context) {

	var req VNPayRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	clientIP := helper.GetClientIP(c)

	data, err := h.PaymentService.RepurchaseOrder(c, &req, clientIP)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}


	helper.SendSuccess(c, http.StatusOK, "success", data)

}
