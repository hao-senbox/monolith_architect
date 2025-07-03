package blog

import (
	"modular_monolith/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *BlogHandler) {

	blogGroup := r.Group("/api/v1/blog")
	{
		blogGroup.POST("", handler.CreateBlog)
		blogGroup.GET("", handler.GetAllBlogs)
		blogGroup.GET("/:id", handler.GetBlogByID)
		blogGroup.PUT("/:id", handler.UpdateBlog)
		blogGroup.DELETE("/:id", handler.DeleteBlog)
		blogGroup.POST("view/:id", handler.ViewBlog)
		blogGroup.POST("/like/:id", middleware.JWTAuthMiddleware(), handler.LikeBlog)
	}

}

