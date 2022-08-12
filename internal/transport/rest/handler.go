package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/popeskul/qna-go/internal/services"
	v1 "github.com/popeskul/qna-go/internal/transport/rest/v1"
)

type Handlers struct {
	service *services.Service
}

func NewHandler(service *services.Service) *Handlers {
	return &Handlers{
		service: service,
	}
}

func (h *Handlers) Init() *gin.Engine {
	router := gin.Default()

	h.initAPI(router)

	return router
}

func (h *Handlers) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.service)

	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
