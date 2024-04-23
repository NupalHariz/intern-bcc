package test

import (
	multipart "mime/multipart"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockISupabase is a mock of ISupabase interface.
type MockISupabase struct {
	ctrl     *gomock.Controller
	recorder *MockISupabaseMockRecorder
}

// MockISupabaseMockRecorder is the mock recorder for MockISupabase.
type MockISupabaseMockRecorder struct {
	mock *MockISupabase
}

// NewMockISupabase creates a new mock instance.
func NewMockISupabase(ctrl *gomock.Controller) *MockISupabase {
	mock := &MockISupabase{ctrl: ctrl}
	mock.recorder = &MockISupabaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockISupabase) EXPECT() *MockISupabaseMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockISupabase) Delete(link string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", link)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockISupabaseMockRecorder) Delete(link interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockISupabase)(nil).Delete), link)
}

// Upload mocks base method.
func (m *MockISupabase) Upload(file *multipart.FileHeader) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upload", file)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Upload indicates an expected call of Upload.
func (mr *MockISupabaseMockRecorder) Upload(file interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upload", reflect.TypeOf((*MockISupabase)(nil).Upload), file)
}