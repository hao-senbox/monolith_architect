package blog

import (
	"modular_monolith/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BlogHandler struct {
	service BlogService
}

func NewBlogHandler(service BlogService) *BlogHandler {
	return &BlogHandler{
		service: service,
	}
}

func (h *BlogHandler) CreateBlog(c *gin.Context) {

	var req CreateBlogRequest
	if err := c.ShouldBind(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	req.UserID = c.PostForm("user_id")
	req.Content = c.PostForm("content")
	req.Title = c.PostForm("title")

	err := c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	mainImage, err := c.FormFile("main_image")
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	err = h.service.CreateBlog(c, &req, mainImage)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}
	helper.SendSuccess(c, http.StatusCreated, "success", nil)

}

func (h *BlogHandler) GetAllBlogs(c *gin.Context) {

	res, err := h.service.GetAllBlog(c)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}
	helper.SendSuccess(c, http.StatusOK, "success", res)

}
