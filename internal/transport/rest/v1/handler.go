// Package v1 defines the handlers for the 1 version.
package v1

import (
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/popeskul/qna-go/internal/logger"
	"github.com/popeskul/qna-go/internal/services"
)

// Handlers defines handlers for v1
type Handlers struct {
	service *services.Service
	store   cookie.Store
	logger  *logger.Logger
}

// NewHandler creates a new Handlers with the necessary dependencies.
func NewHandler(service *services.Service, store cookie.Store, log *logger.Logger) *Handlers {
	return &Handlers{
		service: service,
		store:   store,
		logger:  log,
	}
}

// Init initializes routes for v1
func (h *Handlers) Init(api *gin.RouterGroup) *gin.RouterGroup {
	api.Use(h.loggingMiddleware)

	authAPI := api.Group("/auth")
	{
		authAPI.POST("/sign-up", h.SignUp)
		authAPI.POST("/sign-in", h.SignIn)
		authAPI.GET("/refresh", h.Refresh)
	}

	testsAPI := api.Group("/tests", h.authMiddleware)
	{
		testsAPI.POST("/", h.CreateTest)
		testsAPI.GET("/", h.GetAllTestsByUserID)
		testsAPI.GET("/:id", h.GetTestByID)
		testsAPI.PUT("/:id", h.UpdateTestByID)
		testsAPI.DELETE("/:id", h.DeleteTestByID)
	}

	return api
}
