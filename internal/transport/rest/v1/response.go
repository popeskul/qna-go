package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	ErrTestNotFound = errors.New("test id is required")
	ErrPermission   = errors.New("permission denied")
)

// errorResponse is the error response
type errorResponse struct {
	Message string `json:"message"`
}

// newErrorResponse creates a new error response and logs the error message
func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
