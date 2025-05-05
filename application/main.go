package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"

	"github.com/yuktake/todo-webapp/domain"
	"github.com/yuktake/todo-webapp/domain/todo"
	"github.com/yuktake/todo-webapp/domain/user"
	"github.com/yuktake/todo-webapp/handler"
	"github.com/yuktake/todo-webapp/infrastructure"
	dbschema "github.com/yuktake/todo-webapp/infrastructure/db"
	"github.com/yuktake/todo-webapp/internal"
	"github.com/yuktake/todo-webapp/logger"
	"github.com/yuktake/todo-webapp/router"
	"github.com/yuktake/todo-webapp/service"
)

type Todo = todo.Todo
type User = user.User
type IndexedModel = domain.IndexedModel

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

func main() {
	if err := initEnv(); err != nil {
		log.Fatal(err)
	}

	// envファイルのAPP_ENVでappの起動方法を変える
	if os.Getenv("APP_ENV") == "development" {
		app := fx.New(
			infrastructure.Module,
			service.Module,
			handler.Module,
			domain.Module,
			router.Module,
			logger.Module,
			internal.Module,
			fx.Invoke(func(e *echo.Echo) {
				e.Logger.Fatal(e.Start(":8000"))
			}),
		)

		app.Run()
	} else if os.Getenv("APP_ENV") == "production" {
		// 本番環境用の設定をここに追加

		cliApp := &cli.App{
			Name:  "todo-webapp",
			Usage: "Web server and migration management",
			Commands: []*cli.Command{
				{
					Name:  "serve",
					Usage: "Start the Echo server",
					Action: func(c *cli.Context) error {
						fxApp := fx.New(
							infrastructure.Module,
							service.Module,
							handler.Module,
							domain.Module,
							router.Module,
							logger.Module,
							internal.Module,
							fx.Invoke(func(e *echo.Echo) {
								e.Logger.Fatal(e.Start(":8000"))
							}),
						)
						fxApp.Run()
						return nil
					},
				},
				{
					Name:  "generate-schema",
					Usage: "Generate database schema",
					Action: func(c *cli.Context) error {
						fmt.Println("Generating database schema...")
						var db *bun.DB

						// 依存解決だけしたいので fx.Populate を使う
						var fxApp = fx.New(
							infrastructure.Module,
							fx.Populate(&db),
						)
						// 起動（依存解決）開始
						if err := fxApp.Start(context.Background()); err != nil {
							return fmt.Errorf("failed to start fx: %w", err)
						}
						defer fxApp.Stop(context.Background())

						sql := dbschema.ModelsToByte(db)

						// インデックスの生成
						index := dbschema.IndexesToByte(db)
						sql = append(sql, index...)

						if err := os.WriteFile("schema.sql", sql, 0644); err != nil {
							return fmt.Errorf("failed to write schema.sql: %w", err)
						}

						fmt.Println("Schema written to schema.sql")

						return nil
					},
				},
				{
					Name:  "seed",
					Usage: "Seed the database",
					Action: func(c *cli.Context) error {
						fmt.Println("Seeding database...")
						var db *bun.DB

						// 依存解決だけしたいので fx.Populate を使う
						var fxApp = fx.New(
							infrastructure.Module,
							fx.Populate(&db),
						)
						if err := fxApp.Start(context.Background()); err != nil {
							return fmt.Errorf("failed to start fx: %w", err)
						}
						defer fxApp.Stop(context.Background())

						// データベースのシーディング処理をここに追加
						dbschema.SeedData(db)
						fmt.Println("Database seeded")

						return nil
					},
				},
			},
		}

		if err := cliApp.Run(os.Args); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("APP_ENVが設定されていません")
	}
}
