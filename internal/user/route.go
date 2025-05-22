package user

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, handler *UserHandler) {

	userGroup := r.Group("/api/v1/user") 
	{
		userGroup.GET("/:user_id", handler.GetUserByID)
		userGroup.POST("/register", handler.RegisterUser)
		userGroup.POST("/login", handler.LoginUser)	
	}
}