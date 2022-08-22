// Package v1 defines all middleware for the 1 version.
package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

const (
	authorizationHeader        = "Authorization"
	userContextKey             = "userId"
	unauthenticatedUserError   = "user is unauthenticated"
	invalidAuthorizationHeader = "Authorization header is invalid"
)

var (
	ErrUserIsNotAuthorized = errors.New("user not authorized")
	ErrUserIdNotFound      = errors.New("user id not found")
)

// authMiddleware is a middleware that authenticates the user.
// If the user is not authenticated, the request is aborted with an error.
// If the user is authenticated, the user id is stored in the context.
// The user id is used in the handlers.
// The user id is stored in the context to avoid passing it as a parameter.
func (h *Handlers) authMiddleware(c *gin.Context) {
	token, err := getTokenFromRequest(c.Request)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, invalidAuthorizationHeader)
		return
	}

	userId, err := h.service.Auth.ParseToken(token)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, unauthenticatedUserError)
		return
	}

	c.Set(userContextKey, userId)

	c.Next()
}

// loggingMiddleware is a middleware that logs the request.
func (h *Handlers) loggingMiddleware(c *gin.Context) {
	h.logger.Infof("%s: [%s] - %s ", time.Now().Format(time.RFC3339), c.Request.Method, c.Request.URL.String())

	c.Next()
}

// getTokenFromRequest extracts the token from the request.
// The token is extracted from the Authorization header.
// The token is expected to be in the format: Bearer <token>
// If the token is not in the expected format, an error is returned.
func getTokenFromRequest(r *http.Request) (string, error) {
	header := r.Header.Get(authorizationHeader)
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return headerParts[1], nil
}

// getUserId get the user id from the context.
// If the user id is not found, an error is returned.
// If the user id is zero, an error is returned.
// If the user id is found, it is returned.
func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get("userId")

	if !ok {
		newErrorResponse(c, http.StatusUnauthorized, unauthenticatedUserError)
		return 0, ErrUserIsNotAuthorized
	}

	idInt, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusUnauthorized, unauthenticatedUserError)
		return 0, ErrUserIdNotFound
	}

	return idInt, nil
}
