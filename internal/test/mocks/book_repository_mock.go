package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/osalomon89/go-inventory/internal/domain"
)

// MockBookRepository is a mock of BookRepository interface
type MockBookRepository struct {
	ctrl     *gomock.Controller
	recorder *MockBookRepositoryMockRecorder
}

// MockBookRepositoryMockRecorder is the mock recorder for MockBookRepository
type MockBookRepositoryMockRecorder struct {
	mock *MockBookRepository
}

// NewMockBookRepository creates a new mock instance
func NewMockBookRepository(ctrl *gomock.Controller) *MockBookRepository {
	mock := &MockBookRepository{ctrl: ctrl}
	mock.recorder = &MockBookRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBookRepository) EXPECT() *MockBookRepositoryMockRecorder {
	return m.recorder
}

// CreateBook mocks base method
func (m *MockBookRepository) CreateBook(arg0 *domain.Book) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBook", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateBook indicates an expected call of CreateBook
func (mr *MockBookRepositoryMockRecorder) CreateBook(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBook", reflect.TypeOf((*MockBookRepository)(nil).CreateBook), arg0)
}

// GetBookByID mocks base method
func (m *MockBookRepository) GetBookByID(arg0 uint) (*domain.Book, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBookByID", arg0)
	ret0, _ := ret[0].(*domain.Book)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBookByID indicates an expected call of GetBookByID
func (mr *MockBookRepositoryMockRecorder) GetBookByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBookByID", reflect.TypeOf((*MockBookRepository)(nil).GetBookByID), arg0)
}

// GetBooks mocks base method
func (m *MockBookRepository) GetBooks(arg0 map[string]interface{}) ([]domain.Book, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBooks", arg0)
	ret0, _ := ret[0].([]domain.Book)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBooks indicates an expected call of GetBooks
func (mr *MockBookRepositoryMockRecorder) GetBooks(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBooks", reflect.TypeOf((*MockBookRepository)(nil).GetBooks), arg0)
}

// UpdateBookByParams mocks base method
func (m *MockBookRepository) UpdateBookByParams(arg0 map[string]interface{}, arg1 *domain.Book) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBookByParams", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBookByParams indicates an expected call of UpdateBookByParams
func (mr *MockBookRepositoryMockRecorder) UpdateBookByParams(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBookByParams", reflect.TypeOf((*MockBookRepository)(nil).UpdateBookByParams), arg0, arg1)
}
