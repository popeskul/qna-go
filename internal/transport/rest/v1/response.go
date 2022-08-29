package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/sirupsen/logrus"
)

type signInResponse struct {
	AccessToken string      `json:"accessToken"`
	User        domain.User `json:"user"`
}

// errorResponse is the error response
type errorResponse struct {
	Message string `json:"message"`
}

// newErrorResponse creates a new error response and logs the error message
func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
