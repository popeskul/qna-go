// Package v1 defines the handlers for the 1 version.
package v1

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/util"
	"net/http"
	"os"
	"time"
)

// Auth interface is implemented by the service.
type Auth interface {
	SignUp(ctx context.Context, user domain.User) error
	SignIn(ctx context.Context, user domain.User) (domain.User, error)
}

// SignUp godoc
// @Summary Sign up
// @Tags auth
// @Description Sign up
// @ID sign-up
// @Accept  json
// @Produce  json
// @Param user body domain.User true "user"
// @Success 201
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /sign-up [post]
func (h *Handlers) SignUp(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	password, err := util.HashPassword(user.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	user.Password = password

	if err := h.service.Auth.CreateUser(c, user); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

// SignIn godoc
// @Summary Sign in
// @Tags auth
// @Description Sign in
// @ID sign-in
// @Accept  json
// @Produce  json
// @Param user body domain.User true "user"
// @Success 200 {object} signInResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /sign-in [post]
func (h *Handlers) SignIn(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.Auth.GetUserByEmail(c, user.Email)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	duration, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	accessToken, err := h.service.TokenMaker.CreateToken(user.ID, duration)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, signInResponse{
		AccessToken: accessToken,
		User:        user,
	})
}
