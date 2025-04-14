package service

import (
	mock_repository "avito2/internal/repository/mocks"
	"testing"

	"github.com/golang/mock/gomock"
)

type serviceFixtures struct {
	ctrl     *gomock.Controller
	svc      *Svc
	mockRepo *mock_repository.MockRepository
}

func setUp(t *testing.T) serviceFixtures {
	ctrl := gomock.NewController(t)
	mockRepo := mock_repository.NewMockRepository(ctrl)
	svc := NewService(mockRepo)
	return serviceFixtures{
		ctrl:     ctrl,
		svc:      svc,
		mockRepo: mockRepo,
	}
}

func (s *serviceFixtures) tearDown() {
	s.ctrl.Finish()
}
