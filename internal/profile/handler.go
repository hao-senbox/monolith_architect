package profile

import (
	"modular_monolith/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	ProfileService ProfileService
}

func NewProfileHandler(profileService ProfileService) *ProfileHandler {
	return &ProfileHandler{
		ProfileService: profileService,
	}
}

func (h *ProfileHandler) CreateProfile(c *gin.Context) {
	
	var req CreateProfileRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	err := h.ProfileService.CreateProfile(c, &req)

	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "success", nil)
}