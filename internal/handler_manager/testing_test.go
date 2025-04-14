package handler_manager

import (
	mock_service "avito2/internal/service/mocks"
	mock_jwt "avito2/internal/utils/mocks"
	"testing"

	"github.com/golang/mock/gomock"
)

type handlerManagerFixtures struct {
	ctrl       *gomock.Controller
	hm         HandlerManager
	mockSvc    *mock_service.MockService
	mockJWTGen *mock_jwt.MockJWTGenerator
}

func setUp(t *testing.T) handlerManagerFixtures {
	ctrl := gomock.NewController(t)
	mockSvc := mock_service.NewMockService(ctrl)
	mockJWTGen := mock_jwt.NewMockJWTGenerator(ctrl)
	hm := NewHandlerManager(mockSvc, mockJWTGen)
	return handlerManagerFixtures{
		ctrl:       ctrl,
		hm:         *hm,
		mockSvc:    mockSvc,
		mockJWTGen: mockJWTGen,
	}
}

func (h *handlerManagerFixtures) tearDown() {
	h.ctrl.Finish()
}
