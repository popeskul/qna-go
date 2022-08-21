// Package v1 defines the handlers for the 1 version.
package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/popeskul/qna-go/internal/services"
)

// Handlers defines handlers for v1
type Handlers struct {
	service *services.Service
}

// NewHandler creates a new Handlers with the necessary dependencies.
func NewHandler(service *services.Service) *Handlers {
	return &Handlers{
		service: service,
	}
}

// Init initializes routes for v1
func (h *Handlers) Init(api *gin.RouterGroup) *gin.RouterGroup {
	authAPI := api.Group("/auth")
	{
		authAPI.POST("/sign-up", h.SignUp)
		authAPI.POST("/sign-in", h.SignIn)
	}

	testsAPI := api.Group("/tests", h.authMiddleware)
	{
		testsAPI.POST("/", h.CreateTest)
		testsAPI.GET("/:id", h.GetTestByID)
		testsAPI.PUT("/:id", h.UpdateTestByID)
		testsAPI.DELETE("/:id", h.DeleteTestByID)
	}

	return api
}
