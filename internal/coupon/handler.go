package coupon

import (
	"modular_monolith/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CouponHandler struct {
	CouponService CouponService
}

func NewCouponHandler(couponService CouponService) *CouponHandler {
	return &CouponHandler{
		CouponService: couponService,
	}
}

func (h *CouponHandler) CreateCoupon(c *gin.Context) {

	var req CreateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}
	
	err := h.CouponService.CreateCoupon(c, &req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "success", nil)
	
}

func (h *CouponHandler) GetAllCoupons(c *gin.Context) {

	coupons, err := h.CouponService.GetAllCoupons(c)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", coupons)

}

func (h *CouponHandler) GetCouponByCode(c *gin.Context) {

	code := c.Param("code")

	coupon, err := h.CouponService.GetCouponByCode(c, code)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", coupon)
}

func (h *CouponHandler) GetCouponByUserID(c *gin.Context) {

	id := c.Param("user_id")

	coupon, err := h.CouponService.GetCouponByUserID(c, id)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", coupon)

}

func (h *CouponHandler) CanUseCoupon(c *gin.Context) {

	var req CanUseCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	res, err := h.CouponService.CanUseCoupon(c, &req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", res)

}

func (h *CouponHandler) DeleteCoupon(c *gin.Context) {

	id := c.Param("id")

	err := h.CouponService.DeleteCoupon(c, id)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", nil)
	
}