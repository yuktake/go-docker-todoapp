package handler

import (
	"net/http"

	"github.com/yuktake/todo-webapp/domain/user"
	"github.com/yuktake/todo-webapp/dto"
	"github.com/yuktake/todo-webapp/service"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type AuthHandler struct {
	AuthService service.AuthService
	UserService service.UserService
}

type authHandlerParams struct {
	fx.In
	AuthService service.AuthService
	UserService service.UserService
}

func NewAuthHandler(params authHandlerParams) *AuthHandler {
	return &AuthHandler{
		AuthService: params.AuthService,
		UserService: params.UserService,
	}
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// ユーザー情報をサービスに渡す
	user, err := h.UserService.GetUserByEmail(req.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "ユーザーが見つかりません"})
	}

	// パスワードを比較
	err2 := service.CheckHashPassword(user.Password, req.Password)
	if err2 != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "パスワードが違います"})
	}

	// JWTトークンを生成(JWT Serviceを使って)
	token, err3 := h.AuthService.CreateToken(user)
	if err3 != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "トークンの生成に失敗しました"})
	}

	// トークンを返す
	return c.JSON(http.StatusOK, echo.Map{"token": token})
}

func (h *AuthHandler) Signup(c echo.Context) error {
	var signup_request dto.SignupRequest

	// リクエストのJSONをパース
	err := c.Bind(&signup_request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// バリデーション
	err2 := c.Validate(signup_request)
	if err2 != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": err2.Error()})
	}

	// パスワードを暗号化
	hashPassword, err := service.PasswordEncrypt(signup_request.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "パスワードの暗号化に失敗しました"})
	}

	// dtoの値からUserエンティティを作成
	user := user.User{
		Name:     signup_request.Name,
		Password: hashPassword,
		Email:    signup_request.Email,
	}

	err3 := user.Validate()
	if err3 != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": err3.Error()})
	}

	newUser, err4 := h.UserService.CreateUser(user)
	if err4 != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "ユーザー登録に失敗しました"})
	}

	signup_response := dto.SignupResponse{
		Message: "ユーザー登録が完了しました。ログインを行ってください",
		User:    newUser,
	}

	// JSONを返す
	return c.JSON(http.StatusCreated, signup_response)
}
