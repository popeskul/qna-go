package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/popeskul/qna-go/internal/domain"
	"net/http"
	"strconv"
)

// CreateTest godoc
// @Summary Create test
// @Security ApiKeyAuth
// @Tags tests
// @Description Create test
// @ID create-test
// @Accept  json
// @Produce  json
// @Param test body domain.TestInput true "test"
// @Success 200 {object} domain.Test
// @Failure 400 {object} error: error.Error
// @Failure 500 {object} error: error.Error
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

	id, err := h.service.Tests.CreateTest(userId, test)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"id":     id,
	})
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
// @Failure 400 {object} error: error.Error
// @Failure 500 {object} error: error.Error
// @Router /tests/{id} [get]
func (h *Handlers) GetTestByID(c *gin.Context) {
	_, error := getUserId(c)
	if error != nil {
		newErrorResponse(c, http.StatusUnauthorized, error.Error())
		return
	}

	testID, err := getIdFromRequest(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	test, err := h.service.Tests.GetTest(testID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"test":   test,
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
// @Success 200 {object} domain.Test
// @Failure 400 {object} error: error.Error
// @Failure 500 {object} error: error.Error
// @Router /tests/{id} [put]
func (h *Handlers) UpdateTestByID(c *gin.Context) {
	if _, error := getUserId(c); error != nil {
		newErrorResponse(c, http.StatusUnauthorized, error.Error())
		return
	}

	testID, err := getIdFromRequest(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var test domain.TestInput
	if err := c.ShouldBindJSON(&test); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.service.Tests.UpdateTestByID(testID, test); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"success"})
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
// @Success 200 {object} domain.Test
// @Failure 400 {object} error: error.Error
// @Failure 500 {object} error: error.Error
// @Router /tests/{id} [delete]
func (h *Handlers) DeleteTestByID(c *gin.Context) {
	if _, error := getUserId(c); error != nil {
		newErrorResponse(c, http.StatusUnauthorized, error.Error())
		return
	}

	testID, err := getIdFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.service.Tests.DeleteTestByID(testID); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"success"})
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
