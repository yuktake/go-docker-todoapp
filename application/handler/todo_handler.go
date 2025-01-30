package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yuktake/todo-webapp/service"
	"go.uber.org/fx"
)

type TodoHandler struct {
	Service service.TodoService
}

// `fx.In` で `TaskService` を DI
type todoHandlerParams struct {
	fx.In
	Service service.TodoService
}

func NewTodoHandler(params todoHandlerParams) *TodoHandler {
	return &TodoHandler{Service: params.Service}
}

func (h *TodoHandler) CreateTodo(c echo.Context) error {
	return nil
}

func (h *TodoHandler) GetTodos(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
