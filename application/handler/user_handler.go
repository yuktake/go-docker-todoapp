package handler

import (
	"net/http"

	"github.com/yuktake/todo-webapp/domain/user"
	"github.com/yuktake/todo-webapp/service"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type User = user.User

type UserHandler struct {
	Service service.UserService
}

// `fx.In` で `UserService` を DI
type userHandlerParams struct {
	fx.In
	Service service.UserService
}

func NewUserHandler(params userHandlerParams) *UserHandler {
	return &UserHandler{Service: params.Service}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var user User

	// リクエストのJSONをパース
	err := c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// サービスにUser作成を依頼
	newUser, err2 := h.Service.CreateUser(user)
	if err2 != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// JSONを返す
	return c.JSON(http.StatusCreated, newUser)
}

func (h *UserHandler) GetUser(c echo.Context) error {
	id := c.Param("id")
	user, err := h.Service.GetUserByID(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// JSONを返す
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	users, err := h.Service.GetUsers()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// JSONを返す
	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	var user User

	// リクエストのJSONをパース
	err := c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// サービスにUser更新を依頼
	newUser, err2 := h.Service.UpdateUser(user)
	if err2 != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// JSONを返す
	return c.JSON(http.StatusOK, newUser)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")

	// サービスにUser削除を依頼
	err := h.Service.DeleteUserByID(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// JSONを返す
	return c.JSON(http.StatusNoContent, nil)
}
