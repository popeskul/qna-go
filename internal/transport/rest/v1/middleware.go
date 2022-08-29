// Package v1 defines all middleware for the 1 version.
package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/popeskul/qna-go/internal/token"
	"net/http"
	"strings"
	"time"
)

const (
	authorizationHeader     = "Authorization"
	authorizationType       = "Bearer"
	authorizationPayloadKey = "authorization_payload"
)

var (
	ErrUserIdNotFound    = errors.New("user id not found")
	ErrTokenNotFound     = errors.New("token not found")
	ErrAuthEmptyToken    = errors.New("empty auth header")
	ErrInvalidAuthHeader = errors.New("authorization header is invalid")
)

// authMiddleware is a middleware that authenticates the user based on the token in the request.
func (h *Handlers) authMiddleware(c *gin.Context) {
	tokenFromRequest, err := getTokenFromRequest(c.Request)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, ErrInvalidAuthHeader.Error())
		return
	}

	payload, err := h.service.TokenMaker.VerifyToken(tokenFromRequest)
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

// getTokenFromRequest extracts the token from the request header.
func getTokenFromRequest(r *http.Request) (string, error) {
	header := r.Header.Get(authorizationHeader)
	if header == "" {
		return "", ErrAuthEmptyToken
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != authorizationType {
		return "", ErrInvalidAuthHeader
	}

	if len(headerParts[1]) == 0 {
		return "", ErrTokenNotFound
	}

	return headerParts[1], nil
}

// getUserId get the user id from the context and returns it and an error if it is not found.
func getUserId(c *gin.Context) (int, error) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload == nil {
		return 0, ErrUserIdNotFound
	}

	return authPayload.UserID, nil
}
