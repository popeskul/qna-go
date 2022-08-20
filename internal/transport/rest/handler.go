// Package rest defines the handlers for the rest transport.
// It's define v1 handlers for the rest transport.
package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/popeskul/qna-go/internal/services"
	v1 "github.com/popeskul/qna-go/internal/transport/rest/v1"
)

// Handlers defines the handlers with all the necessary dependencies.
type Handlers struct {
	service *services.Service
}

// NewHandler creates a new Handlers with the necessary dependencies.
func NewHandler(service *services.Service) *Handlers {
	return &Handlers{
		service: service,
	}
}

// Init initializes the rest transport.
// It creates the gin engine and registers all the handlers.
// It returns the gin engine.
func (h *Handlers) Init() *gin.Engine {
	router := gin.Default()

	h.initAPI(router)

	return router
}

// initAPI initializes the api handlers.
// It registers the v1 handlers.
func (h *Handlers) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.service)

	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
