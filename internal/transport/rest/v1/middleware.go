// Package v1 defines all middleware for the 1 version.
package v1

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/popeskul/qna-go/internal/token"
	"net/http"
	"time"
)

const (
	authorizationPayloadKey = "authorization_payload"
)

var (
	ErrUserIdNotFound    = errors.New("user id not found")
	ErrTokenNotFound     = errors.New("accessToken not found")
	ErrAuthEmptyToken    = errors.New("empty auth header")
	ErrInvalidAuthHeader = errors.New("authorization header is invalid")
)

// authMiddleware is a middleware that authenticates the user based on the accessToken in the request.
func (h *Handlers) authMiddleware(c *gin.Context) {
	session := sessions.Default(c)
	token, ok := session.Get(accessTokenName).(string)
	if !ok {
		newErrorResponse(c, http.StatusUnauthorized, ErrTokenNotFound.Error())
		return
	}

	if token == "" {
		newErrorResponse(c, http.StatusUnauthorized, ErrAuthEmptyToken.Error())
		return
	}

	payload, err := h.service.Auth.VerifyToken(c, token)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, ErrInvalidAuthHeader.Error())
		return
	}

	c.Set(authorizationPayloadKey, payload)
	c.Next()
}

// loggingMiddleware is a middleware that logs the request.
func (h *Handlers) loggingMiddleware(c *gin.Context) {
	h.logger.Infof("%s: [%s] - %s ", time.Now().Format(time.RFC3339), c.Request.Method, c.Request.URL.String())

	c.Next()
}

// getUserId get the user id from the context and returns it and an error if it is not found.
func getUserId(c *gin.Context) (int, error) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload == nil {
		return 0, ErrUserIdNotFound
	}

	return authPayload.UserID, nil
}
