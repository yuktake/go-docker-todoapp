package handler

import (
	"net/http"

	"github.com/yuktake/todo-webapp/domain/todo"
	"github.com/yuktake/todo-webapp/service"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type Todo = todo.Todo

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
	var todo Todo

	// リクエストのJSONをパース
	err := c.Bind(&todo)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// サービスにTodo作成を依頼
	newTodo, err2 := h.Service.CreateTodo(todo)
	if err2 != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// JSONを返す
	return c.JSON(http.StatusCreated, newTodo)
}

func (h *TodoHandler) GetTodos(c echo.Context) error {
	var todos []Todo
	todos, err := h.Service.GetTodos()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// JSONを返す
	return c.JSON(http.StatusOK, todos)
}
