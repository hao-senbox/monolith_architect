package user

import (
	"modular_monolith/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService UserService
}

func NewUserHandler(UserService UserService) *UserHandler {
	return &UserHandler{UserService: UserService}
}

func (h *UserHandler) RefreshToken(c *gin.Context) {

	refreshToken := c.Query("refresh_token")

	newToken, newRefreshToken, err := h.UserService.RefreshToken(refreshToken)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", map[string]string{
		"token":        newToken,
		"refreshToken": newRefreshToken,
	})
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {

	users, err := h.UserService.GetAllUsers(c)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", users)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {

	userID := c.Param("user_id")

	user, err := h.UserService.GetUserByID(c, userID)

	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", user)
}

func (h *UserHandler) RegisterUser(c *gin.Context) {

	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}
	
	user, err := h.UserService.RegisterUser(c, &req)

	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "success", user)

}

func (h *UserHandler) LoginUser(c *gin.Context) {

	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}
	
	user, err := h.UserService.LoginUser(c, req.Email, req.Password)

	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	
	userID := c.Param("user_id")

	err := h.UserService.DeleteUser(c, userID)

	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", nil)
}