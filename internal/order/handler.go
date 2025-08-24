package order

import (
	"modular_monolith/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	OrderService OrderService
}

func NewOrderHandler(orderService OrderService) *OrderHandler {
	return &OrderHandler{
		OrderService: orderService,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {

	var req CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	id, err := h.OrderService.CreateOrder(c, &req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "success", id)
	
}

func (h *OrderHandler) GetAllOrders(c *gin.Context) {

	orders, err := h.OrderService.GetAllOrders(c)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", orders)
	
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {

	id := c.Param("id")

	order, err := h.OrderService.GetOrderByID(c, id)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", order)

}

func (h *OrderHandler) UpdateOrder(c *gin.Context) {

	id := c.Param("id")

	var req UpdateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	err := h.OrderService.UpdateOrder(c, &req, id)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", nil)

}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {

	id := c.Param("id")

	err := h.OrderService.DeleteOrder(c, id)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", nil)
	
}

func (h *OrderHandler) GetOrderByUserID(c *gin.Context) {

	userID := c.Param("user_id")

	orders, err := h.OrderService.GetOrderByUserID(c, userID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", orders)
}