package reviews

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *ReviewHandler) {
	reviewGroup := r.Group("/api/v1/review")
	{
		reviewGroup.POST("", handler.CreateReview)
		reviewGroup.GET("", handler.GetAllReviews)
		reviewGroup.GET("/:id", handler.GetReviewByID)
		reviewGroup.PUT("/:id", handler.UpdateReview)
		reviewGroup.DELETE("/:id", handler.DeleteReview)
	}
}
