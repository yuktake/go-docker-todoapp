package handler

import (
	"net/http"

	"github.com/yuktake/todo-webapp/domain/todo"
	"github.com/yuktake/todo-webapp/dto"
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
	var CreateTodoRequest dto.CreateTodoRequest

	// リクエストのJSONをパース
	err := c.Bind(&CreateTodoRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	todo := Todo{
		Content: CreateTodoRequest.Content,
		Done:    CreateTodoRequest.Done,
	}

	// サービスにTodo作成を依頼
	newTodo, err2 := h.Service.CreateTodo(todo)
	if err2 != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// JSONを返す
	return c.JSON(http.StatusCreated, newTodo)
}

func (h *TodoHandler) GetTodo(c echo.Context) error {
	id := c.Param("id")
	todo, err := h.Service.GetTodoByID(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// JSONを返す
	return c.JSON(http.StatusOK, todo)
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

func (h *TodoHandler) UpdateTodo(c echo.Context) error {
	id := c.Param("id")

	todo, err := h.Service.GetTodoByID(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err2 := c.Bind(&todo)
	if err2 != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// サービスにTodo作成を依頼
	newTodo, err3 := h.Service.UpdateTodo(todo)
	if err3 != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// JSONを返す
	return c.JSON(http.StatusOK, newTodo)
}

func (h *TodoHandler) DeleteTodo(c echo.Context) error {
	id := c.Param("id")

	_, err := h.Service.GetTodoByID(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// サービスにTodo削除を依頼
	err2 := h.Service.DeleteTodoByID(id)
	if err2 != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// JSONを返す
	return c.JSON(http.StatusNoContent, nil)
}
