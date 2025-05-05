package e2e

import (
	"context"
	"fmt"
	"net"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
	"go.uber.org/fx"

	"github.com/yuktake/todo-webapp/domain"
	"github.com/yuktake/todo-webapp/handler"
	"github.com/yuktake/todo-webapp/infrastructure"
	"github.com/yuktake/todo-webapp/internal"
	"github.com/yuktake/todo-webapp/logger"
	"github.com/yuktake/todo-webapp/router"
	"github.com/yuktake/todo-webapp/service"
)

var TestServer *httptest.Server

func StartTestServer(m *testing.M) *bun.DB {
	var e *echo.Echo
	var db *bun.DB

	app := fx.New(
		infrastructure.TestModule, // テストDBに接続する
		service.Module,
		handler.Module,
		domain.Module,
		router.Module,
		internal.Module,
		logger.Module,
		fx.Populate(&e),  // Echoインスタンスを取得
		fx.Populate(&db), // DBインスタンスを取得
	)

	startCtx, cancel := context.WithTimeout(context.Background(), fx.DefaultTimeout)
	defer cancel()

	// Fxアプリケーション開始
	if err := app.Start(startCtx); err != nil {
		panic(err)
	}

	// テスト用HTTPサーバを起動
	TestServer = httptest.NewServer(e)
	// サーバーが起動してリクエストを受け付けられるようになるまで待機
	waitForServerReady(TestServer.Listener.Addr().String())

	return db
}

func waitForServerReady(addr string) {
	// 最大5秒間、100ms間隔でリトライ
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			panic("server did not start listening in time")
		case <-ticker.C:
			conn, err := net.Dial("tcp", addr)
			if err == nil {
				conn.Close()
				return // 接続できた → サーバーready
			}
		}
	}
}

func SeedTestData(db *bun.DB) error {
	fmt.Println("Seeding test data...")

	// ctx := context.Background()

	// // トランザクションを開始
	// tx, err := db.BeginTx(ctx, nil)
	// if err != nil {
	// 	panic(err)
	// }

	// users := []user.User{}
	// todos := []todo.Todo{}

	// _, err = db.NewInsert().Model(&users).Exec(ctx)
	// if err != nil {
	// 	tx.Rollback()
	// 	panic(err)
	// }

	// _, err = db.NewInsert().Model(&todos).Exec(ctx)
	// if err != nil {
	// 	tx.Rollback()
	// 	panic(err)
	// }

	// // トランザクションをコミット
	// if err := tx.Commit(); err != nil {
	// 	tx.Rollback()
	// 	panic(err)
	// }

	fmt.Println("Seeding completed.")

	return nil
}

func CleanupAllTables(db *bun.DB) error {
	ctx := context.Background()

	// 外部キー制約を考慮して削除順番を制御するならここでリスト化
	tables := []string{
		"todos",
		"users",
		// other tables...
	}

	for _, table := range tables {
		if _, err := db.NewDelete().Table(table).Where("1=1").Exec(ctx); err != nil {
			return err
		}
	}

	return nil
}
