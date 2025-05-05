package service

import (
	"context"
	"os"
	"testing"

	"github.com/uptrace/bun"
	"go.uber.org/fx"

	"github.com/yuktake/todo-webapp/infrastructure"
)

var testDB *bun.DB

func TestMain(m *testing.M) {
	app := fx.New(
		infrastructure.TestModule,
		fx.Populate(&testDB),
	)

	startCtx := context.Background()
	if err := app.Start(startCtx); err != nil {
		panic(err)
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
