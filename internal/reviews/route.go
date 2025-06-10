package reviews

import (
	"modular_monolith/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *ReviewHandler) {
	reviewGroup := r.Group("/api/v1/review")
	{
		reviewGroup.Use(middleware.JWTAuthMiddleware())
		reviewGroup.POST("", handler.CreateReview)
		// reviewGroup.GET("", handler.GetAllReviews)
		// reviewGroup.GET("/:id", handler.GetReviewByID)
		// reviewGroup.PUT("/:id", handler.UpdateReview)
		// reviewGroup.DELETE("/:id", handler.DeleteReview)
	}
}