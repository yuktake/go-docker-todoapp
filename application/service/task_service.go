package service

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yuktake/todo-webapp/domain/todo"

	"go.uber.org/fx"
)

type Todo = todo.Todo

type TodoService interface {
	CreateTodo() error
}

type todoService struct {
	repo todo.TodoRepository
}

// `fx.In` で `TaskRepository` を自動 DI
type todoServiceParams struct {
	fx.In
	Repo todo.TodoRepository
}

func NewTodoService(params todoServiceParams) TodoService {
	return &todoService{repo: params.Repo}
}

func (s *todoService) CreateTodo() error {
	return s.repo.CreateTodo()
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
