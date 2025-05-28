package profile

import (
	"modular_monolith/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *ProfileHandler) {
	profileGroup := r.Group("/api/v1/profile")
	{
		profileGroup.Use(middleware.JWTAuthMiddleware())
		profileGroup.POST("", handler.CreateProfile)
		profileGroup.PUT("", handler.UpdateProfile)
	}
}