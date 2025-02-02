package router

import (
	"os"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	jwtv5 "github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/yuktake/todo-webapp/domain/auth"
	"github.com/yuktake/todo-webapp/domain/todo"
	"github.com/yuktake/todo-webapp/handler"
)

type Todo = todo.Todo

// ルーティング設定を行う関数
func RegisterRoutes(e *echo.Echo, todoHandler *handler.TodoHandler, authHandler *handler.AuthHandler) {
	// 環境変数からJWTシークレットを取得
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwtv5.Claims {
			return new(auth.JwtCustomClaims)
		},
		SigningKey: jwtSecret,
	})

	// パブリックルート: ログイン
	e.POST("/login", authHandler.Login)
	e.POST("/signup", authHandler.Signup)

	apiGroup := e.Group("/")
	apiGroup.Use(jwtMiddleware)

	apiGroup.GET("", todoHandler.GetTodos)

	apiGroup.POST("todo", todoHandler.CreateTodo)
	apiGroup.GET("todo/:id", todoHandler.GetTodo)
	apiGroup.PATCH("todo/:id", todoHandler.UpdateTodo)
	apiGroup.DELETE("todo/:id", todoHandler.DeleteTodo)
}

// Fx Module
var Module = fx.Module("router",
	fx.Invoke(RegisterRoutes),
)
