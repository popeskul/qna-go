// Package tests is a service with all business logic for tests.
package tests

import (
	"context"
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

// CreateTest create new test in db and return testID and error if test not found.
func (s *ServiceTests) CreateTest(ctx context.Context, userID int, test domain.Test) error {
	return s.repo.CreateTest(ctx, userID, test)
}

// GetTest get test from db by testID and return test and error if test not found.
func (s *ServiceTests) GetTest(ctx context.Context, testID int) (domain.Test, error) {
	return s.repo.GetTest(ctx, testID)
}

func (s *ServiceTests) GetAllTestsByUserID(ctx context.Context, userID int, args domain.GetAllTestsParams) ([]domain.Test, error) {
	return s.repo.GetAllTestsByUserID(ctx, userID, args)
}

// UpdateTestByID update test in db and return test and error if test not found.
func (s *ServiceTests) UpdateTestByID(ctx context.Context, testID int, test domain.Test) error {
	return s.repo.UpdateTestById(ctx, testID, test)
}

// DeleteTestByID delete test in db and return error if test not found.
func (s *ServiceTests) DeleteTestByID(ctx context.Context, testID int) error {
	return s.repo.DeleteTestById(ctx, testID)
}
