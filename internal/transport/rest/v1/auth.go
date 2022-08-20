// Package v1 defines the handlers for the 1 version.
package v1

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/popeskul/qna-go/internal/domain"
	"net/http"
)

type Auth interface {
	SignUp(ctx context.Context, user domain.SignUpInput) error
	SignIn(ctx context.Context, user domain.SignInInput) (domain.User, error)
}

type SignInResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

type SignUpResponse struct {
	Status string `json:"status"`
	ID     int    `json:"id"`
}

// SignUp godoc
// @Tags auth
// @Summary Sign up
// @Description Sign up a new user
// @ID sign-up
// @Accept  json
// @Produce  json
// @Param user body domain.SignUpInput true "User"
// @Success 200 {object} SignUpResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /auth/sign-up [post]
func (h *Handlers) SignUp(c *gin.Context) {
	var user domain.SignUpInput
	if err := c.ShouldBindJSON(&user); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.Auth.CreateUser(user)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, SignUpResponse{
		Status: "success",
		ID:     id,
	})
}

// SignIn
// @Summary Sign in
// @Tags auth
// @Description Sign in user
// @ID sign-in
// @Accept  json
// @Produce  json
// @Param user body domain.SignInInput true "User"
// @Success 200 {object} SignInResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /auth/sign-in [post]
func (h *Handlers) SignIn(c *gin.Context) {
	var input domain.SignInInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.service.Auth.GenerateToken(input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, SignInResponse{
		Status: "success",
		Token:  token,
	})
}
