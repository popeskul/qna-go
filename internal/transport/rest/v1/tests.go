// Package v1 defines the handlers for the 1 version.
package v1

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/repository/tests"
	"net/http"
)

// Tests interface is implemented by the service.
type Tests interface {
	CreateTest(ctx context.Context, test domain.TestInput) error
	GetTestByID(ctx context.Context, id int) (domain.Test, error)
	GetAllTestsByCurrentUser(ctx context.Context, userID int, args domain.GetAllTestsParams) ([]domain.Test, error)
	UpdateTestByID(ctx context.Context, id int, test domain.TestInput) error
	DeleteTestByID(ctx context.Context, id int) error
}

func (h *Handlers) CreateTest(c *gin.Context) {
	userId, error := getUserId(c)
	if error != nil {
		newErrorResponse(c, http.StatusUnauthorized, error.Error())
		return
	}

	var test domain.TestInput
	if err := c.ShouldBindJSON(&test); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Tests.CreateTest(c, userId, test); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, statusResponse{"success"})
}

func (h *Handlers) GetTestByID(c *gin.Context) {
	if _, error := getUserId(c); error != nil {
		newErrorResponse(c, http.StatusUnauthorized, error.Error())
		return
	}

	var request domain.GetTestByIDRequest
	if err := c.ShouldBindUri(&request); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	test, err := h.service.Tests.GetTest(c, request.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			newErrorResponse(c, http.StatusNotFound, "test not found")
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, domain.GetTestByIDResponse{
		Test: test,
	})
}

func (h *Handlers) GetAllTestsByCurrentUser(c *gin.Context) {
	userID, error := getUserId(c)
	if error != nil {
		newErrorResponse(c, http.StatusUnauthorized, error.Error())
		return
	}

	var request domain.GetAllTestsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		fmt.Println(request, err)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	args := domain.GetAllTestsParams{
		Limit:  request.PageSize,
		Offset: (request.PageID - 1) * request.PageSize,
	}

	tests, err := h.service.Tests.GetAllTestsByCurrentUser(c, userID, args)
	if err != nil {
		if err == sql.ErrNoRows {
			newErrorResponse(c, http.StatusNotFound, "tests not found")
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, domain.AllTestResponse{
		Tests: tests,
	})
}

func (h *Handlers) UpdateTestByID(c *gin.Context) {
	if _, error := getUserId(c); error != nil {
		newErrorResponse(c, http.StatusUnauthorized, error.Error())
		return
	}

	var request domain.GetTestByIDRequest
	if err := c.ShouldBindUri(&request); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var test domain.TestInput
	if err := c.ShouldBindJSON(&test); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Tests.UpdateTestByID(c, request.ID, test); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"success"})
}

func (h *Handlers) DeleteTestByID(c *gin.Context) {
	if _, error := getUserId(c); error != nil {
		newErrorResponse(c, http.StatusUnauthorized, error.Error())
		return
	}

	var request domain.GetTestByIDRequest
	if err := c.ShouldBindUri(&request); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Tests.DeleteTestByID(c, request.ID); err != nil {
		if err == tests.ErrDeleteTest {
			newErrorResponse(c, http.StatusNotFound, tests.ErrDeleteTest.Error())
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"success"})
}
