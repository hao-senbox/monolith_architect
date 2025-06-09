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

func (h *CartHandler) GetCart(c *gin.Context) {

	userID := c.Param("user-id")

	carts, err := h.service.GetCartByUserID(c, userID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", carts)

}

func (h *CartHandler) UpdateCart(c *gin.Context) {

	var req UpdateCartRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	err := h.service.UpdateCart(c, &req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", nil)
}

func (h *CartHandler) DeleteItemCart(c *gin.Context) {

	var req DeleteItemCartRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	err := h.service.DeleteItemCart(c, &req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", nil)
}

func (h *CartHandler) DeleteCart(c *gin.Context) {

	userID := c.Param("user-id")

	err := h.service.DeleteCart(c, userID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", nil)
}