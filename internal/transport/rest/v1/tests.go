// Package v1 defines the handlers for the 1 version.
package v1

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/repository/tests"
	"net/http"
	"strconv"
)

// Tests interface is implemented by the service.
type Tests interface {
	CreateTest(ctx context.Context, test domain.Test) error
	GetTestByID(ctx context.Context, id int) (domain.Test, error)
	GetAllTestsByUserID(ctx context.Context, userID int, params domain.GetAllTestsParams) ([]domain.Test, error)
	UpdateTestByID(ctx context.Context, id int, test domain.Test) error
	DeleteTestByID(ctx context.Context, id int) error
}

func (h *Handlers) CreateTest(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	var test domain.Test
	if err = c.ShouldBindJSON(&test); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if test.Title == "" {
		newErrorResponse(c, http.StatusBadRequest, ErrTestNotFound.Error())
		return
	}

	if err = h.service.Tests.CreateTest(c, userId, test); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handlers) GetTestByID(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	testID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	test, err := h.service.Tests.GetTest(c, testID)
	if err != nil {
		if err == tests.ErrTestNotFound {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if test.AuthorID != userId {
		newErrorResponse(c, http.StatusUnauthorized, ErrPermission.Error())
		return
	}

	c.JSON(http.StatusOK, test)
}

func (h *Handlers) GetAllTestsByUserID(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	var request domain.GetAllTestsRequest
	if err = c.ShouldBindQuery(&request); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	args := domain.GetAllTestsParams{
		Limit:  request.PageSize,
		Offset: (request.PageID - 1) * request.PageSize,
	}

	tests, err := h.service.Tests.GetAllTestsByUserID(c, userID, args)
	if err != nil {
		if err == sql.ErrNoRows {
			newErrorResponse(c, http.StatusNotFound, "tests not found")
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tests)
}

func (h *Handlers) UpdateTestByID(c *gin.Context) {
	testID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var test domain.Test
	if err := c.ShouldBindJSON(&test); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Tests.UpdateTestByID(c, testID, test); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handlers) DeleteTestByID(c *gin.Context) {
	testID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if testID == 0 {
		newErrorResponse(c, http.StatusBadRequest, "test id is required")
		return
	}

	if err = h.service.Tests.DeleteTestByID(c, testID); err != nil {
		if err == tests.ErrDeleteTest {
			newErrorResponse(c, http.StatusNotFound, tests.ErrDeleteTest.Error())
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
