package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/popeskul/qna-go/internal/logger"
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
	logger := logger.GetLogger()
	logger.WithFields(logrus.Fields{
		"url":    c.Request.URL.String(),
		"method": c.Request.Method,
	}).Error(message)

	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
