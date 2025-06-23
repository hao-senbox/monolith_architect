package reviews

import (
	"fmt"
	"modular_monolith/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	ReviewService ReviewService
}

func NewReviewHandler(reviewService ReviewService) *ReviewHandler {
	return &ReviewHandler{ReviewService: reviewService}
}

func (r *ReviewHandler) CreateReview(c *gin.Context) {
	
	var req CreateReviewRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	err := r.ReviewService.CreateReview(c, &req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "success", nil)

}

func (r *ReviewHandler) GetAllReviews(c *gin.Context) {

	productID := c.Query("product_id")

	if productID == "" {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("product_id is required"), helper.ErrInvalidRequest)
		return
	}

	reviews, err := r.ReviewService.GetAllReviews(c, productID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", reviews)
	
}

func (r *ReviewHandler) GetReviewByID(c *gin.Context) {
	
	id := c.Param("id")

	review, err := r.ReviewService.GetReviewByID(c, id)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}	

	helper.SendSuccess(c, http.StatusOK, "success", review)
	
}

func (r *ReviewHandler) UpdateReview(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("id is required"), helper.ErrInvalidRequest)
		return
	}

	var req UpdateReviewRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	err := r.ReviewService.UpdateReview(c, id, &req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", nil)

}

func (r *ReviewHandler) DeleteReview(c *gin.Context) {

	id := c.Param("id")

	err := r.ReviewService.DeleteReview(c, id)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", nil)
}
