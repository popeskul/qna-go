package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/popeskul/qna-go/internal/services"
)

type Handlers struct {
	service *services.Service
}

func NewHandler(service *services.Service) *Handlers {
	return &Handlers{
		service: service,
	}
}

func (h *Handlers) Init(api *gin.RouterGroup) *gin.RouterGroup {
	v1 := api.Group("/v1")
	{
		h.InitAuthRouter(v1)
	}

	return v1
}
