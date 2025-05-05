package db

import (
	"fmt"

	"github.com/uptrace/bun"
)

// SeedDataは主にマスタデータなどをDBに挿入します
func SeedData(db *bun.DB) error {
	fmt.Println("Seeding data...")

	return nil
}
