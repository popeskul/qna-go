package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// errorResponse is the error response
// It contains the error message
type errorResponse struct {
	Message string `json:"message"`
}

// newErrorResponse creates a new error response
// It logs the error and returns a new error response
func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
