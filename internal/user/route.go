package user

import (
	"modular_monolith/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *UserHandler) {

	userGroup := r.Group("/api/v1/user") 
	{
		userGroup.POST("/login", handler.LoginUser)
		userGroup.POST("/register", handler.RegisterUser)
		userGroup.POST("/logout", middleware.JWTAuthMiddleware(), handler.LogoutUser)
		userGroup.GET("", middleware.JWTAuthMiddleware(), handler.GetAllUsers)
		userGroup.GET("/:user_id", middleware.JWTAuthMiddleware(), handler.GetUserByID)
		userGroup.DELETE("/:user_id", middleware.JWTAuthMiddleware(), handler.DeleteUser)
		userGroup.GET("/refresh", handler.RefreshToken)	
	}
}