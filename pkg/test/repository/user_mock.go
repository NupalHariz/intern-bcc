package test

import (
	context "context"
	domain "intern-bcc/domain"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockIUserRepository is a mock of IUserRepository interface.
type MockIUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIUserRepositoryMockRecorder
}

// MockIUserRepositoryMockRecorder is the mock recorder for MockIUserRepository.
type MockIUserRepositoryMockRecorder struct {
	mock *MockIUserRepository
}

// NewMockIUserRepository creates a new mock instance.
func NewMockIUserRepository(ctrl *gomock.Controller) *MockIUserRepository {
	mock := &MockIUserRepository{ctrl: ctrl}
	mock.recorder = &MockIUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIUserRepository) EXPECT() *MockIUserRepositoryMockRecorder {
	return m.recorder
}

// CreateHasMentor mocks base method.
func (m *MockIUserRepository) CreateHasMentor(mentor *domain.HasMentor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateHasMentor", mentor)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateHasMentor indicates an expected call of CreateHasMentor.
func (mr *MockIUserRepositoryMockRecorder) CreateHasMentor(mentor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateHasMentor", reflect.TypeOf((*MockIUserRepository)(nil).CreateHasMentor), mentor)
}

// CreatePasswordVerification mocks base method.
func (m *MockIUserRepository) CreatePasswordVerification(ctx context.Context, emailVerHash, userName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePasswordVerification", ctx, emailVerHash, userName)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreatePasswordVerification indicates an expected call of CreatePasswordVerification.
func (mr *MockIUserRepositoryMockRecorder) CreatePasswordVerification(ctx, emailVerHash, userName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePasswordVerification", reflect.TypeOf((*MockIUserRepository)(nil).CreatePasswordVerification), ctx, emailVerHash, userName)  
}

// DeleteLikeProduct mocks base method.
func (m *MockIUserRepository) DeleteLikeProduct(likedProduct *domain.LikeProduct) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLikeProduct", likedProduct)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLikeProduct indicates an expected call of DeleteLikeProduct.
func (mr *MockIUserRepositoryMockRecorder) DeleteLikeProduct(likedProduct interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLikeProduct", reflect.TypeOf((*MockIUserRepository)(nil).DeleteLikeProduct), likedProduct)
}

// GetLikeProduct mocks base method.
func (m *MockIUserRepository) GetLikeProduct(likedProduct *domain.LikeProduct, likeProductParam domain.LikeProduct) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLikeProduct", likedProduct, likeProductParam)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetLikeProduct indicates an expected call of GetLikeProduct.
func (mr *MockIUserRepositoryMockRecorder) GetLikeProduct(likedProduct, likeProductParam interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLikeProduct", reflect.TypeOf((*MockIUserRepository)(nil).GetLikeProduct), likedProduct, likeProductParam)
}

// GetLikeProducts mocks base method.
func (m *MockIUserRepository) GetLikeProducts(user *domain.Users, userId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLikeProducts", user, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetLikeProducts indicates an expected call of GetLikeProducts.
func (mr *MockIUserRepositoryMockRecorder) GetLikeProducts(user, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLikeProducts", reflect.TypeOf((*MockIUserRepository)(nil).GetLikeProducts), user, userId)
}

// GetOwnMentors mocks base method.
func (m *MockIUserRepository) GetOwnMentors(user *domain.Users, userId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOwnMentors", user, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetOwnMentors indicates an expected call of GetOwnMentors.
func (mr *MockIUserRepositoryMockRecorder) GetOwnMentors(user, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOwnMentors", reflect.TypeOf((*MockIUserRepository)(nil).GetOwnMentors), user, userId)
}

// GetOwnProducts mocks base method.
func (m *MockIUserRepository) GetOwnProducts(user *domain.Users, userId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOwnProducts", user, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetOwnProducts indicates an expected call of GetOwnProducts.
func (mr *MockIUserRepositoryMockRecorder) GetOwnProducts(user, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOwnProducts", reflect.TypeOf((*MockIUserRepository)(nil).GetOwnProducts), user, userId)
}

// GetPasswordVerification mocks base method.
func (m *MockIUserRepository) GetPasswordVerification(ctx context.Context, userName string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPasswordVerification", ctx, userName)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPasswordVerification indicates an expected call of GetPasswordVerification.
func (mr *MockIUserRepositoryMockRecorder) GetPasswordVerification(ctx, userName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPasswordVerification", reflect.TypeOf((*MockIUserRepository)(nil).GetPasswordVerification), ctx, userName)
}

// GetUser mocks base method.
func (m *MockIUserRepository) GetUser(user *domain.Users, param domain.UserParam) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", user, param)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetUser indicates an expected call of GetUser.
func (mr *MockIUserRepositoryMockRecorder) GetUser(user, param interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockIUserRepository)(nil).GetUser), user, param)
}

// LikeProduct mocks base method.
func (m *MockIUserRepository) LikeProduct(likeProduct *domain.LikeProduct) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LikeProduct", likeProduct)
	ret0, _ := ret[0].(error)
	return ret0
}

// LikeProduct indicates an expected call of LikeProduct.
func (mr *MockIUserRepositoryMockRecorder) LikeProduct(likeProduct interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LikeProduct", reflect.TypeOf((*MockIUserRepository)(nil).LikeProduct), likeProduct)
}

// Register mocks base method.
func (m *MockIUserRepository) Register(newUser *domain.Users) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", newUser)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockIUserRepositoryMockRecorder) Register(newUser interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockIUserRepository)(nil).Register), newUser)
}

// UpdateUser mocks base method.
func (m *MockIUserRepository) UpdateUser(userUpdate *domain.UserUpdate, userId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", userUpdate, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockIUserRepositoryMockRecorder) UpdateUser(userUpdate, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockIUserRepository)(nil).UpdateUser), userUpdate, userId)
}
