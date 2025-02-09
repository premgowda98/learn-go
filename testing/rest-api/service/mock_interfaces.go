// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go
//
// Generated by this command:
//
//	mockgen -destination=mock_interfaces.go -package=service -source=interface.go
//

// Package service is a generated GoMock package.
package service

import (
	reflect "reflect"
	models "test/restapi/models"

	gomock "go.uber.org/mock/gomock"
)

// MockUserStore is a mock of UserStore interface.
type MockUserStore struct {
	ctrl     *gomock.Controller
	recorder *MockUserStoreMockRecorder
	isgomock struct{}
}

// MockUserStoreMockRecorder is the mock recorder for MockUserStore.
type MockUserStoreMockRecorder struct {
	mock *MockUserStore
}

// NewMockUserStore creates a new mock instance.
func NewMockUserStore(ctrl *gomock.Controller) *MockUserStore {
	mock := &MockUserStore{ctrl: ctrl}
	mock.recorder = &MockUserStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserStore) EXPECT() *MockUserStoreMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUserStore) Create(u *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", u)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockUserStoreMockRecorder) Create(u any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserStore)(nil).Create), u)
}

// Get mocks base method.
func (m *MockUserStore) Get(id int) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", id)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockUserStoreMockRecorder) Get(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUserStore)(nil).Get), id)
}
