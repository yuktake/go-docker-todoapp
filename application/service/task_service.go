package service

import (
	"github.com/yuktake/todo-webapp/domain/todo"
	"github.com/yuktake/todo-webapp/logger"

	"go.uber.org/fx"
)

type Todo = todo.Todo

type TodoService interface {
	CreateTodo(todo Todo) (Todo, error)
	GetTodoByID(id string) (Todo, error)
	GetTodos() ([]Todo, error)
	UpdateTodo(todo Todo) (Todo, error)
	DeleteTodoByID(id string) error
}

type todoService struct {
	repo   todo.TodoRepository
	logger logger.Logger
}

// `fx.In` で `TaskRepository` を自動 DI
type todoServiceParams struct {
	fx.In
	Repo   todo.TodoRepository
	Logger logger.Logger
}

func NewTodoService(params todoServiceParams) TodoService {
	return &todoService{repo: params.Repo, logger: params.Logger}
}

func (s *todoService) CreateTodo(todo Todo) (Todo, error) {
	todo, err := s.repo.CreateTodo(&todo)

	if err != nil {
		s.logger.Error("failed to create todo", err)
		return Todo{}, err
	}

	return todo, nil
}

func (s *todoService) GetTodoByID(id string) (Todo, error) {
	todo, err := s.repo.GetTodoByID(id)
	if err != nil {
		s.logger.Error("failed to get todo", err)
		return Todo{}, err
	}

	return todo, nil
}

func (s *todoService) GetTodos() ([]Todo, error) {
	todos, err := s.repo.GetTodos()
	if err != nil {
		s.logger.Error("failed to get todos", err)
		return nil, err
	}

	return todos, nil
}

func (s *todoService) UpdateTodo(todo Todo) (Todo, error) {
	todo, err := s.repo.UpdateTodo(todo)
	if err != nil {
		s.logger.Error("failed to update todo", err)
		return Todo{}, err
	}

	return todo, nil
}

func (s *todoService) DeleteTodoByID(id string) error {
	err := s.repo.DeleteTodoByID(id)
	if err != nil {
		s.logger.Error("failed to delete todo", err)
		return err
	}

	return nil
}
