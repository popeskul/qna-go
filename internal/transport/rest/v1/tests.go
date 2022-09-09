// Package v1 defines the handlers for the 1 version.
package v1

import (
	"context"
	"database/sql"
	"errors"
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

// CreateTest godoc
// @Summary Create test
// @Security ApiKeyAuth
// @Tags tests
// @Description Create test
// @ID create-test
// @Accept  json
// @Produce  json
// @Param test body domain.Test true "test"
// @Success 201
// @Failure 400,401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /tests [post]
func (h *Handlers) CreateTest(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	if userId == 0 {
		newErrorResponse(c, http.StatusBadRequest, "user id is required")
		return
	}

	var test domain.Test
	if err = c.ShouldBindJSON(&test); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if test.Title == "" {
		newErrorResponse(c, http.StatusBadRequest, "title is required")
		return
	}

	if err = h.service.Tests.CreateTest(c, userId, test); err != nil {
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
// @Success 200 {object} domain.Test
// @Failure 400,401,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /tests/{id} [get]
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
		newErrorResponse(c, http.StatusUnauthorized, "you are not allowed to get this test")
		return
	}

	c.JSON(http.StatusOK, test)
}

// GetAllTestsByUserID godoc
// @Summary Get all tests by current user
// @Tags tests
// @Security ApiKeyAuth
// @Description Get all tests by current user
// @ID get-all-tests-by-current-user
// @Accept  json
// @Produce  json
// @Param page_id query int false "page id"
// @Param page_size query int false "page size"
// @Success 200 {object} []domain.Test
// @Failure 400,401,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /tests [get]
func (h *Handlers) GetAllTestsByUserID(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	if userID == 0 {
		newErrorResponse(c, http.StatusBadRequest, "user id is required")
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

// UpdateTestByID godoc
// @Summary Update test by id
// @Tags tests
// @Security ApiKeyAuth
// @Description Update test by id
// @ID update-test-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "id"
// @Param test body domain.Test true "test"
// @Success 200
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /tests/{id} [put]
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

// DeleteTestByID godoc
// @Summary Delete test by id
// @Tags tests
// @Description Delete test by id
// @ID delete-test-by-id
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param id path int true "id"
// @Success 200
// @Failure 400,401,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /tests/{id} [delete]
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
		if errors.Unwrap(err) == tests.ErrTest {
			newErrorResponse(c, http.StatusNotFound, tests.ErrTest.Error())
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
