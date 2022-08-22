// Package v1 defines the handlers for the 1 version.
package v1

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/popeskul/qna-go/internal/domain"
	"net/http"
)

// Auth interface is implemented by the service.
type Auth interface {
	SignUp(ctx context.Context, user domain.SignUpInput) error
	SignIn(ctx context.Context, user domain.SignInInput) (domain.User, error)
}

type SignInResponse struct {
	Token string `json:"token"`
}

type SignUpResponse struct {
	Status string `json:"status"`
}

func (h *Handlers) SignUp(c *gin.Context) {
	var user domain.SignUpInput
	if err := c.ShouldBindJSON(&user); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Auth.CreateUser(c, user); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handlers) SignIn(c *gin.Context) {
	var input domain.SignInInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.service.Auth.GenerateToken(c, input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, SignInResponse{
		Token: token,
	})
}
