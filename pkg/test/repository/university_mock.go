package test

import (
	domain "intern-bcc/domain"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIUniversityRepository is a mock of IUniversityRepository interface.
type MockIUniversityRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIUniversityRepositoryMockRecorder
}

// MockIUniversityRepositoryMockRecorder is the mock recorder for MockIUniversityRepository.
type MockIUniversityRepositoryMockRecorder struct {
	mock *MockIUniversityRepository
}

// NewMockIUniversityRepository creates a new mock instance.
func NewMockIUniversityRepository(ctrl *gomock.Controller) *MockIUniversityRepository {
	mock := &MockIUniversityRepository{ctrl: ctrl}
	mock.recorder = &MockIUniversityRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIUniversityRepository) EXPECT() *MockIUniversityRepositoryMockRecorder {
	return m.recorder
}

// CreateUniversity mocks base method.
func (m *MockIUniversityRepository) CreateUniversity(university *domain.Universities) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUniversity", university)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUniversity indicates an expected call of CreateUniversity.
func (mr *MockIUniversityRepositoryMockRecorder) CreateUniversity(university *domain.Universities) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUniversity", reflect.TypeOf((*MockIUniversityRepository)(nil).CreateUniversity), university)
}

// GetUniversity mocks base method.
func (m *MockIUniversityRepository) GetUniversity(university *domain.Universities, universityParam domain.Universities) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUniversity", university, universityParam)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetUniversity indicates an expected call of GetUniversity.
func (mr *MockIUniversityRepositoryMockRecorder) GetUniversity(university, universityParam interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUniversity", reflect.TypeOf((*MockIUniversityRepository)(nil).GetUniversity), university, universityParam)    
}