package db

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"go.uber.org/fx"

	"github.com/yuktake/todo-webapp/domain/todo"
	"github.com/yuktake/todo-webapp/domain/user"
)

type Todo = todo.Todo
type User = user.User

// DBConfig を提供
type DBConfig struct {
	DNS string
}

func NewDBConfig() DBConfig {
	return DBConfig{
		DNS: os.Getenv("DATABASE_URL"),
	}
}

func InitDB(config DBConfig) (*sql.DB, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(config.DNS)))

	return sqldb, nil
}

// bun.DB を初期化 (スキーマ作成は別で管理)
func NewBunDB(sqldb *sql.DB) *bun.DB {
	db := bun.NewDB(sqldb, pgdialect.New())

	// クエリログ出力
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	return db
}

// スキーマを作成する関数
func CreateSchema(lc fx.Lifecycle, db *bun.DB) {
	ctx := context.Background()

	for _, model := range []interface{}{(*User)(nil), (*Todo)(nil)} {
		_, err := db.NewCreateTable().
			Model(model).
			IfNotExists().
			Exec(ctx)
		if err != nil {
			log.Fatalf("failed to create table for %T: %v", model, err)
		}
	}
}
