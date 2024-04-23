package test

import (
	context "context"
	domain "intern-bcc/domain"
	reflect "reflect"

	gin "github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockIProductRepository is a mock of IProductRepository interface.
type MockIProductRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIProductRepositoryMockRecorder
}

// MockIProductRepositoryMockRecorder is the mock recorder for MockIProductRepository.
type MockIProductRepositoryMockRecorder struct {
	mock *MockIProductRepository
}

// NewMockIProductRepository creates a new mock instance.
func NewMockIProductRepository(ctrl *gomock.Controller) *MockIProductRepository {
	mock := &MockIProductRepository{ctrl: ctrl}
	mock.recorder = &MockIProductRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIProductRepository) EXPECT() *MockIProductRepositoryMockRecorder {
	return m.recorder
}

// CreateProduct mocks base method.
func (m *MockIProductRepository) CreateProduct(newProduct *domain.Products) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProduct", newProduct)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateProduct indicates an expected call of CreateProduct.
func (mr *MockIProductRepositoryMockRecorder) CreateProduct(newProduct interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProduct", reflect.TypeOf((*MockIProductRepository)(nil).CreateProduct), newProduct)
}

// GetProduct mocks base method.
func (m *MockIProductRepository) GetProduct(product *domain.Products, productParam domain.ProductParam) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProduct", product, productParam)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetProduct indicates an expected call of GetProduct.
func (mr *MockIProductRepositoryMockRecorder) GetProduct(product, productParam interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProduct", reflect.TypeOf((*MockIProductRepository)(nil).GetProduct), product, productParam)
}

// GetProducts mocks base method.
func (m *MockIProductRepository) GetProducts(c *gin.Context, ctx context.Context, product *[]domain.Products, productParam domain.ProductParam) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProducts", c, ctx, product, productParam)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetProducts indicates an expected call of GetProducts.
func (mr *MockIProductRepositoryMockRecorder) GetProducts(c, ctx, product, productParam interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProducts", reflect.TypeOf((*MockIProductRepository)(nil).GetProducts), c, ctx, product, productParam)
}

// GetTotalProduct mocks base method.
func (m *MockIProductRepository) GetTotalProduct(totalProduct *int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTotalProduct", totalProduct)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetTotalProduct indicates an expected call of GetTotalProduct.
func (mr *MockIProductRepositoryMockRecorder) GetTotalProduct(totalProduct interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTotalProduct", reflect.TypeOf((*MockIProductRepository)(nil).GetTotalProduct), totalProduct)
}

// UpdateProduct mocks base method.
func (m *MockIProductRepository) UpdateProduct(product *domain.ProductUpdate, productId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProduct", product, productId)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProduct indicates an expected call of UpdateProduct.
func (mr *MockIProductRepositoryMockRecorder) UpdateProduct(product, productId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProduct", reflect.TypeOf((*MockIProductRepository)(nil).UpdateProduct), product, productId)
}