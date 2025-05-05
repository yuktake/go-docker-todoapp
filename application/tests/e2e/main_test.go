package e2e

import (
	"os"

	"testing"

	"github.com/uptrace/bun"
)

var testDB *bun.DB

func TestMain(m *testing.M) {
	testDB = StartTestServer(m)

	// Seeder
	SeedTestData(testDB)

	// 実際のテストを実行
	code := m.Run()

	// 必要なら後処理
	// --- クリーンアップ ---
	err := CleanupAllTables(testDB)
	if err != nil {
		panic(err)
	}

	if err := testDB.Close(); err != nil {
		panic(err)
	}

	os.Exit(code)
}
