// Package v1 defines the handlers for the 1 version.
package v1

import (
	"context"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/popeskul/qna-go/internal/domain"
	"net/http"
	"os"
	"time"
)

// Auth interface is implemented by the service.
type Auth interface {
	SignUp(ctx context.Context, user domain.User) error
	SignIn(ctx context.Context, user domain.User) error
	Refresh(ctx context.Context, sessionKey string) error
}

func (h *Handlers) SignUp(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Auth.CreateUser(c, user); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handlers) SignIn(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, refreshToken, err := h.service.Auth.SignIn(c, user)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Set-Cookie", fmt.Sprintf("refresh-accessToken='%s'; HttpOnly", refreshToken))
	if err = updateSession(c, accessToken); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken": accessToken,
	})
}

func (h *Handlers) Refresh(c *gin.Context) {
	token, err := c.Cookie("refresh-accessToken")
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	accessToken, refreshToken, err := h.service.GenerateAccessRefreshTokens(c, token)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Set-Cookie", fmt.Sprintf("refresh-accessToken='%s'; HttpOnly", refreshToken))

	if err = updateSession(c, accessToken); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken": accessToken,
	})
}

func updateSession(c *gin.Context, token string) error {
	sessionTTL, err := time.ParseDuration(os.Getenv("SESSION_HOUR_TTL"))
	if err != nil {
		return err
	}

	session := sessions.Default(c)
	session.Options(sessions.Options{
		Path:     "/",
		MaxAge:   int(sessionTTL),
		HttpOnly: true,
	})
	session.Set(accessTokenName, token)

	if err = session.Save(); err != nil {
		return err
	}

	return nil
}
