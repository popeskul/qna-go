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
	v1 := api.Group("/v1")
	{
		h.InitAuthRouter(v1)
		h.InitTestsRouter(v1)
	}

	return v1
}
