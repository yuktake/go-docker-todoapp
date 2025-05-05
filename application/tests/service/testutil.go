package service

import (
	"context"

	"github.com/uptrace/bun"
)

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
