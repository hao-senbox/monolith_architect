package cart

import (
	"modular_monolith/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	service CartService
}

func NewCartHandler(service CartService) *CartHandler {
	return &CartHandler{
		service: service,
	}
}

func (h *CartHandler) CreateCart(c *gin.Context) {
	
	var req AddtoCartRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	err := h.service.CreateCart(c, &req)

	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "success", nil)
	
}