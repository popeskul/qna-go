package v1

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/popeskul/qna-go/internal/domain"
	"net/http"
	"strconv"
)

type Tests interface {
	CreateTest(ctx context.Context, test domain.TestInput) error
}

func (h *Handlers) InitTestsRouter(v1 *gin.RouterGroup) {
	testsAPI := v1.Group("/tests", h.authMiddleware)
	{
		testsAPI.POST("/", h.CreateTest)
		testsAPI.GET("/:id", h.GetTest)
	}
}

func (h *Handlers) CreateTest(c *gin.Context) {
	userId, error := getUserId(c)
	if error != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": error})
		return
	}

	var test domain.TestInput
	if err := c.ShouldBindJSON(&test); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.service.Tests.CreateTest(userId, test)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"id":     id,
	})
}

func (h *Handlers) GetTest(c *gin.Context) {
	_, error := getUserId(c)
	if error != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": error})
		return
	}

	testID, err := getIdFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	test, err := h.service.Tests.GetTest(testID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"test":   test,
	})
}

func getIdFromRequest(r *gin.Context) (int, error) {
	id, err := strconv.Atoi(r.Param("id"))
	if err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, errors.New("id can't be 0")
	}

	return id, nil
}
