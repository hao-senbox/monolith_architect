package payment

import (
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

	var redirectURL string
	switch callback.ResponseCode {
	case "00":
		redirectURL = "/payment/success?txn_ref=" + callback.TransactionRef
	case "24":
		redirectURL = "/payment/cancelled?txn_ref=" + callback.TransactionRef
	default:
		redirectURL = "/payment/failed?txn_ref=" + callback.TransactionRef
	}

	c.Redirect(http.StatusSeeOther, redirectURL)

}

func (h *PaymentHandler) HandleVNPayIPN(c *gin.Context) {
	
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

	helper.SendSuccess(c, http.StatusOK, "success", nil)
	
}