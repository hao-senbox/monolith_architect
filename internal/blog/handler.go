package blog

import (
	"fmt"
	"modular_monolith/helper"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (h *BlogHandler) GetBlogByID(c *gin.Context) {

	id := c.Param("id")

	res, err := h.service.GetBlogByID(c, id)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", res)

}

func (h *BlogHandler) UpdateBlog(c *gin.Context) {

	id := c.Param("id")

	var req UpdateBlogRequest
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

	err = h.service.UpdateBlog(c, id, &req, mainImage)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "success", nil)

}

func (h *BlogHandler) DeleteBlog(c *gin.Context) {

	id := c.Param("id")

	err := h.service.DeleteBlog(c, id)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", nil)

}

func (h *BlogHandler) LikeBlog(c *gin.Context) {

	id := c.Param("id")

	userIDRaw, exists := c.Get("user_id")
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("user_id not found"), helper.ErrInvalidRequest)
		return
	}
	
	userIDStr, ok := userIDRaw.(string)
	if !ok {
		helper.SendError(c, 400, fmt.Errorf("user_id must be string"), helper.ErrInvalidRequest)
		return
	}

	if strings.HasPrefix(userIDStr, "ObjectID(\"") && strings.HasSuffix(userIDStr, "\")") {
		userIDStr = userIDStr[10 : len(userIDStr)-2] 
	}
	
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	err = h.service.LikeBlog(c, id, userID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", nil)

}

func (h *BlogHandler) ViewBlog(c *gin.Context) {

	id := c.Param("id")

	err := h.service.ViewBlog(c, id)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", nil)

}