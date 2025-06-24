package profile

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *ProfileHandler) {
	profileGroup := r.Group("/api/v1/profile")
	{
		profileGroup.POST("", handler.CreateProfile)
		profileGroup.PUT("", handler.UpdateProfile)
	}
}
