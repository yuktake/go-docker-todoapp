package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"

	"github.com/yuktake/todo-webapp/domain"
	"github.com/yuktake/todo-webapp/domain/todo"
	"github.com/yuktake/todo-webapp/handler"
	"github.com/yuktake/todo-webapp/infrastructure"
	"github.com/yuktake/todo-webapp/logger"
	"github.com/yuktake/todo-webapp/router"
	"github.com/yuktake/todo-webapp/service"
	"github.com/yuktake/todo-webapp/validator"
)

type Todo = todo.Todo

// initEnvは.envファイルから環境変数を読み込みます
func initEnv() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf(".envファイルの読み込みに失敗しました: %v", err)
	}

	if os.Getenv("JWT_SECRET") == "" {
		return fmt.Errorf("JWT_SECRETが.envファイルに設定されていません")
	}

	return nil
}

func NewEcho() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = validator.NewValidator()

	return e
}

func main() {
	// 環境変数の初期化
	if err := initEnv(); err != nil {
		log.Fatal(err)
	}

	app := fx.New(
		infrastructure.Module,
		service.Module,
		handler.Module,
		domain.Module,
		router.Module,
		logger.Module,
		fx.Provide(NewEcho),
		fx.Invoke(func(e *echo.Echo) {
			e.Logger.Fatal(e.Start(":8000"))
		}),
	)

	app.Run()
}
