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

func (h *Handlers) InitAuthRouter(v1 *gin.RouterGroup) {
	usersAPI := v1.Group("/auth")
	{
		usersAPI.POST("/sign-up", h.SignUp)
		usersAPI.POST("/sign-in", h.SignIn)
	}
}

func (h *Handlers) SignUp(c *gin.Context) {
	var user domain.SignUpInput
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.service.Auth.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"id":     id,
	})
}

func (h *Handlers) SignIn(c *gin.Context) {
	var input domain.SignInInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Auth.GenerateToken(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"token":  token,
	})
}