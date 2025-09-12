package user

import (
	"errors"
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

func (h *UserHandler) LogoutUser(c *gin.Context) {

	userIdInterface, ok := c.Get("user_id")

	if !ok {
		helper.SendError(c, http.StatusUnauthorized, nil, helper.ErrInvalidOperation)
		return
	}

	userId, ok := userIdInterface.(string)
	if !ok {
		helper.SendError(c, http.StatusInternalServerError, nil, helper.ErrInvalidOperation)
		return
	}

	err := h.UserService.LogoutUser(c, userId)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", nil)

}

func (h *UserHandler) ChangePassword(c *gin.Context) {

	var req ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, errors.New("unauthorized"), helper.ErrInvalidOperation)
		return
	}

	err := h.UserService.ChangePassword(c, req, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", nil)

}

func (h *UserHandler) ForgotPassword(c *gin.Context) {

	var req ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	err := h.UserService.ForgotPassword(c, req.Email)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", nil)

}

func (h *UserHandler) ResetPassword(c *gin.Context) {

	var req ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	err := h.UserService.ResetPassword(c, req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", nil)
	
}
