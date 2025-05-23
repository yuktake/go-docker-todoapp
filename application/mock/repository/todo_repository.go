// Code generated by MockGen. DO NOT EDIT.
// Source: domain/todo/todo_repository.go
//
// Generated by this command:
//
//	mockgen -source=domain/todo/todo_repository.go -destination=./mock/repository/todo_repository.go -package=repository
//

// Package repository is a generated GoMock package.
package repository

import (
	reflect "reflect"

	todo "github.com/yuktake/todo-webapp/domain/todo"
	gomock "go.uber.org/mock/gomock"
)

// MockTodoRepository is a mock of TodoRepository interface.
type MockTodoRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTodoRepositoryMockRecorder
	isgomock struct{}
}

// MockTodoRepositoryMockRecorder is the mock recorder for MockTodoRepository.
type MockTodoRepositoryMockRecorder struct {
	mock *MockTodoRepository
}

// NewMockTodoRepository creates a new mock instance.
func NewMockTodoRepository(ctrl *gomock.Controller) *MockTodoRepository {
	mock := &MockTodoRepository{ctrl: ctrl}
	mock.recorder = &MockTodoRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTodoRepository) EXPECT() *MockTodoRepositoryMockRecorder {
	return m.recorder
}

// CreateTodo mocks base method.
func (m *MockTodoRepository) CreateTodo(arg0 *todo.Todo) (todo.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTodo", arg0)
	ret0, _ := ret[0].(todo.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTodo indicates an expected call of CreateTodo.
func (mr *MockTodoRepositoryMockRecorder) CreateTodo(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTodo", reflect.TypeOf((*MockTodoRepository)(nil).CreateTodo), arg0)
}

// DeleteTodoByID mocks base method.
func (m *MockTodoRepository) DeleteTodoByID(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTodoByID", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTodoByID indicates an expected call of DeleteTodoByID.
func (mr *MockTodoRepositoryMockRecorder) DeleteTodoByID(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTodoByID", reflect.TypeOf((*MockTodoRepository)(nil).DeleteTodoByID), id)
}

// GetTodoByID mocks base method.
func (m *MockTodoRepository) GetTodoByID(id string) (todo.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTodoByID", id)
	ret0, _ := ret[0].(todo.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTodoByID indicates an expected call of GetTodoByID.
func (mr *MockTodoRepositoryMockRecorder) GetTodoByID(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTodoByID", reflect.TypeOf((*MockTodoRepository)(nil).GetTodoByID), id)
}

// GetTodos mocks base method.
func (m *MockTodoRepository) GetTodos() ([]todo.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTodos")
	ret0, _ := ret[0].([]todo.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTodos indicates an expected call of GetTodos.
func (mr *MockTodoRepositoryMockRecorder) GetTodos() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTodos", reflect.TypeOf((*MockTodoRepository)(nil).GetTodos))
}

// UpdateTodo mocks base method.
func (m *MockTodoRepository) UpdateTodo(arg0 todo.Todo) (todo.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTodo", arg0)
	ret0, _ := ret[0].(todo.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTodo indicates an expected call of UpdateTodo.
func (mr *MockTodoRepositoryMockRecorder) UpdateTodo(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTodo", reflect.TypeOf((*MockTodoRepository)(nil).UpdateTodo), arg0)
}
