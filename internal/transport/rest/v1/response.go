package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// errorResponse - response for error
type errorResponse struct {
	Message string `json:"message"`
}

// statusResponse - response for status
type statusResponse struct {
	Status string `json:"status"`
}

// newErrorResponse - create new error response
func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
