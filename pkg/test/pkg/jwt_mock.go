package test

import (
	domain "intern-bcc/domain"
	reflect "reflect"

	gin "github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockIJwt is a mock of IJwt interface.
type MockIJwt struct {
	ctrl     *gomock.Controller
	recorder *MockIJwtMockRecorder
}

// MockIJwtMockRecorder is the mock recorder for MockIJwt.
type MockIJwtMockRecorder struct {
	mock *MockIJwt
}

// NewMockIJwt creates a new mock instance.
func NewMockIJwt(ctrl *gomock.Controller) *MockIJwt {
	mock := &MockIJwt{ctrl: ctrl}
	mock.recorder = &MockIJwtMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIJwt) EXPECT() *MockIJwtMockRecorder {
	return m.recorder
}

// GenerateToken mocks base method.
func (m *MockIJwt) GenerateToken(userId uuid.UUID) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateToken", userId)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MockIJwtMockRecorder) GenerateToken(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockIJwt)(nil).GenerateToken), userId)
}

// GetLoginUser mocks base method.
func (m *MockIJwt) GetLoginUser(ctx *gin.Context) (domain.Users, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLoginUser", ctx)
	ret0, _ := ret[0].(domain.Users)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLoginUser indicates an expected call of GetLoginUser.
func (mr *MockIJwtMockRecorder) GetLoginUser(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLoginUser", reflect.TypeOf((*MockIJwt)(nil).GetLoginUser), ctx)
}

// ValidateToken mocks base method.
func (m *MockIJwt) ValidateToken(tokenString string) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateToken", tokenString)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateToken indicates an expected call of ValidateToken.
func (mr *MockIJwtMockRecorder) ValidateToken(tokenString interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateToken", reflect.TypeOf((*MockIJwt)(nil).ValidateToken), tokenString)
}
