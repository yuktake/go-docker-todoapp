package db

import (
	"context"
	"database/sql"
	"os"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
	"go.uber.org/fx"

	"github.com/yuktake/todo-webapp/domain/todo"
)

type Todo = todo.Todo

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
	sqldb, err := sql.Open("postgres", config.DNS)
	if err != nil {
		return nil, err
	}

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
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			_, err := db.NewCreateTable().
				Model((*Todo)(nil)).
				IfNotExists().
				Exec(ctx)
			if err != nil {
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})
}
