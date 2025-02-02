package handler

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

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
	// リクエストのJSONをパース
	req := new(User)
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
	err2 := CheckHashPassword(user.Password, req.Password)
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
	var user User

	// リクエストのJSONをパース
	err := c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// パスワードを暗号化
	hashPassword, err := PasswordEncrypt(user.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "パスワードの暗号化に失敗しました"})
	}
	user.Password = hashPassword

	newUser, err := h.UserService.CreateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "ユーザー登録に失敗しました"})
	}

	// JSONを返す
	return c.JSON(http.StatusCreated, echo.Map{
		"message": "ユーザー登録が完了しました。ログインを行ってください",
		"user":    newUser,
	})
}

// 暗号化 (hash)
func PasswordEncrypt(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashPassword), err
}

// 暗号化パスワードと比較
func CheckHashPassword(hashPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}
