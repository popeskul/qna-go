// Package tests is a service with all business logic for tests.
package tests

import (
	"context"
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/repository"

	"github.com/popeskul/cache"
)

// ServiceTests compose all functions for tests.
type ServiceTests struct {
	repo  repository.Tests
	cache *cache.Cache
}

// NewServiceTests create service with all fields.
func NewServiceTests(repo repository.Tests, cache *cache.Cache) *ServiceTests {
	return &ServiceTests{
		repo:  repo,
		cache: cache,
	}
}

// CreateTest create new test in db and return testID and error if test not found.
func (s *ServiceTests) CreateTest(ctx context.Context, userID int, test domain.Test) error {
	return s.repo.CreateTest(ctx, userID, test)
}

// GetTest get test from db by testID and return test and error if test not found.
func (s *ServiceTests) GetTest(ctx context.Context, testID int) (domain.Test, error) {
	if v, ok := s.cache.Get(testID); ok {
		if test, ok := v.(domain.Test); ok {
			return test, nil
		}
	}

	test, err := s.repo.GetTest(ctx, testID)
	if err != nil {
		return domain.Test{}, err
	}

	s.cache.Set(testID, test)

	return test, nil
}

func (s *ServiceTests) GetAllTestsByUserID(ctx context.Context, userID int, args domain.GetAllTestsParams) ([]domain.Test, error) {
	return s.repo.GetAllTestsByUserID(ctx, userID, args)
}

// UpdateTestByID update test in db and return test and error if test not found.
func (s *ServiceTests) UpdateTestByID(ctx context.Context, testID int, test domain.Test) error {
	if err := s.repo.UpdateTestById(ctx, testID, test); err != nil {
		return err
	}

	if v, ok := s.cache.Get(testID); ok {
		cachedTest, ok := v.(domain.Test)
		if !ok {
			return nil
		}

		if cachedTest.Title != test.Title && test.Title != "" {
			cachedTest.Title = test.Title
		}
		if cachedTest.AuthorID != test.AuthorID && test.AuthorID != 0 {
			cachedTest.AuthorID = test.AuthorID
		}
		s.cache.Set(testID, cachedTest)
	}

	return nil
}

// DeleteTestByID delete test in db and return error if test not found.
func (s *ServiceTests) DeleteTestByID(ctx context.Context, testID int) error {
	if err := s.repo.DeleteTestById(ctx, testID); err != nil {
		return err
	}

	s.cache.Delete(testID)

	return nil
}
