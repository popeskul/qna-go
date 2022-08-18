package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/popeskul/qna-go/internal/domain"
	"net/http"
)

// SignUp godoc
// @Tags auth
// @Summary Sign up
// @Description Sign up a new user
// @ID sign-up
// @Accept  json
// @Produce  json
// @Param user body domain.SignUpInput true "User"
// @Success 200 {object} domain.User
// @Failure 400 {object} error: error.Error
// @Failure 500 {object} error: error.Error
// @Router /auth/sign-up [post]
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

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"id":     id,
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
// @Success 200 {object} domain.User
// @Failure 400 {object} error: error.Error
// @Failure 500 {object} error: error.Error
// @Router /auth/sign-in [post]
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
