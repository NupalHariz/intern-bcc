package test

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIGoMail is a mock of IGoMail interface.
type MockIGoMail struct {
	ctrl     *gomock.Controller
	recorder *MockIGoMailMockRecorder
}

// MockIGoMailMockRecorder is the mock recorder for MockIGoMail.
type MockIGoMailMockRecorder struct {
	mock *MockIGoMail
}

// NewMockIGoMail creates a new mock instance.
func NewMockIGoMail(ctrl *gomock.Controller) *MockIGoMail {
	mock := &MockIGoMail{ctrl: ctrl}
	mock.recorder = &MockIGoMailMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIGoMail) EXPECT() *MockIGoMailMockRecorder {
	return m.recorder
}

// SendGoMail mocks base method.
func (m *MockIGoMail) SendGoMail(subject, htmlBody, toEmail string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendGoMail", subject, htmlBody, toEmail)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendGoMail indicates an expected call of SendGoMail.
func (mr *MockIGoMailMockRecorder) SendGoMail(subject, htmlBody, toEmail interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendGoMail", reflect.TypeOf((*MockIGoMail)(nil).SendGoMail), subject, htmlBody, toEmail)
}