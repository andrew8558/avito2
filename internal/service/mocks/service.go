// Code generated by MockGen. DO NOT EDIT.
// Source: ./service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	model "avito2/internal/model"
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// AddProduct mocks base method.
func (m *MockService) AddProduct(ctx context.Context, pvzId uuid.UUID, productType model.ProductType) (*model.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddProduct", ctx, pvzId, productType)
	ret0, _ := ret[0].(*model.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddProduct indicates an expected call of AddProduct.
func (mr *MockServiceMockRecorder) AddProduct(ctx, pvzId, productType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddProduct", reflect.TypeOf((*MockService)(nil).AddProduct), ctx, pvzId, productType)
}

// CloseLastReception mocks base method.
func (m *MockService) CloseLastReception(ctx context.Context, pvzId uuid.UUID) (*model.Reception, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseLastReception", ctx, pvzId)
	ret0, _ := ret[0].(*model.Reception)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloseLastReception indicates an expected call of CloseLastReception.
func (mr *MockServiceMockRecorder) CloseLastReception(ctx, pvzId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseLastReception", reflect.TypeOf((*MockService)(nil).CloseLastReception), ctx, pvzId)
}

// CreatePvz mocks base method.
func (m *MockService) CreatePvz(ctx context.Context, city model.City) (*model.Pvz, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePvz", ctx, city)
	ret0, _ := ret[0].(*model.Pvz)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePvz indicates an expected call of CreatePvz.
func (mr *MockServiceMockRecorder) CreatePvz(ctx, city interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePvz", reflect.TypeOf((*MockService)(nil).CreatePvz), ctx, city)
}

// CreateReception mocks base method.
func (m *MockService) CreateReception(ctx context.Context, pvzId uuid.UUID) (*model.Reception, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateReception", ctx, pvzId)
	ret0, _ := ret[0].(*model.Reception)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateReception indicates an expected call of CreateReception.
func (mr *MockServiceMockRecorder) CreateReception(ctx, pvzId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateReception", reflect.TypeOf((*MockService)(nil).CreateReception), ctx, pvzId)
}

// DeleteLastProduct mocks base method.
func (m *MockService) DeleteLastProduct(ctx context.Context, pvzId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLastProduct", ctx, pvzId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLastProduct indicates an expected call of DeleteLastProduct.
func (mr *MockServiceMockRecorder) DeleteLastProduct(ctx, pvzId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLastProduct", reflect.TypeOf((*MockService)(nil).DeleteLastProduct), ctx, pvzId)
}

// GetPvzInfo mocks base method.
func (m *MockService) GetPvzInfo(ctx context.Context, startDate, endDate time.Time, page, limit int32) (*model.GetPvzInfoResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPvzInfo", ctx, startDate, endDate, page, limit)
	ret0, _ := ret[0].(*model.GetPvzInfoResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPvzInfo indicates an expected call of GetPvzInfo.
func (mr *MockServiceMockRecorder) GetPvzInfo(ctx, startDate, endDate, page, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPvzInfo", reflect.TypeOf((*MockService)(nil).GetPvzInfo), ctx, startDate, endDate, page, limit)
}
