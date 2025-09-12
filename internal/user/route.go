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
		userGroup.GET("", handler.GetAllUsers)
		userGroup.GET("/:user_id", middleware.JWTAuthMiddleware(), handler.GetUserByID)
		userGroup.DELETE("/:user_id", handler.DeleteUser)
		userGroup.GET("/refresh", handler.RefreshToken)	
		userGroup.POST("/change-password", middleware.JWTAuthMiddleware(), handler.ChangePassword)
		userGroup.POST("/forgot-password", handler.ForgotPassword)
		userGroup.POST("/reset-password", handler.ResetPassword)
	}
}