package reviews

import (
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