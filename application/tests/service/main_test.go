package service

import (
	"context"
	"os"
	"testing"

	"github.com/uptrace/bun"
	"go.uber.org/fx"

	"github.com/yuktake/todo-webapp/domain"
	"github.com/yuktake/todo-webapp/domain/todo"
	"github.com/yuktake/todo-webapp/domain/user"
	"github.com/yuktake/todo-webapp/infrastructure"
)

var testDB *bun.DB
var userRepository user.UserRepository
var todoRepository todo.TodoRepository

func TestMain(m *testing.M) {
	app := fx.New(
		infrastructure.TestModule,
		domain.Module,
		fx.Populate(&testDB),
		fx.Populate(&userRepository),
		fx.Populate(&todoRepository),
	)

	startCtx := context.Background()
	if err := app.Start(startCtx); err != nil {
		panic("failed to start app: " + err.Error())
	}

	// 実際のテストを実行
	code := m.Run()

	stopCtx := context.Background()
	if err := app.Stop(stopCtx); err != nil {
		panic(err)
	}

	err := testDB.Close()
	if err != nil {
		panic(err)
	}

	os.Exit(code)
}
