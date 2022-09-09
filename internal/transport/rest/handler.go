// Package rest defines the handlers for the rest transport.
// It's define v1 handlers for the rest transport.
package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/popeskul/qna-go/internal/logger"
	"github.com/popeskul/qna-go/docs"
	"github.com/popeskul/qna-go/internal/services"
	v1 "github.com/popeskul/qna-go/internal/transport/rest/v1"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Handlers defines the handlers with all the necessary dependencies.
type Handlers struct {
	service *services.Service
	logger  *logger.Logger
}

// NewHandler creates a new Handlers with the necessary dependencies.
func NewHandler(service *services.Service, logger *logger.Logger) *Handlers {
	return &Handlers{
		service: service,
		logger:  logger,
	}
}

// Init initializes the rest transport handlers and returns a gin engine.
func (h *Handlers) Init() *gin.Engine {
	router := gin.Default()

	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.Title = "Qna API"
	docs.SwaggerInfo.Description = "Qna API"
	docs.SwaggerInfo.Version = "1.0"

	apiV1 := router.Group("/api/v1")
	{
		handlersV1 := v1.NewHandler(h.service, h.logger)
		handlersV1.Init(apiV1)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router
}
