package category

import (
	"modular_monolith/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService CategoryService
}

func NewCategoryHandler(categoryService CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	
	var req CreateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	err := h.categoryService.CreateCategory(c, &req)

	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "success", nil)
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	
	categories, err := h.categoryService.GetCategories(c)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", categories)

}

func (h *CategoryHandler) GetCategory(c *gin.Context) {

	categoryID := c.Param("id")	

	category, err := h.categoryService.GetCategory(c, categoryID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", category)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {

	categoryID := c.Param("id")
	var req UpdateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	err := h.categoryService.UpdateCategory(c, &req, categoryID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", nil)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	
	categoryID := c.Param("id")

	err := h.categoryService.DeleteCategory(c, categoryID)

	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", nil)
	
}