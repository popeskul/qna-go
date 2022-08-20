// Package tests is a service with all business logic for tests.
package tests

import (
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/repository"
)

// ServiceTests compose all functions for tests.
type ServiceTests struct {
	repo repository.Tests
}

// NewServiceTests create service with all fields.
func NewServiceTests(repo repository.Tests) *ServiceTests {
	return &ServiceTests{
		repo: repo,
	}
}

// CreateTest create new test in db.
// It's return error and test id if test created.
func (s *ServiceTests) CreateTest(userID int, test domain.TestInput) (int, error) {
	return s.repo.CreateTest(userID, test)
}

// GetTest get test from db.
// It's return domain.Test and error if test not found.
func (s *ServiceTests) GetTest(testID int) (domain.Test, error) {
	return s.repo.GetTest(testID)
}

// UpdateTestByID update test in db.
// It's return testID and error if test not found.
func (s *ServiceTests) UpdateTestByID(testID int, test domain.TestInput) error {
	return s.repo.UpdateTestById(testID, test)
}

// DeleteTestByID delete test in db.
// It's return an error if test not found.
func (s *ServiceTests) DeleteTestByID(testID int) error {
	return s.repo.DeleteTestById(testID)
}
