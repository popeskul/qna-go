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

// CreateTest godoc
// @Summary Create test
// @Security ApiKeyAuth
// @Tags tests
// @Description Create test
// @ID create-test
// @Accept  json
// @Produce  json
// @Param test body domain.TestInput true "test"
// @Success 201
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /tests [post]
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

	c.Status(http.StatusCreated)
}

// GetTestByID godoc
// @Summary Get test by id
// @Tags tests
// @Security ApiKeyAuth
// @Description Get test by id
// @ID get-test-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "id"
// @Success 200 {object} getTestByIDResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /tests/{id} [get]
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

// GetAllTestsByCurrentUser godoc
// @Summary Get all tests by current user
// @Tags tests
// @Security ApiKeyAuth
// @Description Get all tests by current user
// @ID get-all-tests-by-current-user
// @Accept  json
// @Produce  json
// @Param page_id query int false "page id"
// @Param page_size query int false "page size"
// @Success 200 {object} getAllTestsByCurrentUserResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /tests [get]
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

// UpdateTestByID godoc
// @Summary Update test by id
// @Tags tests
// @Security ApiKeyAuth
// @Description Update test by id
// @ID update-test-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "id"
// @Param test body domain.TestInput true "test"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /tests/{id} [put]
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

	c.Status(http.StatusOK)
}

// DeleteTestByID godoc
// @Summary Delete test by id
// @Tags tests
// @Description Delete test by id
// @ID delete-test-by-id
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param id path int true "id"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /tests/{id} [delete]
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

	c.Status(http.StatusOK)
}
