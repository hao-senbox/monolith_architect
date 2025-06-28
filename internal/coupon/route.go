package coupon

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, handler *CouponHandler) {
	couponGroup := r.Group("/api/v1/coupon")
	{
		couponGroup.POST("", handler.CreateCoupon)
		couponGroup.GET("", handler.GetAllCoupons)
		couponGroup.GET("/:code", handler.GetCouponByCode)
		couponGroup.PUT("/:id", handler.UpdateCoupon)
		couponGroup.DELETE("/:id", handler.DeleteCoupon)
		couponGroup.GET("user/:user_id", handler.GetCouponByUserID)
		couponGroup.POST("/can_use_coupon", handler.CanUseCoupon)
	}
}