// Code generated by MockGen. DO NOT EDIT.
// Source: service/user.go

// Package mock is a generated GoMock package.
package mock

import (
	model "github.com/amikai/gogolive/model"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockIUserService is a mock of IUserService interface
type MockIUserService struct {
	ctrl     *gomock.Controller
	recorder *MockIUserServiceMockRecorder
}

// MockIUserServiceMockRecorder is the mock recorder for MockIUserService
type MockIUserServiceMockRecorder struct {
	mock *MockIUserService
}

// NewMockIUserService creates a new mock instance
func NewMockIUserService(ctrl *gomock.Controller) *MockIUserService {
	mock := &MockIUserService{ctrl: ctrl}
	mock.recorder = &MockIUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIUserService) EXPECT() *MockIUserServiceMockRecorder {
	return m.recorder
}

// Register mocks base method
func (m *MockIUserService) Register(user model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register
func (mr *MockIUserServiceMockRecorder) Register(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockIUserService)(nil).Register), user)
}

// VerifyPassword mocks base method
func (m *MockIUserService) VerifyPassword(user model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyPassword", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyPassword indicates an expected call of VerifyPassword
func (mr *MockIUserServiceMockRecorder) VerifyPassword(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyPassword", reflect.TypeOf((*MockIUserService)(nil).VerifyPassword), user)
}
