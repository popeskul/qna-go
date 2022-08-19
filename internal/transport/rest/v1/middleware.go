package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader        = "Authorization"
	userContextKey             = "userId"
	unauthenticatedUserError   = "user is unauthenticated"
	invalidAuthorizationHeader = "Authorization header is invalid"
)

type CtxValue int

var (
	ErrUserIsNotAuthorized = errors.New("user not authorized")
	ErrUserIdNotFound      = errors.New("user id not found")
)

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
