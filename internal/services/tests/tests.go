package tests

import (
	"github.com/popeskul/qna-go/internal/domain"
	"github.com/popeskul/qna-go/internal/repository"
)

type ServiceTests struct {
	repo repository.Tests
}

func NewServiceTests(repo repository.Tests) *ServiceTests {
	return &ServiceTests{
		repo: repo,
	}
}

func (s *ServiceTests) CreateTest(userID int, test domain.TestInput) (int, error) {
	return s.repo.CreateTest(userID, test)
}

func (s *ServiceTests) GetTest(testID int) (domain.Test, error) {
	return s.repo.GetTest(testID)
}
