package service

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yuktake/todo-webapp/domain/todo"
	"github.com/yuktake/todo-webapp/logger"

	"go.uber.org/fx"
)

type Todo = todo.Todo

type TodoService interface {
	CreateTodo(todo Todo) (Todo, error)
	GetTodos() ([]Todo, error)
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

func (s *todoService) GetTodos() ([]Todo, error) {
	todos, err := s.repo.GetTodos()
	if err != nil {
		s.logger.Error("failed to get todos", err)
		return nil, err
	}

	return todos, nil
}

func customFunc(todo *Todo) func([]string) []error {
	return func(values []string) []error {
		if len(values) == 0 || values[0] == "" {
			return nil
		}
		dt, err := time.Parse("2006-01-02T15:04 MST", values[0]+" JST")
		if err != nil {
			return []error{echo.NewBindingError("until", values[0:1], "failed to decode time", err)}
		}
		todo.Until = dt
		return nil
	}
}
